package main

import (
	"database/sql"
	"fmt"
	"log"
	"runtime"
	"strconv"
	"time"
	"unsafe"

	_ "github.com/lib/pq"
	"gonum.org/v1/gonum/mat"
)

// --------------------------------------------------
// CONFIGURATION
// --------------------------------------------------
const (
	MAX_ROWS    = 10
	MAX_COLS    = 11 // Nombre de champs dans votre table
	POOL_SIZE   = 1024
	BUFFER_SIZE = 4096
)

// --------------------------------------------------
// STRUCTURES DE DONNÉES
// --------------------------------------------------
type HotelData struct {
	encodedPtr uint64
	rawData    string
	length     int
}

type HotelMatrix struct {
	matrix     *mat.Dense
	dataPool   []HotelData
	index      map[uint64]int
	bufferPool [][]byte
}

// --------------------------------------------------
// INITIALISATION
// --------------------------------------------------
func NewHotelMatrix() *HotelMatrix {
	// Modification: Matrice de taille MAX_ROWS x MAX_COLS (au lieu de MAX_COLS*2)
	hm := &HotelMatrix{
		matrix:     mat.NewDense(MAX_ROWS, MAX_COLS, nil),
		dataPool:   make([]HotelData, 0, POOL_SIZE),
		index:      make(map[uint64]int),
		bufferPool: make([][]byte, POOL_SIZE),
	}

	for i := range hm.bufferPool {
		hm.bufferPool[i] = make([]byte, 0, BUFFER_SIZE)
	}
	return hm
}

// --------------------------------------------------
// ENCODAGE CORRIGÉ
// --------------------------------------------------
func (hm *HotelMatrix) encodeToMatrix(row int, data []string) error {
	if row >= MAX_ROWS {
		return fmt.Errorf("row index out of bounds")
	}

	// Modification: Vérification du nombre de colonnes
	if len(data) < MAX_COLS {
		return fmt.Errorf("not enough columns in data")
	}

	for col := 0; col < MAX_COLS; col++ {
		field := data[col]
		ptr, _ := hm.encodeField(field)
		hm.matrix.Set(row, col, float64(ptr)) // Modification: Utilisation des colonnes 0-10 seulement
	}
	return nil
}

func (hm *HotelMatrix) encodeField(data string) (uint64, float64) {
	buf := hm.bufferPool[len(hm.dataPool)%POOL_SIZE]
	buf = buf[:0]
	buf = append(buf, data...)

	ptr := uint64(uintptr(unsafe.Pointer(&buf[0])))

	hm.dataPool = append(hm.dataPool, HotelData{
		encodedPtr: ptr,
		rawData:    data,
		length:     len(data),
	})
	hm.index[ptr] = len(hm.dataPool) - 1

	return ptr, float64(len(data))
}

// --------------------------------------------------
// DÉCODAGE
// --------------------------------------------------
func (hm *HotelMatrix) decodeByPointer(ptr uint64) string {
	if idx, exists := hm.index[ptr]; exists {
		return hm.dataPool[idx].rawData
	}
	return ""
}

// --------------------------------------------------
// INTERFACE POSTGRES
// --------------------------------------------------
func fetchHotelData(db *sql.DB) ([][]string, error) {
	rows, err := db.Query(
		`SELECT hotel_key, name, city, country, stars, address, 
		  mail, latitude, longitude, phone, description, short_description 
		  FROM hotel LIMIT $1`, MAX_ROWS)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results [][]string
	for rows.Next() {
		var (
			key, name, city, country, address, mail, phone, desc, shortDesc string
			stars                                                           int
			latitude, longitude                                             float64
		)

		err := rows.Scan(&key, &name, &city, &country, &stars, &address, &mail,
			&latitude, &longitude, &phone, &desc, &shortDesc)
		if err != nil {
			return nil, err
		}

		record := []string{
			key, name, city, country, fmt.Sprint(stars), address, mail,
			fmt.Sprint(latitude), fmt.Sprint(longitude), phone, desc, shortDesc,
		}
		results = append(results, record)
	}
	return results, nil
}

/*
// Ajoutez cette structure pour l'index rapide
type HotelIndex struct {
	idToRow map[string]int // Map ID → Ligne dans la matrice
	matrix  *HotelMatrix   // Référence à votre matrice existante
}

func NewHotelIndex(hm *HotelMatrix, data [][]string) *HotelIndex {
	hi := &HotelIndex{
		idToRow: make(map[string]int),
		matrix:  hm,
	}

	for i, record := range data {
		if i >= MAX_ROWS {
			break
		}
		// L'ID est toujours dans la première colonne
		hi.idToRow[record[0]] = i
	}
	return hi
}

// Fonction de recherche ultra-optimisée
func (hi *HotelIndex) FindHotel(hotelID string) (map[string]string, time.Duration) {
	start := time.Now()

	result := make(map[string]string)
	row, exists := hi.idToRow[hotelID]
	if !exists {
		return nil, time.Since(start)
	}

	// Noms des champs (à adapter selon vos besoins)
	fieldNames := []string{
		"ID", "Nom", "Ville", "Pays", "Étoiles", "Adresse",
		"Email", "Latitude", "Longitude", "Téléphone",
		"Description", "Description courte",
	}

	// Décodage ciblé uniquement sur la ligne trouvée
	for col := 0; col < MAX_COLS; col++ {
		ptr := uint64(hi.matrix.matrix.At(row, col))
		result[fieldNames[col]] = hi.matrix.decodeByPointer(ptr)
	}

	return result, time.Since(start)
}*/

// Ajoutez à HotelIndex
type HotelIndex struct {
	idToRow    map[string]int // ID → Ligne
	starsIndex map[int][]int  // Étoiles → Liste de lignes
	matrix     *HotelMatrix
}

func NewHotelIndex(hm *HotelMatrix, data [][]string) *HotelIndex {
	hi := &HotelIndex{
		idToRow:    make(map[string]int),
		starsIndex: make(map[int][]int),
		matrix:     hm,
	}

	for i, record := range data {
		if i >= MAX_ROWS {
			break
		}

		// Index par ID
		hi.idToRow[record[0]] = i

		// Index par étoiles (4ème colonne, index 3 en 0-based)
		stars, _ := strconv.Atoi(record[4])
		hi.starsIndex[stars] = append(hi.starsIndex[stars], i)
	}
	return hi
}

// Recherche par étoiles (O(1) pour l'accès)
func (hi *HotelIndex) FindByStars(stars int) ([]map[string]string, time.Duration) {
	start := time.Now()

	var results []map[string]string
	rows, exists := hi.starsIndex[stars]
	if !exists {
		return nil, time.Since(start)
	}

	fieldNames := []string{"ID", "Nom", "Ville", "Pays", "Étoiles", "Adresse",
		"Email", "Latitude", "Longitude", "Téléphone",
		"Description", "Description courte"}

	for _, row := range rows {
		hotelData := make(map[string]string)
		for col := 0; col < MAX_COLS; col++ {
			ptr := uint64(hi.matrix.matrix.At(row, col))
			hotelData[fieldNames[col]] = hi.matrix.decodeByPointer(ptr)
		}
		results = append(results, hotelData)
	}

	return results, time.Since(start)
}

// --------------------------------------------------
// MAIN CORRIGÉ
// --------------------------------------------------
func main() {
	db, err := sql.Open("postgres", "host=localhost port=5433 user=myuser password=mypassword dbname=hotel_db sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	hm := NewHotelMatrix()
	data, err := fetchHotelData(db)
	if err != nil {
		log.Fatal(err)
	}

	// Encodage sécurisé
	for i, record := range data {
		if i >= MAX_ROWS {
			break
		}
		if err := hm.encodeToMatrix(i, record); err != nil {
			log.Printf("Error encoding row %d: %v", i, err)
			continue
		}
	}

	// Test de décodage
	if len(hm.dataPool) > 0 {
		ptr := uint64(hm.matrix.At(0, 0))
		decoded := hm.decodeByPointer(ptr)
		fmt.Printf("First field decoded: %s\n", decoded)
	}

	// Affichage mémoire
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Memory Usage: Alloc=%.2fMB\n", float64(m.Alloc)/1024/1024)

	// Création de l'index
	index := NewHotelIndex(hm, data)

	// Test recherche par étoiles
	stars := 2
	if hotels, duration := index.FindByStars(stars); len(hotels) > 0 {
		fmt.Printf("\nTrouvé %d hôtels avec %d étoiles en %v :\n",
			len(hotels), stars, duration)
		for i, h := range hotels {
			fmt.Printf("\nHôtel %d:\n", i+1)
			fmt.Printf("  Nom: %s\n", h["Nom"])
			fmt.Printf("  Ville: %s\n", h["Ville"])
			// Affichez d'autres champs au besoin
		}
	} else {
		fmt.Printf("Aucun hôtel trouvé avec %d étoiles\n", stars)
	}

	/*
		// Création de l'index
		index := NewHotelIndex(hm, data)

		// Test de recherche
		if result, duration := index.FindHotel("HOT317"); result != nil {
			fmt.Printf("\nRésultat trouvé en %v :\n", duration)
			for k, v := range result {
				fmt.Printf("%-15s: %.30s\n", k, v)
			}
		} else {
			fmt.Println("Hôtel non trouvé")
		}*/

}
