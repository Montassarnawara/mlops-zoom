package models

import "time"

// Country corresponds to your real data structure
type Country struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

// City corresponds to your real data structure  
type City struct {
	ID          int    `json:"id"`
	CityCode    string `json:"citycode"`
	Name        string `json:"name"`
	CountryName string `json:"countryname"`
	CountryID   int    `json:"countryid"`
}

// Hotel corresponds to your real data structure
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

// Role corresponds to your real data structure
type Role struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// User corresponds to your real data structure
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

// Contract corresponds to your real data structure
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
