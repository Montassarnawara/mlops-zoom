// hotel-sync/cache/hazelcast.go
package cache

import (
	"context"

	"github.com/hazelcast/hazelcast-go-client"
)

type HazelcastManager struct {
	Client *hazelcast.Client
}

func NewHazelcastManager() (*HazelcastManager, error) {
	cfg := hazelcast.Config{}

	// ⚠️ Fix cluster name ici
	cfg.Cluster.Name = "dev-go"

	// Adresse du conteneur Hazelcast
	cfg.Cluster.Network.SetAddresses("localhost:5701")

	client, err := hazelcast.StartNewClientWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, err
	}
	return &HazelcastManager{Client: client}, nil
}

func (hm *HazelcastManager) Close() {
	if hm.Client != nil {
		hm.Client.Shutdown(context.Background())
	}
}
