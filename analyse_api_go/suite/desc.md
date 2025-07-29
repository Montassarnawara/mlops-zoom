++++++++++++++++++++++++++++
go-hazelcast-api/
│
├── main.go
├── go.mod
├── models/
│   └── user.go
├── cache/
│   └── hazelcast.go
├── db/
│   └── postgres.go
├── handlers/
│   └── user_handler.go

+++++++++++++++++++++
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
main.go

But : Point d’entrée principal de l’application

Contient :

    La fonction main()

    L’appel aux fonctions d'initialisation :

        Connexion à Hazelcast

        Connexion à PostgreSQL

        Démarrage des routes API avec Gin

💡 C’est le fichier coordonnateur du projet. Il ne contient presque aucune logique, juste les appels aux composants définis ailleurs.
🔹 go.mod

But : Fichier de configuration Go Modules

Contient :

    Le nom du module (module go-hazelcast-api)

    Les dépendances : Gin, Hazelcast, GORM, PostgreSQL driver…

💡 Ce fichier est automatiquement généré avec go mod init puis mis à jour avec go get ....
🔹 models/user.go

But : Définir la structure de données principale, ici un User.

Contient :

type User struct {
    ID         string `json:"id" gorm:"primaryKey"`
    Name       string `json:"name"`
    Validation bool   `json:"validation"`
}

💡 Cette structure sera utilisée dans :

    Les requêtes JSON de l’API

    La table users dans PostgreSQL

    Les objets en cache dans Hazelcast

🔹 cache/hazelcast.go

But : Initialiser Hazelcast et exposer ses fonctions.

Contient :

    Une fonction InitHazelcast() qui connecte Hazelcast

    Des fonctions comme :

        GetUserFromCache(id string)

        SetUserToCache(user User)

        DeleteUserFromCache(id string) (optionnel)

💡 Tout ce qui touche à Hazelcast est centralisé ici.
🔹 db/postgres.go

But : Gérer la connexion PostgreSQL + CRUD

Contient :

    La fonction InitPostgres() pour la connexion

    Des fonctions :

        GetUserFromDB(id string)

        CreateUserInDB(user User)

        GetAllUsersFromDB()

💡 Cette couche est responsable de toute interaction avec la vraie base de données.
🔹 handlers/user_handler.go

But : Contenir les fonctions utilisées dans les routes HTTP

Exemples :

    func GetUser(c *gin.Context)

        Regarde dans Hazelcast

        Sinon, va chercher dans PostgreSQL et met en cache

    func CreateUser(c *gin.Context)

        Ajoute dans PostgreSQL

        Puis dans Hazelcast

    func GetAllUsers(c *gin.Context)

        Retourne tous les users en DB ou cache (selon logique choisie)

💡 Cette couche est l’interface entre l’utilisateur final et la logique serveur.
🔁 Résumé : rôles par couche
Couche / Dossier	Rôle principal
main.go	Coordination générale, démarrage du serveur
models/	Définition des structures de données (User)
cache/	Connexion Hazelcast, lecture/écriture cache
db/	Connexion PostgreSQL, lecture/écriture DB
handlers/	Logique métier de l'API (GET, POST…)
go.mod	Configuration Go Modules + dépendances
+++++++++++++++++++++++++++++++++++++++++++++++++++++++++






package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"go/mutitrade/cache"
	"go/mutitrade/tasks"
)

func main() {
	// Initialisation globale de Hazelcast
	cache.InitHazelcast()

	// Compte le nombre de clés CLGxxxx présentes dans Hazelcast
	nbColleges := 0
	for i := 1; ; i++ {
		key := fmt.Sprintf("CLG%04d", i)
		exists, err := cache.CacheMap.ContainsKey(context.Background(), key)
		if err != nil || !exists {
			break
		}
		nbColleges++
	}
	if nbColleges == 0 {
		log.Fatalf("❌ Aucune donnée CLGxxxx trouvée dans Hazelcast !")
	}
	log.Printf("✅ %d colleges trouvés dans Hazelcast.", nbColleges)

	// Limite sizes à ce qui est possible
	sizes := []int{1, 2, 4, 6, 8, 10, 20, 25, 30, 35, 40, 45, 50, 55, 60, 65, 70, 75, 80, 85, 90, 95, 100}
	var filteredSizes []int
	for _, n := range sizes {
		if n <= nbColleges {
			filteredSizes = append(filteredSizes, n)
		}
	}
	if len(filteredSizes) == 0 {
		log.Fatalf("❌ Pas assez de données pour le benchmark. Ajoute plus de colleges dans Hazelcast.")
	}

	// ...

	file, err := os.Create("results/benchmark.csv")
	if err != nil {
		log.Fatalf(" Erreur création fichier CSV : %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Write([]string{"N", "Colonne", "Durée(µs)", "Somme", "Moyenne"})

	fields := []string{"IQ", "CGPA", "AcademicPerformance"}
	for _, n := range filteredSizes {
		fmt.Printf("\n📊 Benchmark pour N = %d\n", n)
		for _, field := range fields {
			sum, avg, dur := tasks.TaskSumField(n, field)
			writer.Write([]string{
				strconv.Itoa(n), field,
				fmt.Sprintf("%.0f", dur.Seconds()*1_000_000),
				fmt.Sprintf("%.2f", sum),
				fmt.Sprintf("%.2f", avg),
			})
		}
	}
	fmt.Println(" Benchmarks terminés. Résultats : results/benchmark.csv")

	cache.HazelcastClient.Shutdown(context.Background()) // Fermeture propre du client
}
