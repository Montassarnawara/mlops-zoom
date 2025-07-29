package models

import (
	"context"
	"hotel-sync/cache"
)

type City struct {
	ID          int    `json:"id"`
	CityCode    string `json:"citycode"`
	Name        string `json:"name"`
	CountryName string `json:"countryname"`
	CountryID   int    `json:"countryid"`
}

func (hm *cache.HazelcastManager) AddCity(c City) error {
	m, err := hm.client.GetMap(context.Background(), "cities")
	if err != nil {
		return err
	}
	_, err = m.Put(context.Background(), c.ID, c)
	return err
}

func (hm *cache.HazelcastManager) GetCity(id int) (City, error) {
	m, err := hm.client.GetMap(context.Background(), "cities")
	if err != nil {
		return City{}, err
	}
	val, err := m.Get(context.Background(), id)
	return val.(City), err
}
