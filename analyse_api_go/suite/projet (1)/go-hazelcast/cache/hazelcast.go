package cache

import (
	"context"
	"encoding/json"
	"log"

	"go/hazelcast/api/models"

	"github.com/hazelcast/hazelcast-go-client"
)

var (
	ctx      = context.Background()
	Client   *hazelcast.Client
	CacheMap *hazelcast.Map
)

// Initialiser la connexion Hazelcast
func InitHazelcast() {
	config := hazelcast.Config{}
	config.Cluster.Network.SetAddresses("127.0.0.1:5701")
	config.Cluster.Name = "dev-go"

	var err error
	Client, err = hazelcast.StartNewClientWithConfig(ctx, config)
	if err != nil {
		log.Fatalf(" Connexion Hazelcast échouée : %v", err)
	}

	CacheMap, err = Client.GetMap(ctx, "college")
	if err != nil {
		log.Fatalf(" Impossible d'accéder à la map Hazelcast : %v", err)
	}

	log.Println("Connexion Hazelcast établie")
}

// Lire un élément du cache par ID
func GetFromCache(id string) (*models.College, error) {
	raw, err := CacheMap.Get(ctx, id)
	if err != nil || raw == nil {
		return nil, err // nil si pas trouvé
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

// Enregistrer un élément dans le cache
func SetToCache(college *models.College) error {
	jsonBytes, err := json.Marshal(college)
	if err != nil {
		return err
	}

	return CacheMap.Set(ctx, college.CollegeID, jsonBytes)
}
