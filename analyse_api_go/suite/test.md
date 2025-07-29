package cache

import (
	"log"
	"go-hazelcast-api/db"
)

// Injecte les 100 premières lignes dans Hazelcast
func PopulateHazelcastWithColleges() {
	colleges, err := db.GetFirst100Colleges()
	if err != nil {
		log.Fatalf("❌ Erreur lors de la récupération des données : %v", err)
	}

	for _, college := range colleges {
		err := SetToCache(&college)
		if err != nil {
			log.Printf("⚠️ Erreur insertion dans Hazelcast pour %s : %v", college.CollegeID, err)
		}
	}

	log.Printf("✅ %d colleges insérés dans Hazelcast", len(colleges))
}







package main

import (
        "fmt"
        "go-hazelcast-api/cache"
        "go-hazelcast-api/routes"
)

func main() {
        // Connexion Hazelcast
        cache.InitHazelcast()
        defer cache.Client.Shutdown(cache.Ctx)

        fmt.Println(" ^=^z^` API d  marr  e avec Hazelcast uniquement")

        // Lancer le serveur web
        router := routes.SetupRoutes()
        router.Run(":8080")
}




package cache

import (
	"go-hazelcast-api/db"
	"log"
)

// Injecte les 100 premières lignes dans Hazelcast
func PopulateHazelcastWithColleges() {
	colleges, err := db.GetFirst100Colleges()
	if err != nil {
		log.Fatalf("❌ Erreur lors de la récupération des données : %v", err)
	}

	for _, college := range colleges {
		err := SetToCache(&college)
		if err != nil {
			log.Printf("⚠️ Erreur insertion dans Hazelcast pour %s : %v", college.CollegeID, err)
		}
	}

	log.Printf("✅ %d colleges insérés dans Hazelcast", len(colleges))
}








package main

import (
        "log"

        "go-hazelcast-api/cache"
        "go-hazelcast-api/db"
)

func main() {
        if err := db.InitDB(); err != nil {
                log.Fatalf(" ^}^l Erreur connexion PostgreSQL : %v", err)
        }
        defer db.DB.Close()

        cache.InitHazelcast()
        defer cache.Client.Shutdown(cache.Ctx)

        cache.PopulateHazelcastWithColleges()

        log.Println(" ^=^z^` Initialisation termin  e.")
}





package cache

import (
        "go-hazelcast-api/db"
        "log"
)

// Injecte les 100 premi  res lignes dans Hazelcast
func PopulateHazelcastWithColleges() {
        colleges, err := db.GetFirst100Colleges()
        if err != nil {
                log.Fatalf(" ^}^l Erreur lors de la r  cup  ration des donn  es : %v", err)
        }

        for _, college := range colleges {
                err := SetToCache(&college)
                if err != nil {
                        log.Printf(" ^z   ^o Erreur insertion dans Hazelcast pour %s : %v", college.CollegeID, err)
                }
        }

        log.Printf(" ^|^e %d colleges ins  r  s dans Hazelcast", len(colleges))
}








///////////////
package controllers

import (
	"net/http"

	"go-hazelcast-api/cache"

	"github.com/gin-gonic/gin"
)

// GET /college/:id
func GetCollegeByID(c *gin.Context) {
	id := c.Param("id")

	college, err := cache.GetFromCache(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de l'accès au cache"})
		return
	}

	if college == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Collège non trouvé dans le cache"})
		return
	}

	c.JSON(http.StatusOK, college)
}
s


package main

import (
	"fmt"
	"go-hazelcast-api/cache"
	"go-hazelcast-api/routes"
)

func main() {
	// Connexion Hazelcast
	cache.InitHazelcast()
	defer cache.Client.Shutdown(cache.Ctx)

	fmt.Println(" ^=^z^` API d  marr  e avec Hazelcast uniquement")

	// Lancer le serveur web
	router := routes.SetupRoutes()
	router.Run(":8080")
}




/////////////////////////////////////////////
