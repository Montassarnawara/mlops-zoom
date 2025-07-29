package main

import (
	"fmt"
	"strings"
	"time"
)

func encode(s string) []int {
	encoded := make([]int, 0, len(s))
	for _, c := range s {
		encoded = append(encoded, int(c))
	}
	return encoded
}

func decode(encoded []int) string {
	var builder strings.Builder
	builder.Grow(len(encoded))
	for _, v := range encoded {
		builder.WriteRune(rune(v))
	}
	return builder.String()
}

func processHotelRecord(record string) ([][]int, []string) {
	start := time.Now()
	fields := strings.Split(record, ",")

	// Correction: Ajustement dynamique des noms de champs
	fieldNames := []string{
		"ID", "Nom", "Ville", "Pays", "Étoiles", "Adresse",
		"Email", "Latitude", "Longitude", "Téléphone",
		"Description", "Description courte",
	}

	// S'assurer qu'on a assez de noms de champs
	if len(fields) > len(fieldNames) {
		// Générer des noms supplémentaires si nécessaire
		for i := len(fieldNames); i < len(fields); i++ {
			fieldNames = append(fieldNames, fmt.Sprintf("Champ %d", i+1))
		}
	}

	vector := make([][]int, len(fields))

	for i, field := range fields {
		field = strings.TrimSpace(field)
		encodeStart := time.Now()
		vector[i] = encode(field)
		fmt.Printf("[%s] Encodé en %v\n", fieldNames[i], time.Since(encodeStart))
	}

	fmt.Printf("\nTemps total d'encodage : %v\n", time.Since(start))
	return vector, fieldNames[:len(fields)] // Retourne seulement les noms utilisés
}

func displayEncodedVector(vector [][]int, fieldNames []string) {
	fmt.Println("\n=== Vecteur Encodé Complet ===")
	for i, v := range vector {
		fmt.Printf("%-18s: %v\n", fieldNames[i], v)
	}
}

func main() {
	hotelRecord := `HOT061,Lane-Short,North Elizabethburgh,Saint Barthelemy,3,"42591 Boyd Streets Apt. 987 Michaelmouth, MH 73574",davidbeard@fields-morris.com,22.761561,33.697786,475.244.2084,Write movie collection process authority news democratic. Recently especially medical cover performance read respond.,View quality use society real require investment.`

	// Encodage
	fmt.Println("=== Encodage ===")
	vector, fieldNames := processHotelRecord(hotelRecord)
	displayEncodedVector(vector, fieldNames)

	// Décodage
	fmt.Println("\n=== Décodage ===")
	for i, field := range vector {
		start := time.Now()
		decoded := decode(field)
		fmt.Printf("[%-18s] Décodé en %-10v: %q\n",
			fieldNames[i], time.Since(start),
			truncateString(decoded, 30))
	}
}

func truncateString(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
