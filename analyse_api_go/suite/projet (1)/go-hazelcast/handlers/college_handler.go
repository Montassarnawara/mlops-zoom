package handlers

import (
	"net/http"

	"go/hazelcast/api/cache" // Chemin corrigé selon module
	"go/hazelcast/api/db"    // idem

	"github.com/gin-gonic/gin"
)

// GET /college/:id → cherche d'abord dans le cache, sinon depuis PostgreSQL
func GetCollege(c *gin.Context) {
	id := c.Param("id")

	// Étape 1 : Lire depuis le cache Hazelcast
	college, err := cache.GetFromCache(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lecture cache"})
		return
	}
	if college != nil {
		c.JSON(http.StatusOK, gin.H{
			"source":  "cache",
			"college": college,
		})
		return
	}

	// Étape 2 : Lire depuis la base de données
	college, err = db.GetCollegeByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lecture DB"})
		return
	}
	if college == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Aucun collège trouvé"})
		return
	}

	// Étape 3 : Sauvegarder en cache pour les futurs accès
	_ = cache.SetToCache(college)

	c.JSON(http.StatusOK, gin.H{
		"source":  "db",
		"college": college,
	})
}
