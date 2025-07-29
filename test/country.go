package models

import (
	"context"
	"hotel-sync/cache"
)

// 4. COUNTRY (~20 lignes - plus simple)
type Country struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

func (hm *cache.HazelcastManager) AddCountry(c Country) error {
	m, err := hm.client.GetMap(context.Background(), "countries")
	if err != nil {
		return err
	}
	_, err = m.Put(context.Background(), c.ID, c)
	return err
}

func (hm *cache.HazelcastManager) GetCountry(id int) (Country, error) {
	m, err := hm.client.GetMap(context.Background(), "countries")
	if err != nil {
		return Country{}, err
	}
	val, err := m.Get(context.Background(), id)
	return val.(Country), err
}
