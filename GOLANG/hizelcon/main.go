package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hazelcast/hazelcast-go-client"
)

// --- STRUCTURE DES DONNÉES ---

type Data struct {
	Id         string `json:"Id"`
	Name       string `json:"Name"`
	Validation bool   `json:"Validation"`
}

// --- VARIABLES GLOBALES ---

var (
	ctx      = context.Background()
	hzClient *hazelcast.Client
	usersMap *hazelcast.Map
)

// --- INITIALISATION DE HAZELCAST ---

func initHazelcast() {
	config := hazelcast.Config{}
	config.Cluster.Network.SetAddresses("127.0.0.1:5701")
	config.Cluster.Name = "dev-go"

	var err error
	hzClient, err = hazelcast.StartNewClientWithConfig(ctx, config)
	if err != nil {
		log.Fatalf(" Erreur de connexion Hazelcast : %v", err)
	}

	usersMap, err = hzClient.GetMap(ctx, "users")
	if err != nil {
		log.Fatalf(" Erreur GetMap : %v", err)
	}

	log.Println(" Connexion Hazelcast établie")
}

// --- HANDLER : Récupérer un utilisateur par ID ---

func getUserByID(c *gin.Context) {
	id := c.Param("id")
	raw, err := usersMap.Get(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if raw == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Utilisateur non trouvé"})
		return
	}

	jsonBytes, ok := raw.([]byte)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur de type (pas []byte)"})
		return
	}

	var user Data
	if err := json.Unmarshal(jsonBytes, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur JSON: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// --- HANDLER : Récupérer tous les utilisateurs (GET /user) ---

func getAllUsers(c *gin.Context) {
	keys, _ := usersMap.GetKeySet(ctx)
	var all []Data

	for _, key := range keys {
		raw, _ := usersMap.Get(ctx, key)
		jsonBytes, ok := raw.([]byte)
		if ok {
			var user Data
			json.Unmarshal(jsonBytes, &user)
			all = append(all, user)
		}
	}
	c.JSON(http.StatusOK, all)
}

// --- HANDLER : Ajouter un utilisateur (POST /user) ---

func createUser(c *gin.Context) {
	var newUser Data
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jsonBytes, _ := json.Marshal(newUser)
	if err := usersMap.Set(ctx, newUser.Id, jsonBytes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Utilisateur ajouté"})
}

// --- HANDLER : Vérification de santé (GET /health) ---

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// --- HANDLER : Accès à tous les "fichiers" (GET /user/all_files) ---

func getAllFiles(c *gin.Context) {
	//  Supposé ici que chaque utilisateur représente un "fichier"
	// Si les fichiers sont une autre structure, il faudra adapter.
	keys, _ := usersMap.GetKeySet(ctx)
	c.JSON(http.StatusOK, gin.H{
		"files_ids": keys,
		"count":     len(keys),
	})
}

// --- ROUTES PRINCIPALES ---

func setupRoutes() *gin.Engine {
	router := gin.Default()

	router.GET("/user", getAllUsers)
	router.GET("/user/:id", getUserByID)
	router.POST("/user", createUser)
	router.GET("/health", healthCheck)
	router.GET("/user/all_files", getAllFiles)

	return router
}

// --- MAIN PROGRAMME ---

func main() {
	initHazelcast()
	defer hzClient.Shutdown(ctx)

	router := setupRoutes()
	router.Run(":8081")
}
