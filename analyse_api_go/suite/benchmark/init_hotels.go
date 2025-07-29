package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	hz "github.com/hazelcast/hazelcast-go-client"
)

type Hotel struct {
	IDHotel      int `json:"id_hotel"`
	NBRoom       int `json:"nb_room"`
	TypeRoom     int `json:"type_room"`
	OfferPerRoom int `json:"offer_per_room"`
	Saison       int `json:"saison"`
}

func randomHotel(id int) Hotel {
	return Hotel{
		IDHotel:      id,
		NBRoom:       rand.Intn(200) + 20,
		TypeRoom:     rand.Intn(4) + 1,
		OfferPerRoom: rand.Intn(10) + 1,
		Saison:       rand.Intn(3) + 1,
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	ctx := context.Background()

	config := hz.Config{}
	config.Cluster.Name = "dev-go"
	client, err := hz.StartNewClientWithConfig(ctx, config)
	if err != nil {
		panic(err)
	}
	defer client.Shutdown(ctx)

	hotelMap, err := client.GetMap(ctx, "hotels")
	if err != nil {
		panic(err)
	}

	hotelMap.Clear(ctx)

	for i := 1; i <= 100; i++ {
		hotel := randomHotel(i)
		jsonHotel, err := json.Marshal(hotel)
		if err != nil {
			fmt.Println("Erreur JSON:", err)
			continue
		}
		if err := hotelMap.Set(ctx, i, jsonHotel); err != nil {
			fmt.Println("Erreur d'insertion Hazelcast:", err)
		}
	}

	fmt.Println("✅ Hôtels insérés dans Hazelcast au format JSON.")
}
