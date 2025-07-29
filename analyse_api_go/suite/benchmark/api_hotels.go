package main

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	hz "github.com/hazelcast/hazelcast-go-client"
)

type Hotel struct {
	IDHotel      int `json:"id_hotel"`
	NBRoom       int `json:"nb_room"`
	TypeRoom     int `json:"type_room"`
	OfferPerRoom int `json:"offer_per_room"`
	Saison       int `json:"saison"`
}

func initHazelcastClient() (*hz.Client, *hz.Map, context.Context) {
	ctx := context.Background()
	config := hz.Config{}
	config.Cluster.Name = "dev-go"

	client, err := hz.StartNewClientWithConfig(ctx, config)
	if err != nil {
		panic(err)
	}

	hotelMap, err := client.GetMap(ctx, "hotels")
	if err != nil {
		panic(err)
	}

	return client, hotelMap, ctx
}

func getAllHotelsHandler(ctx context.Context, hotelMap *hz.Map) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now() // Début mesure
		keys, err := hotelMap.GetKeySet(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur Hazelcast"})
			return
		}

		hotels := make([]Hotel, 0, len(keys))
		for _, key := range keys {
			raw, err := hotelMap.Get(ctx, key)
			if err != nil || raw == nil {
				continue
			}
			jsonHotel, ok := raw.([]byte)
			if !ok {
				continue
			}
			var hotel Hotel

			if err := json.Unmarshal(jsonHotel, &hotel); err == nil {
				hotels = append(hotels, hotel)
			}
			// Enregistrement de la durée dans un tableau ou un système de monitoring
		}
		duration := time.Since(start).Microseconds() // Fin mesure
		c.JSON(http.StatusOK, gin.H{
			"hotels":                hotels,
			"duration_microseconds": duration,
		})
	}
}

func getHotelByIDHandler(ctx context.Context, hotelMap *hz.Map) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
			return
		}

		raw, err := hotelMap.Get(ctx, id)
		if err != nil || raw == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Hôtel non trouvé"})
			return
		}

		jsonHotel, ok := raw.([]byte)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur de type"})
			return
		}

		var hotel Hotel
		if err := json.Unmarshal(jsonHotel, &hotel); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur JSON"})
			return
		}
		c.JSON(http.StatusOK, hotel)
	}
}

// === NOUVEAU HANDLER POUR LA SOMME DES ROOMS ===
func getSumRoomsHandler(ctx context.Context, hotelMap *hz.Map) gin.HandlerFunc {
	return func(c *gin.Context) {
		n, err := strconv.Atoi(c.Param("n"))
		if err != nil || n <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Paramètre invalide"})
			return
		}

		totalRooms := 0
		for i := 1; i <= n; i++ {
			raw, err := hotelMap.Get(ctx, i)
			if err != nil || raw == nil {
				continue
			}
			jsonHotel, ok := raw.([]byte)
			if !ok {
				continue
			}
			var hotel Hotel
			if err := json.Unmarshal(jsonHotel, &hotel); err == nil {
				totalRooms += hotel.NBRoom
			}
		}

		c.JSON(http.StatusOK, gin.H{"nb_total_room": totalRooms})
	}
}

func getHotelConditionHandler(ctx context.Context, hotelMap *hz.Map) gin.HandlerFunc {
	return func(c *gin.Context) {
		saison, err1 := strconv.Atoi(c.Param("saison"))
		choix, err2 := strconv.Atoi(c.Param("choix"))
		if err1 != nil || err2 != nil || saison < 1 || saison > 4 || choix < 1 || choix > 3 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Paramètres invalides"})
			return
		}

		// Règles par saison (map de slice)
		saisonTypes := map[int][]int{
			1: {1, 2},
			2: {2, 3},
			3: {1, 3},
			4: {1, 2, 3},
		}

		typesAcceptes := saisonTypes[saison]
		match := func(t int) bool {
			for _, v := range typesAcceptes {
				if v == t {
					return true
				}
			}
			return false
		}

		keys, err := hotelMap.GetKeySet(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur Hazelcast"})
			return
		}

		result := []int{}
		for _, key := range keys {
			raw, err := hotelMap.Get(ctx, key)
			if err != nil || raw == nil {
				continue
			}
			jsonHotel, ok := raw.([]byte)
			if !ok {
				continue
			}
			var hotel Hotel
			if err := json.Unmarshal(jsonHotel, &hotel); err == nil {
				if hotel.Saison == saison && hotel.TypeRoom == choix && match(hotel.TypeRoom) {
					result = append(result, hotel.IDHotel)
				}
			}
		}

		c.JSON(http.StatusOK, gin.H{"id_hotel_cond": result})
	}
}

func main() {
	rand.Seed(time.Now().UnixNano()) // Initialisation aléatoire
	client, hotelMap, ctx := initHazelcastClient()
	defer client.Shutdown(ctx)

	r := gin.Default()
	r.GET("/hotels", getAllHotelsHandler(ctx, hotelMap))
	r.GET("/hotels/:id", getHotelByIDHandler(ctx, hotelMap))
	r.GET("/nbroom/:n", getSumRoomsHandler(ctx, hotelMap)) // NOUVEL ENDPOINT
	r.GET("/condition/:saison/:choix", getHotelConditionHandler(ctx, hotelMap))

	r.Run(":8080")
}

// Pour corriger l'erreur d'import Gin :
// 1. Ouvre un terminal dans /home/montassar/Bureau/go/benchmark
// 2. Lance : go mod init benchmark
// 3. Puis : go get github.com/gin-gonic/gin
// 4. Compile ensuite normalement : go run api_hotels.go

// Exemple d'utilisation de l'API avec curl :
//
// Pour obtenir la liste de tous les hôtels :
// curl http://localhost:8080/hotels
//
// Pour obtenir un hôtel par son ID (exemple pour l'hôtel d'ID 5) :
// curl http://localhost:8080/hotels/5
// 2. Lance : go mod init benchmark
// 3. Puis : go get github.com/gin-gonic/gin
// 4. Compile ensuite normalement : go run api_hotels.go

// Exemple d'utilisation de l'API avec curl :
//
// Pour obtenir la liste de tous les hôtels :
// curl http://localhost:8080/hotels
//
// Pour obtenir un hôtel par son ID (exemple pour l'hôtel d'ID 5) :
// curl http://localhost:8080/hotels/5
// Pour obtenir la somme des chambres pour les N premiers hôtels (exemple pour N=10) :
// curl http://localhost:8080/nbroom/10
// Pour obtenir les hôtels selon la saison et le type de chambre (exemple pour saison 1 et type 2) :
// curl http://localhost:8080/condition/1/2
