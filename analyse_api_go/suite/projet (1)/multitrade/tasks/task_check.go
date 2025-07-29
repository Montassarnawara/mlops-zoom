package tasks

import (
	"context"
	"fmt"
	"go/mutitrade/cache"
	"go/mutitrade/models"
	"log"
	"time"
)

func TaskCheckCondition(n int) (int, time.Duration) {
	start := time.Now()
	ctx := context.Background()

	colMap, err := cache.HazelcastClient.GetMap(ctx, "colleges")
	if err != nil {
		log.Fatalf("❌ Erreur récupération map Hazelcast : %v", err)
	}

	count := 0

	for i := 1; i <= n; i++ {
		id := fmt.Sprintf("CLG%04d", i)
		val, err := colMap.Get(ctx, id)
		if err != nil {
			log.Printf("❌ Erreur lecture %s : %v", id, err)
			continue
		}

		col, ok := val.(models.College)
		if !ok {
			log.Printf("❌ Donnée invalide %s", id)
			continue
		}

		if col.Placement == "Yes" || col.IQ > 100 {
			count++
		}
	}

	duration := time.Since(start)
	return count, duration
}
