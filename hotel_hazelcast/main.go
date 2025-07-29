package main

import (
	"encoding/gob"
	"hotel-sync/cache"
	"hotel-sync/config"
	"hotel-sync/models"
	"hotel-sync/services"
	"log"
)

func init() {
	// Enregistrer tous les types utilisés dans les maps Hazelcast
	gob.Register(models.Country{})
	gob.Register(models.City{})
	gob.Register(models.Hotel{})
	gob.Register(models.Role{})
	gob.Register(models.User{})
	gob.Register(models.Contract{})
}

func main() {
	// Connexion PostgreSQL
	db, err := config.ConnectPostgres()
	if err != nil {
		log.Fatalf("Erreur PostgreSQL: %v", err)
	}
	defer db.Close()

	// Connexion Hazelcast
	hz, err := cache.NewHazelcastManager()
	if err != nil {
		log.Fatalf("Erreur Hazelcast: %v", err)
	}
	defer hz.Close()

	// Importation
	importer := services.NewImporter(db, hz)
	if err := importer.ImportAll(); err != nil {
		log.Fatalf("Erreur importation: %v", err)
	}

	log.Println("✅ Importation terminée avec succès")
}
