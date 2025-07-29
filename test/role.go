package models

import (
	"context"
	"hotel-sync/cache"
)

// 6. ROLE (~15 lignes - plus simple)
type Role struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (hm *cache.HazelcastManager) AddRole(r Role) error {
	m, err := hm.client.GetMap(context.Background(), "roles")
	if err != nil {
		return err
	}
	_, err = m.Put(context.Background(), r.ID, r)
	return err
}

func (hm *cache.HazelcastManager) GetRole(id int) (Role, error) {
	m, err := hm.client.GetMap(context.Background(), "roles")
	if err != nil {
		return Role{}, err
	}
	val, err := m.Get(context.Background(), id)
	return val.(Role), err
}
