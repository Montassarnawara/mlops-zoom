package models

import (
	"context"
	"hotel-sync/cache"
)

type Hotel struct {
	HotelKey         string  `json:"hotel_key"`
	Name             string  `json:"name"`
	City             string  `json:"city"`
	Country          string  `json:"country"`
	Stars            int     `json:"stars"`
	Address          string  `json:"address"`
	Mail             string  `json:"mail"`
	Latitude         float64 `json:"latitude"`
	Longitude        float64 `json:"longitude"`
	Phone            string  `json:"phone"`
	Description      string  `json:"description"`
	ShortDescription string  `json:"short_description"`
}

func (hm *cache.HazelcastManager) AddHotel(h Hotel) error {
	m, err := hm.client.GetMap(context.Background(), "hotels")
	if err != nil {
		return err
	}
	_, err = m.Put(context.Background(), h.HotelKey, h)
	return err
}

func (hm *cache.HazelcastManager) GetHotel(key string) (Hotel, error) {
	m, err := hm.client.GetMap(context.Background(), "hotels")
	if err != nil {
		return Hotel{}, err
	}
	val, err := m.Get(context.Background(), key)
	return val.(Hotel), err
}
