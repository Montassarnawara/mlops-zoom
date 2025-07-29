package db

// Package db provides functions to interact with the PostgreSQL database.
// It includes initialization of the database connection and retrieval of college data.

import (
	"database/sql"
	"fmt"
	"log"

	"go-hazelcast-api/models"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() error {
	connStr := "host=localhost port=5432 user=montassar password=123mont@456 dbname=ma_base sslmode=disable"
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	err = DB.Ping()
	if err != nil {
		return err
	}
	fmt.Println(" Connexion PostgreSQL établie")
	return nil
}

func GetFirst100Colleges() ([]models.College, error) {
	query := `SELECT * FROM college`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var colleges []models.College
	for rows.Next() {
		var c models.College
		err := rows.Scan(
			&c.CollegeID, &c.IQ, &c.CGPA, &c.MinMarks, &c.Attendance,
			&c.Backlogs, &c.Communication, &c.Logical, &c.Projects,
			&c.Certifications,
		)
		if err != nil {
			log.Println("Erreur de scan :", err)
			continue
		}
		colleges = append(colleges, c)
	}
	return colleges, nil
}
