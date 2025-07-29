package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	hz "github.com/hazelcast/hazelcast-go-client"
	_ "github.com/lib/pq"
)

type College struct {
	CollegeID            string  `json:"College_ID"`
	IQ                   int     `json:"IQ"`
	PrevSemResult        float64 `json:"Prev_Sem_Result"`
	CGPA                 float64 `json:"CGPA"`
	AcademicPerformance  int     `json:"Academic_Performance"`
	InternshipExperience string  `json:"Internship_Experience"`
	ExtraCurricularScore int     `json:"Extra_Curricular_Score"`
	CommunicationSkills  int     `json:"Communication_Skills"`
	ProjectsCompleted    int     `json:"Projects_Completed"`
	Placement            string  `json:"Placement"`
}

func main() {
	// Connexion PostgreSQL
	connStr := "host=localhost port=5432 user=montassar password=123mont@456 dbname=ma_base sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Erreur connexion PostgreSQL : %v", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalf("Erreur ping PostgreSQL : %v", err)
	}
	fmt.Println("✅ Connexion PostgreSQL OK")

	// Connexion Hazelcast
	config := hz.Config{}
	config.Cluster.Name = "dev-go"
	config.Cluster.Network.SetAddresses("127.0.0.1:5701")
	client, err := hz.StartNewClientWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Erreur connexion Hazelcast : %v", err)
	}
	defer client.Shutdown(context.Background())
	cacheMap, err := client.GetMap(context.Background(), "college_cache")
	if err != nil {
		log.Fatalf("Erreur accès map Hazelcast : %v", err)
	}
	fmt.Println("✅ Connexion Hazelcast OK")

	// Efface la map Hazelcast avant import (sécurité)
	err = cacheMap.Clear(context.Background())
	if err != nil {
		log.Fatalf("Erreur lors du clear de la map Hazelcast : %v", err)
	}
	fmt.Println("🧹 Map Hazelcast vidée avant import.")

	// Lecture des colleges depuis PostgreSQL
	rows, err := db.Query(`SELECT college_id, iq, prev_sem_result, cgpa, academic_performance, internship_experience, extra_curricular_score, communication_skills, projects_completed, placement FROM college`)
	if err != nil {
		log.Fatalf("Erreur requête SELECT : %v", err)
	}
	defer rows.Close()

	nb := 0
	for rows.Next() {
		var c College
		err := rows.Scan(
			&c.CollegeID, &c.IQ, &c.PrevSemResult, &c.CGPA, &c.AcademicPerformance,
			&c.InternshipExperience, &c.ExtraCurricularScore, &c.CommunicationSkills,
			&c.ProjectsCompleted, &c.Placement,
		)
		if err != nil {
			log.Printf("Erreur scan : %v", err)
			continue
		}
		jsonBytes, err := json.Marshal(c)
		if err != nil {
			log.Printf("Erreur JSON : %v", err)
			continue
		}
		key := fmt.Sprintf("CLG%04d", nb+1) // CLG0001, CLG0002, ...
		if _, err := cacheMap.Put(context.Background(), key, jsonBytes); err != nil {
			log.Printf("Erreur insertion Hazelcast pour %s : %v", key, err)
			continue
		}
		nb++
	}
	fmt.Printf("✅ %d colleges insérés dans Hazelcast sous les clés CLGxxxx.\n", nb)
}
