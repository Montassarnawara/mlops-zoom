package models

import (
	"context"
	"hotel-sync/cache"
	"time"
)

type Contract struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	HotelID  string    `json:"hotel_id"`
	StartAt  time.Time `json:"start_at"`
	EndAt    time.Time `json:"end_at"`
	Access   string    `json:"access"`
	Active   bool      `json:"active"`
	Currency string    `json:"currency"`
	Market   int       `json:"market"`
	ClientID int       `json:"client_id"`
}

func (hm *cache.HazelcastManager) AddContract(c Contract) error {
	m, err := hm.client.GetMap(context.Background(), "contracts")
	if err != nil {
		return err
	}
	_, err = m.Put(context.Background(), c.ID, c)
	return err
}

func (hm *cache.HazelcastManager) GetContract(id int) (Contract, error) {
	m, err := hm.client.GetMap(context.Background(), "contracts")
	if err != nil {
		return Contract{}, err
	}
	val, err := m.Get(context.Background(), id)
	return val.(Contract), err
}
