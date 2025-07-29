package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/hazelcast/hazelcast-go-client"
	_ "github.com/lib/pq"
)

// Structure de connexion PostgreSQL
type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// [Vos structures Contract, Hotel, City, Country, User, Role et HazelcastManager restent identiques...]
// 1. CONTRACT (~30 lignes)
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

func (hm *HazelcastManager) AddContract(c Contract) error {
	m, err := hm.client.GetMap(context.Background(), "contracts")
	if err != nil {
		return err
	}
	_, err = m.Put(context.Background(), c.ID, c)
	return err
}

func (hm *HazelcastManager) GetContract(id int) (Contract, error) {
	m, err := hm.client.GetMap(context.Background(), "contracts")
	if err != nil {
		return Contract{}, err
	}
	val, err := m.Get(context.Background(), id)
	return val.(Contract), err
}

// 2. HOTEL (~30 lignes)
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

func (hm *HazelcastManager) AddHotel(h Hotel) error {
	m, err := hm.client.GetMap(context.Background(), "hotels")
	if err != nil {
		return err
	}
	_, err = m.Put(context.Background(), h.HotelKey, h)
	return err
}

func (hm *HazelcastManager) GetHotel(key string) (Hotel, error) {
	m, err := hm.client.GetMap(context.Background(), "hotels")
	if err != nil {
		return Hotel{}, err
	}
	val, err := m.Get(context.Background(), key)
	return val.(Hotel), err
}

// 3. CITY (~30 lignes)
type City struct {
	ID          int    `json:"id"`
	CityCode    string `json:"citycode"`
	Name        string `json:"name"`
	CountryName string `json:"countryname"`
	CountryID   int    `json:"countryid"`
}

func (hm *HazelcastManager) AddCity(c City) error {
	m, err := hm.client.GetMap(context.Background(), "cities")
	if err != nil {
		return err
	}
	_, err = m.Put(context.Background(), c.ID, c)
	return err
}

func (hm *HazelcastManager) GetCity(id int) (City, error) {
	m, err := hm.client.GetMap(context.Background(), "cities")
	if err != nil {
		return City{}, err
	}
	val, err := m.Get(context.Background(), id)
	return val.(City), err
}

// 4. COUNTRY (~20 lignes - plus simple)
type Country struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

func (hm *HazelcastManager) AddCountry(c Country) error {
	m, err := hm.client.GetMap(context.Background(), "countries")
	if err != nil {
		return err
	}
	_, err = m.Put(context.Background(), c.ID, c)
	return err
}

func (hm *HazelcastManager) GetCountry(id int) (Country, error) {
	m, err := hm.client.GetMap(context.Background(), "countries")
	if err != nil {
		return Country{}, err
	}
	val, err := m.Get(context.Background(), id)
	return val.(Country), err
}

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

func (hm *HazelcastManager) AddUser(u User) error {
	m, err := hm.client.GetMap(context.Background(), "users")
	if err != nil {
		return err
	}
	_, err = m.Put(context.Background(), u.ID, u)
	return err
}

func (hm *HazelcastManager) GetUser(id int) (User, error) {
	m, err := hm.client.GetMap(context.Background(), "users")
	if err != nil {
		return User{}, err
	}
	val, err := m.Get(context.Background(), id)
	return val.(User), err
}

// 6. ROLE (~15 lignes - plus simple)
type Role struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (hm *HazelcastManager) AddRole(r Role) error {
	m, err := hm.client.GetMap(context.Background(), "roles")
	if err != nil {
		return err
	}
	_, err = m.Put(context.Background(), r.ID, r)
	return err
}

func (hm *HazelcastManager) GetRole(id int) (Role, error) {
	m, err := hm.client.GetMap(context.Background(), "roles")
	if err != nil {
		return Role{}, err
	}
	val, err := m.Get(context.Background(), id)
	return val.(Role), err
}

// Manager Hazelcast (~20 lignes)
type HazelcastManager struct {
	client *hazelcast.Client
}

func NewHazelcastManager() (*HazelcastManager, error) {
	config := hazelcast.Config{}
	config.Cluster.Name = "dev-go"
	config.Cluster.Network.SetAddresses("localhost:5701")
	client, err := hazelcast.StartNewClientWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}
	return &HazelcastManager{client: client}, nil
}

func (hm *HazelcastManager) Close() {
	if hm.client != nil {
		hm.client.Shutdown(context.Background())
	}
}

func connectPostgres(cfg PostgresConfig) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func importFromPostgresToHazelcast(db *sql.DB, hm *HazelcastManager) error {
	// 1. Importer les pays (COUNTRY)
	rows, err := db.Query("SELECT id, name, code FROM COUNTRY")
	if err != nil {
		return fmt.Errorf("erreur COUNTRY: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var c Country
		if err := rows.Scan(&c.ID, &c.Name, &c.Code); err != nil {
			return err
		}
		if err := hm.AddCountry(c); err != nil {
			return err
		}
	}

	// 2. Importer les villes (CITY)
	rows, err = db.Query("SELECT id, citycode, name, countryname, countryid FROM CITY")
	if err != nil {
		return fmt.Errorf("erreur CITY: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var c City
		if err := rows.Scan(&c.ID, &c.CityCode, &c.Name, &c.CountryName, &c.CountryID); err != nil {
			return err
		}
		if err := hm.AddCity(c); err != nil {
			return err
		}
	}

	// 3. Importer les hôtels (HOTEL)
	rows, err = db.Query("SELECT hotel_key, name, city, country, stars, address, mail, latitude, longitude, phone, description, short_description FROM HOTEL")
	if err != nil {
		return fmt.Errorf("erreur HOTEL: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var h Hotel
		if err := rows.Scan(&h.HotelKey, &h.Name, &h.City, &h.Country, &h.Stars, &h.Address, &h.Mail, &h.Latitude, &h.Longitude, &h.Phone, &h.Description, &h.ShortDescription); err != nil {
			return err
		}
		if err := hm.AddHotel(h); err != nil {
			return err
		}
	}

	// 4. Importer les rôles (ROLE)
	rows, err = db.Query("SELECT id, name FROM ROLE")
	if err != nil {
		return fmt.Errorf("erreur ROLE: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var r Role
		if err := rows.Scan(&r.ID, &r.Name); err != nil {
			return err
		}
		if err := hm.AddRole(r); err != nil {
			return err
		}
	}

	// 5. Importer les utilisateurs (USER_APP)
	rows, err = db.Query(`SELECT id, status, marge, marge_operation, solde, solde_rouge, currency, maxrequest, "group", marge_b2b, marge_xml, role_id FROM USER_APP`)
	if err != nil {
		return fmt.Errorf("erreur USER_APP: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Status, &u.Marge, &u.MargeOperation, &u.Solde, &u.SoldeRouge, &u.Currency, &u.MaxRequest, &u.Group, &u.MargeB2B, &u.MargeXML, &u.RoleID); err != nil {
			return err
		}
		if err := hm.AddUser(u); err != nil {
			return err
		}
	}

	// 6. Importer les contrats (CONTRACT)
	rows, err = db.Query("SELECT id, name, hotel_id, start_at, end_at, access, active, currency, market, client_id FROM CONTRACT")
	if err != nil {
		return fmt.Errorf("erreur CONTRACT: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var c Contract
		var startAt, endAt time.Time
		if err := rows.Scan(&c.ID, &c.Name, &c.HotelID, &startAt, &endAt, &c.Access, &c.Active, &c.Currency, &c.Market, &c.ClientID); err != nil {
			return err
		}
		c.StartAt = startAt
		c.EndAt = endAt
		if err := hm.AddContract(c); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	// Configuration PostgreSQL
	pgConfig := PostgresConfig{
		Host:     "localhost",
		Port:     "5433",
		User:     "myuser",
		Password: "mypassword",
		DBName:   "hotel_db",
	}

	// Connexion PostgreSQL
	db, err := connectPostgres(pgConfig)
	if err != nil {
		log.Fatalf("Erreur PostgreSQL: %v", err)
	}
	defer db.Close()

	// Connexion Hazelcast
	hm, err := NewHazelcastManager()
	if err != nil {
		log.Fatalf("Erreur Hazelcast: %v", err)
	}
	defer hm.Close()

	// Import des données
	if err := importFromPostgresToHazelcast(db, hm); err != nil {
		log.Fatalf("Erreur import: %v", err)
	}

	log.Println("Import terminé avec succès!")
}
