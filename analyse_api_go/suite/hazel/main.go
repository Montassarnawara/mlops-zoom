package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	hz "github.com/hazelcast/hazelcast-go-client"
	_ "github.com/lib/pq"
)

// Contexte & client Hazelcast global
var (
	Ctx      = context.Background()
	Client   *hz.Client
	CacheMap *hz.Map
)

// 🎓 Modèle College
type College struct {
	CollegeID      string
	IQ             int
	CGPA           float64
	MinMarks       float64
	Attendance     int
	Backlogs       string
	Communication  int
	Logical        int
	Projects       int
	Certifications string
}

// Initialiser la connexion Hazelcast
func InitHazelcast() {
	config := hz.Config{}
	config.Cluster.Network.SetAddresses("127.0.0.1:5701")
	config.Cluster.Name = "dev"
	config.Security.Credentials.Username = "dev"
	config.Security.Credentials.Password = "dev-pass"

	var err error
	Client, err = hz.StartNewClientWithConfig(Ctx, config)
	if err != nil {
		log.Fatalf(" Connexion Hazelcast échouée : %v", err)
	}

	CacheMap, err = Client.GetMap(Ctx, "college_cache")
	if err != nil {
		log.Fatalf(" Impossible d'accéder à la map Hazelcast : %v", err)
	}

	log.Println(" Connexion Hazelcast établie")
}

// Connexion à PostgreSQL
func FetchCollegesFromPostgres() []College {
	connStr := "host=localhost port=5432 user=montassar password=123mont@456 dbname=ma_base sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf(" Erreur connexion PostgreSQL : %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf(" PostgreSQL indisponible : %v", err)
	}
	fmt.Println(" Connexion PostgreSQL établie")

	query := `SELECT * FROM college LIMIT 100`
	rows, err := db.Query(query)
	if err != nil {
		log.Fatalf(" Erreur requête SQL : %v", err)
	}
	defer rows.Close()

	var colleges []College
	for rows.Next() {
		var c College
		err := rows.Scan(
			&c.CollegeID, &c.IQ, &c.CGPA, &c.MinMarks, &c.Attendance,
			&c.Backlogs, &c.Communication, &c.Logical, &c.Projects, &c.Certifications,
		)
		if err != nil {
			log.Printf(" Erreur lecture ligne : %v", err)
			continue
		}
		colleges = append(colleges, c)
	}
	log.Printf(" %d colleges récupérés depuis PostgreSQL", len(colleges))
	return colleges
}

// Insertion dans Hazelcast
func PopulateHazelcast(colleges []College) {
	count := 0
	for _, c := range colleges {
		err := CacheMap.Set(Ctx, c.CollegeID, c)
		if err != nil {
			log.Printf(" Erreur insertion pour %s : %v", c.CollegeID, err)
			continue
		}
		count++
	}
	log.Printf("%d colleges insérés dans Hazelcast", count)
}

func main() {
	InitHazelcast()
	defer Client.Shutdown(Ctx)

	colleges := FetchCollegesFromPostgres()
	PopulateHazelcast(colleges)

	log.Println("Initialisation terminée.")
}
