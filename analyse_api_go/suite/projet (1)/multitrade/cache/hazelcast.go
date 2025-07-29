package cache

import (
	"context"
	"encoding/json"
	"go/mutitrade/models"
	"log"

	hz "github.com/hazelcast/hazelcast-go-client"
)

// GetCollegesFromCacheBulk récupère plusieurs colleges d'un coup via Hazelcast GetAll
func GetCollegesFromCacheBulk(keys []string) ([]*models.College, error) {
	colleges := make([]*models.College, 0, len(keys))
	keyInterfaces := make([]interface{}, len(keys))
	for i, k := range keys {
		keyInterfaces[i] = k
	}
	entries, err := CacheMap.GetAll(context.Background(), keyInterfaces...)
	if err != nil {
		return nil, err
	}
	// Map clé -> valeur pour accès rapide
	valueMap := make(map[string][]byte)
	for _, entry := range entries {
		k, ok1 := entry.Key.(string)
		v, ok2 := entry.Value.([]byte)
		if ok1 && ok2 {
			valueMap[k] = v
		}
	}
	for _, k := range keys {
		v, ok := valueMap[k]
		if !ok {
			colleges = append(colleges, nil)
			continue
		}
		var college models.College
		if err := json.Unmarshal(v, &college); err != nil {
			colleges = append(colleges, nil)
			continue
		}
		colleges = append(colleges, &college)
	}
	return colleges, nil
}

// SanityCheckColleges vérifie la présence des N premiers colleges dans Hazelcast

var (
	HazelcastClient *hz.Client
	CacheMap        *hz.Map
)

func InitHazelcast() *hz.Client {
	config := hz.Config{}
	config.Cluster.Name = "dev-go"

	client, err := hz.StartNewClientWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Erreur de connexion Hazelcast : %v", err)
	}

	HazelcastClient = client

	// 🟠 Important : on utilise la même map que le projet 1
	CacheMap, err = HazelcastClient.GetMap(context.Background(), "college_cache")
	if err != nil {
		log.Fatalf("Erreur lors de la récupération de la map Hazelcast : %v", err)
	}

	log.Println("✅ Hazelcast connecté avec succès.")
	return client
}

// GetCollegeFromCache récupère un college en JSON et le désérialise
func GetCollegeFromCache(id string) (*models.College, error) {
	raw, err := CacheMap.Get(context.Background(), id)
	if err != nil || raw == nil {
		return nil, err
	}

	jsonBytes, ok := raw.([]byte)
	if !ok {
		return nil, nil
	}

	var college models.College
	if err := json.Unmarshal(jsonBytes, &college); err != nil {
		return nil, err
	}

	return &college, nil
}
