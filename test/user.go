package models

import (
	"context"
	"hotel-sync/cache"
)

// 5. USER (~30 lignes)
type User struct {
	ID             int     `json:"id"`
	Status         string  `json:"status"`
	Marge          float64 `json:"marge"`
	MargeOperation float64 `json:"marge_operation"`
	Solde          float64 `json:"solde"`
	SoldeRouge     float64 `json:"solde_rouge"`
	Currency       string  `json:"currency"`
	MaxRequest     int     `json:"maxrequest"`
	Group          string  `json:"group"`
	MargeB2B       float64 `json:"marge_b2b"`
	MargeXML       float64 `json:"marge_xml"`
	RoleID         int     `json:"role_id"`
}

func AddUser(hm *cache.HazelcastManager, u User) error {
	m, err := hm.Client().GetMap(context.Background(), "users")
	if err != nil {
		return err
	}
	_, err = m.Put(context.Background(), u.ID, u)
	return err
}

func GetUser(hm *cache.HazelcastManager, id int) (User, error) {
	m, err := hm.Client().GetMap(context.Background(), "users")
	if err != nil {
		return User{}, err
	}
	val, err := m.Get(context.Background(), id)
	return val.(User), err
}
