package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"hotel-sync/cache"
	"hotel-sync/models" // Ensure 'models' is a valid package and contains the required definitions.
	// If the 'models' package does not exist, create it and define the required structs like Country, City, Hotel, Role, User, and Contract.
)

type Importer struct {
	db *sql.DB
	hz *cache.HazelcastManager
}

func NewImporter(db *sql.DB, hz *cache.HazelcastManager) *Importer {
	return &Importer{db: db, hz: hz}
}

func (imp *Importer) ImportAll() error {
	ctx := context.Background()

	// --- COUNTRY ---
	rows, err := imp.db.Query("SELECT id, name, code FROM COUNTRY")
	if err != nil {
		return fmt.Errorf("COUNTRY query error: %w", err)
	}
	defer rows.Close()

	mCountries, err := imp.hz.Client.GetMap(ctx, "countries")
	if err != nil {
		return err
	}

	for rows.Next() {
		var c models.Country // Ensure 'Country' is defined in the 'models' package.
		if err := rows.Scan(&c.ID, &c.Name, &c.Code); err != nil {
			return err
		}
		if _, err := mCountries.Put(ctx, c.ID, c); err != nil {
			return err
		}
	}
	log.Println("Import COUNTRY terminé")

	// --- CITY ---
	rows, err = imp.db.Query("SELECT id, citycode, name, countryname, countryid FROM CITY")
	if err != nil {
		return fmt.Errorf("CITY query error: %w", err)
	}
	defer rows.Close()

	mCities, err := imp.hz.Client.GetMap(ctx, "cities")
	if err != nil {
		return err
	}

	for rows.Next() {
		var c models.City // Ensure 'City' is defined in the 'models' package.
		if err := rows.Scan(&c.ID, &c.CityCode, &c.Name, &c.CountryName, &c.CountryID); err != nil {
			return err
		}
		if _, err := mCities.Put(ctx, c.ID, c); err != nil {
			return err
		}
	}
	log.Println("Import CITY terminé")

	// --- HOTEL ---
	rows, err = imp.db.Query("SELECT hotel_key, name, city, country, stars, address, mail, latitude, longitude, phone, description, short_description FROM HOTEL")
	if err != nil {
		return fmt.Errorf("HOTEL query error: %w", err)
	}
	defer rows.Close()

	mHotels, err := imp.hz.Client.GetMap(ctx, "hotels")
	if err != nil {
		return err
	}

	for rows.Next() {
		var h models.Hotel // Ensure 'Hotel' is defined in the 'models' package.
		if err := rows.Scan(&h.HotelKey, &h.Name, &h.City, &h.Country, &h.Stars, &h.Address, &h.Mail, &h.Latitude, &h.Longitude, &h.Phone, &h.Description, &h.ShortDescription); err != nil {
			return err
		}
		if _, err := mHotels.Put(ctx, h.HotelKey, h); err != nil {
			return err
		}
	}
	log.Println("Import HOTEL terminé")

	// --- ROLE ---
	rows, err = imp.db.Query("SELECT id, name FROM ROLE")
	if err != nil {
		return fmt.Errorf("ROLE query error: %w", err)
	}
	defer rows.Close()

	mRoles, err := imp.hz.Client.GetMap(ctx, "roles")
	if err != nil {
		return err
	}

	for rows.Next() {
		var r models.Role // Ensure 'Role' is defined in the 'models' package.
		if err := rows.Scan(&r.ID, &r.Name); err != nil {
			return err
		}
		if _, err := mRoles.Put(ctx, r.ID, r); err != nil {
			return err
		}
	}
	log.Println("Import ROLE terminé")

	// --- USER ---
	rows, err = imp.db.Query(`SELECT id, status, marge, marge_operation, solde, solde_rouge, currency, maxrequest, "group", marge_b2b, marge_xml, role_id FROM USER_APP`)
	if err != nil {
		return fmt.Errorf("USER_APP query error: %w", err)
	}
	defer rows.Close()

	mUsers, err := imp.hz.Client.GetMap(ctx, "users")
	if err != nil {
		return err
	}

	for rows.Next() {
		var u models.User // Ensure 'User' is defined in the 'models' package.
		if err := rows.Scan(&u.ID, &u.Status, &u.Marge, &u.MargeOperation, &u.Solde, &u.SoldeRouge, &u.Currency, &u.MaxRequest, &u.Group, &u.MargeB2B, &u.MargeXML, &u.RoleID); err != nil {
			return err
		}
		if _, err := mUsers.Put(ctx, u.ID, u); err != nil {
			return err
		}
	}
	log.Println("Import USER terminé")

	// --- CONTRACT ---
	rows, err = imp.db.Query("SELECT id, name, hotel_id, start_at, end_at, access, active, currency, market, client_id FROM contract")
	if err != nil {
		return fmt.Errorf("contract query error: %w", err)
	}
	defer rows.Close()

	mContracts, err := imp.hz.Client.GetMap(ctx, "contracts")
	if err != nil {
		return err
	}

	for rows.Next() {
		var c models.Contract // Ensure 'Contract' is defined in the 'models' package.
		var startAt, endAt time.Time
		if err := rows.Scan(&c.ID, &c.Name, &c.HotelID, &startAt, &endAt, &c.Access, &c.Active, &c.Currency, &c.Market, &c.ClientID); err != nil {
			return err
		}
		c.StartAt = startAt
		c.EndAt = endAt
		if _, err := mContracts.Put(ctx, c.ID, c); err != nil {
			return err
		}
	}
	log.Println("Import CONTRACT terminé")

	return nil
}
