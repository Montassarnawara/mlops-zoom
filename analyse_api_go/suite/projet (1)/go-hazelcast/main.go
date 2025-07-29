package main

import (
	"log"

	"go/hazelcast/api/cache" //    selon go.mod
	"go/hazelcast/api/db"
	"go/hazelcast/api/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialisation du cache Hazelcast et de PostgreSQL
	cache.InitHazelcast()
	db.InitDB()

	router := gin.Default()

	// Route principale
	router.GET("/college/:id", handlers.GetCollege)

	log.Println(" API démarrée sur :8080")
	router.Run(":8080")
}
