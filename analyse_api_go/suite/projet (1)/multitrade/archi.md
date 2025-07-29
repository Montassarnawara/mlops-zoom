bnj a ce moment on va fair un projet de nom multitrade qui faite a conecter a hazelcast et par ce data et ne pas selon la sql mais a laider de reference (banchmark) elle prendre sleon un nb n max 1000 les id des coolege puis dans un code go elle faair la sommme de iq puuis le moy de iq de se college . 
et un autre code qui prand de meme idee un nb n desiner la nb de id  collge on va faiter des matrice matire va remplir va Prev_Sem_Result et un matrice  remplir de CGPA puis een va fait la produit . 
et en fin en va fair  un fin file go qui selon l id de college voir si Placement est yes affiche un msg sion voir si iq > 100 si oui affihce un msg 
fair en debut temsp larchitceture de fiches et que le donner sont toujour dans hazelcast et tout test se par id ne pas par qsl la code et le demarche est toujours 0% sql et a jouter suusi une fonction que selon la file (fonction soit somme soit prosuit de matice ssoit if/eslse ) e le colonnne se sont le nb n par exmp 10 100 200 300 400 ... 1000 
et le table elle remplir par le temps de exuction de chaaque ficher . en fin pour la syntrheser fair une anaylyser total de archi puis en va fair la code un a un pour quiel bien macrhe et la fremwrok est gin debut .



multitrade/
│
├── main.go                    # Point d'entrée (benchmark controller)
├── models/college.go          # Définition de la structure College
├── cache/hazelcast.go         # Connexion + chargement Hazelcast
│
├── tasks/
│   ├── task_sum_iq.go         # Somme + moyenne IQ pour N colleges
├── results/
│   └── benchmark.csv          # Temps d'exécution (par colonne, fonction)



│   ├── task_matrix_product.go  # Produit matriciel Prev_Sem x CGPA
│   ├── task_check.go          # Test Placement == Yes ou IQ > 100
│
++++++++++++++++++++





alors il ya de err musbe dans mac code ou dans l  archi tu ne pas complete myby la map colleges par ce que elle afiiche un err et de meme pou id elle affihce un err dans la focntion et dans la mein afficeh sa en err client := cache.InitHazelcast()  et sa est larchit : multitrade/
│
├── main.go                    # Point d'entrée (benchmark controller)
├── models/college.go          # Définition de la structure College
├── cache/hazelcast.go         # Connexion + chargement Hazelcast
│
├── tasks/
│   ├── task_sum_iq.go         # Somme + moyenne IQ pour N colleges
│   ├── task_matrix_product.go  # Produit matriciel Prev_Sem x CGPA
│   ├── task_check.go          # Test Placement == Yes ou IQ > 100
│
├── results/
│   └── benchmark.csv    et sa le code un par un 1////// package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"go/mutitrade/cache"
	"go/mutitrade/tasks"
)

func main() {
	client := cache.InitHazelcast()

	defer client.Shutdown()

	// Liste des tailles à tester
	sizes := []int{10, 100, 200, 500, 1000}

	// Créer fichier CSV
	file, err := os.Create("results/benchmark.csv")
	if err != nil {
		log.Fatalf("❌ Erreur création fichier CSV : %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Écrire l'entête
	writer.Write([]string{"N", "Tâche", "Durée(ms)", "Résultat"})

	for _, n := range sizes {
		fmt.Printf("\nBenchmark pour N = %d\n", n)

		// ➤ Task 1 : Somme et moyenne IQ
		sum, avg, dur1 := tasks.TaskSumIQ(client, n)
		writer.Write([]string{strconv.Itoa(n), "Somme/Moy IQ", fmt.Sprintf("%.2f", dur1.Seconds()*1000), fmt.Sprintf("Sum=%d, Moy=%.2f", sum, avg)})

		// ➤ Task 2 : Produit matriciel
		prod, dur2 := tasks.TaskMatrixProduct(client, n)
		writer.Write([]string{strconv.Itoa(n), "Produit Matriciel", fmt.Sprintf("%.2f", dur2.Seconds()*1000), fmt.Sprintf("%.2f", prod)})

		// ➤ Task 3 : Vérification conditions
		count, dur3 := tasks.TaskCheckCondition(client, n)
		writer.Write([]string{strconv.Itoa(n), "Placement/Cond IQ", fmt.Sprintf("%.2f", dur3.Seconds()*1000), strconv.Itoa(count)})
	}

	fmt.Println("✅ Benchmarks terminés. Résultats : results/benchmark.csv")
}
2/////// File: projet/multitrade/models/college.go
package models

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
 3///////////package cache

import (
	"context"
	"fmt"
	"go/mutitrade/models"
	"log"
	"math/rand"

	hz "github.com/hazelcast/hazelcast-go-client"
)

var HazelcastClient *hz.Client

// InitHazelcastClient initialise la connexion Hazelcast
func InitHazelcastClient() {
	cfg := hz.Config{}
	client, err := hz.StartNewClientWithConfig(context.Background(), cfg) // Ajout du contexte ici
	if err != nil {
		log.Fatalf("Erreur de connexion Hazelcast : %v", err)
	}
	HazelcastClient = client
	fmt.Println("Hazelcast connecté avec succès.")
}

// LoadCollegesToHazelcast génère et stocke N colleges simulés dans la map Hazelcast
func LoadCollegesToHazelcast(n int) {
	colMap, err := HazelcastClient.GetMap(context.Background(), "colleges") // Contexte ajouté
	if err != nil {
		log.Fatalf("Erreur lors de la récupération de la map Hazelcast : %v", err)
	}

	for i := 1; i <= n; i++ {
		id := fmt.Sprintf("CLG%04d", i)
		col := models.College{
			CollegeID:            id,
			IQ:                   rand.Intn(61) + 80,   // [80, 140]
			PrevSemResult:        rand.Float64()*5 + 4, // [4.0, 9.0]
			CGPA:                 rand.Float64()*5 + 4,
			AcademicPerformance:  rand.Intn(11), // 0–10
			InternshipExperience: randomYesNo(),
			ExtraCurricularScore: rand.Intn(11),
			CommunicationSkills:  rand.Intn(11),
			ProjectsCompleted:    rand.Intn(6),
			Placement:            randomYesNo(),
		}

		if err := colMap.Set(context.Background(), id, col); err != nil { // contexte ajouté
			log.Printf("Erreur lors de l’insertion de %s : %v", id, err)
		}
	}
	fmt.Printf("%d Colleges chargés dans Hazelcast.\n", n)
}

// randomYesNo retourne aléatoirement "Yes" ou "No"
func randomYesNo() string {
	if rand.Intn(2) == 0 {
		return "Yes"
	}
	return "No"
}
4///////////package tasks

import (
	"fmt"
	"go/mutitrade/models"
	"log"
	"time"

	hz "github.com/hazelcast/hazelcast-go-client"
)

func TaskCheckCondition(client *hz.Client, n int) (int, time.Duration) {
	start := time.Now()

	colMap, err := client.GetMap("colleges")
	if err != nil {
		log.Fatalf(" Erreur récupération map Hazelcast : %v", err)
	}

	count := 0

	for i := 1; i <= n; i++ {
		id := fmt.Sprintf("CLG%04d", i)
		val, err := colMap.Get(id)
		if err != nil {
			log.Printf(" Erreur lecture %s : %v", id, err)
			continue
		}

		col, ok := val.(models.College)
		if !ok {
			log.Printf(" Donnée invalide %s", id)
			continue
		}

		if col.Placement == "Yes" || col.IQ > 100 {
			count++
		}
	}

	duration := time.Since(start)
	return count, duration
}
5///////////////package tasks

import (
	"fmt"
	"go/mutitrade/models"
	"log"
	"time"

	hz "github.com/hazelcast/hazelcast-go-client"
)

func TaskMatrixProduct(client *hz.Client, n int) (float64, time.Duration) {
	start := time.Now()

	colMap, err := client.GetMap("colleges")
	if err != nil {
		log.Fatalf(" Erreur récupération map Hazelcast : %v", err)
	}

	var result float64 = 0.0

	for i := 1; i <= n; i++ {
		id := fmt.Sprintf("CLG%04d", i)
		val, err := colMap.Get(id)
		if err != nil {
			log.Printf(" Erreur lecture %s : %v", id, err)
			continue
		}

		col, ok := val.(models.College)
		if !ok {
			log.Printf("Donnée invalide %s", id)
			continue
		}

		result += col.PrevSemResult * col.CGPA
	}

	duration := time.Since(start)
	return result, duration
}
6//////////////////////package tasks

import (
	"fmt"
	"go/mutitrade/models"
	"log"
	"time"

	hz "github.com/hazelcast/hazelcast-go-client"
)

func TaskSumIQ(client *hz.Client, n int) (totalIQ int, averageIQ float64, duration time.Duration) {
	start := time.Now()

	colMap, err := client.GetMap("colleges")
	if err != nil {
		log.Fatalf("Erreur récupération map Hazelcast : %v", err)
	}

	totalIQ = 0

	for i := 1; i <= n; i++ {
		id := fmt.Sprintf("CLG%04d", i)
		val, err := colMap.Get(id)
		if err != nil {
			log.Printf(" Erreur lors de la lecture de %s : %v", id, err)
			continue
		}

		col, ok := val.(models.College)
		if !ok {
			log.Printf(" Donnée invalide pour %s", id)
			continue
		}

		totalIQ += col.IQ
	}

	duration = time.Since(start)
	averageIQ = float64(totalIQ) / float64(n)

	return totalIQ, averageIQ, duration
}
 anlsyer et corriger les fauter dans la ficher mal si non le complet la file absente si exucter et ecrire la code taole de ficher corriger et merci bacq pour votre aider 






 tout sa marche mais unn err est dana main dans la line 17 et si je tape elle ouvert sa ype Config ¶ added in v1.0.0

type Config struct {
	NearCaches            []nearcache.Config                `json:",omitempty"`
	FlakeIDGenerators     map[string]FlakeIDGeneratorConfig `json:",omitempty"`
	Labels                []string                          `json:",omitempty"`
	ClientName            string                            `json:",omitempty"`
	Logger                logger.Config                     `json:",omitempty"`
	Failover              cluster.FailoverConfig            `json:",omitempty"`
	Serialization         serialization.Config              `json:",omitempty"`
	Cluster               cluster.Config                    `json:",omitempty"`
	Stats                 StatsConfig                       `json:",omitempty"`
	NearCacheInvalidation NearCacheInvalidationConfig       `json:",omitempty"`
	// contains filtered or unexported fields
}

Config contains configuration for a client. Zero value of Config is the default configuration.
func NewConfig ¶

func NewConfig() Config

NewConfig creates the default configuration.
func (*Config) AddFlakeIDGenerator ¶ added in v1.1.0

func (c *Config) AddFlakeIDGenerator(name string, prefetchCount int32, prefetchExpiry types.Duration) error

AddFlakeIDGenerator validates the values and adds new FlakeIDGeneratorConfig with the given name.
func (*Config) AddLifecycleListener ¶ added in v1.0.0

func (c *Config) AddLifecycleListener(handler LifecycleStateChangeHandler) types.UUID

AddLifecycleListener adds a lifecycle listener. The listener is attached to the client before the client starts, so all lifecycle events can be received. Use the returned subscription ID to remove the listener. The handler must not block.
func (*Config) AddMembershipListener ¶ added in v1.0.0

func (c *Config) AddMembershipListener(handler cluster.MembershipStateChangeHandler) types.UUID

AddMembershipListener adds a membership listener. The listener is attached to the client before the client starts, so all membership events can be received. Use the returned subscription ID to remove the listener.
func (*Config) AddNearCache ¶ added in v1.3.0

func (c *Config) AddNearCache(cfg nearcache.Config)

AddNearCache adds a near cache configuration.
func (*Config) Clone ¶ added in v1.0.0

func (c *Config) Clone() Config

Clone returns a copy of the configuration.
func (*Config) GetNearCache ¶ added in v1.3.0

func (c *Config) GetNearCache(pattern string) (nearcache.Config, bool, error)

GetNearCache returns the first configuration that matches the given pattern. Returns hzerrors.ErrInvalidConfiguration if the pattern matches more than one configuration.
func (Config) MarshalJSON ¶ added in v1.3.0

func (c Config) MarshalJSON() ([]byte, error)

MarshalJSON marshals the configuration to JSON.
func (*Config) SetLabels ¶ added in v1.0.0

func (c *Config) SetLabels(labels ...string)

SetLabels sets the labels for the client. These labels are displayed in the Hazelcast Management Center.
func (*Config) Validate ¶ added in v1.0.0

func (c *Config) Validate() error

Validate validates the configuration and replaces missing configuration with defaults.
type DistributedObjectEventType ¶ added in v1.0.0

type DistributedObjectEventType string

DistributedObjectEventType describes event type of a distributed object.

const (
	// DistributedObjectCreated is the event type when a distributed object is created.
	DistributedObjectCreated DistributedObjectEventType = "CREATED"
	// DistributedObjectDestroyed is the event type when a distributed object is destroyed.
	DistributedObjectDestroyed DistributedObjectEventType = "DESTROYED"
)

type DistributedObjectNotified ¶ added in v1.0.0

type DistributedObjectNotified struct {
	ServiceName string
	ObjectName  string
	EventType   DistributedObjectEventType
}

 et sa  line lerr efer client.Shutdown() et  sa lerr ::: not enough arguments in call to client.Shutdown
	have ()
	want (context.Context)compilerWrongArgCount











bnj a ce moment on va fair un projet de nom multitrade qui faite a conecter a hazelcast et par ce data et ne pas selon la sql mais a laider de reference (banchmark) elle prendre sleon un nb n les id des colege puis dans un code go elle fair la sommme de iq puuis le moy de iq de se college . 
 
et voila  la architceture des fiches et que le donnee sont toujour dans hazelcast juster fair une fonctio des voir si les donner est bien la et tout test sera  par id ne pas par sql la code et le demarche est toujours 0% sql et ajouter une fonction que selon la file (fonction soit somme ) e le colonnne se sont le nb n par exmp 10 100 200 300 400 ... 1000 
et le table elle remplir par le temps de exuction de chaaque ficher . en fin pour la syntehese fair une anaylse de ce code que je fai est le corrifger et le optmizer par ce que elle representer boucoup des err et tu pas auusi modifer archi si ajuter si supprimer    et voila les code a order comme archi : 



multitrade/
│
├── main.go                    # Point d'entrée (benchmark controller)
├── models/college.go          # Définition de la structure College
├── cache/hazelcast.go         # Connexion + chargement Hazelcast
│
├── tasks/
│   ├── task_sum_iq.go         # Somme + moyenne IQ pour N colleges
├── results/
│   └── benchmark.csv          # Temps d'exécution (par colonne, fonction)

1///package main

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

	sizes := []int{10, 40, 80, 160, 320, 400, 500, 700, 900, 1000, 1200, 1400, 1600, 1800, 2000, 2500, 3000, 4000, 5000, 6000, 7000, 8000, 9000, 10000}

	file, err := os.Create("results/benchmark.csv")
	if err != nil {
		log.Fatalf(" Erreur création fichier CSV : %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Write([]string{"N", "Tâche", "Durée(µs)", "Résultat"})

	for _, n := range sizes {
		fmt.Printf("\n📊 Benchmark pour N = %d\n", n)

		sum, avg, dur1 := tasks.TaskSumIQ(n)
		writer.Write([]string{
			strconv.Itoa(n), "Somme/Moy IQ",
			fmt.Sprintf("%.2f", dur1.Seconds()*1000),
			fmt.Sprintf("Sum=%d, Moy=%.2f", sum, avg),
		})
	}
	fmt.Println(" Benchmarks terminés. Résultats : results/benchmark.csv")

	cache.HazelcastClient.Shutdown(context.Background()) // Fermeture propre du client
}
2/////// File: projet/multitrade/models/college.go
package models

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
3/////package cache

import (
	"context"
	"encoding/json"
	"go/mutitrade/models"
	"log"

	hz "github.com/hazelcast/hazelcast-go-client"
)

var (
	HazelcastClient *hz.Client
	CacheMap        *hz.Map
)

func InitHazelcast() *hz.Client {
	config := hz.Config{}
	config.Cluster.Name = "dev-go"

	client, err := hz.StartNewClientWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Erreur de connexion Hazelcast : %v", err)
	}

	HazelcastClient = client

	// 🟠 Important : on utilise la même map que le projet 1
	CacheMap, err = HazelcastClient.GetMap(context.Background(), "college_cache")
	if err != nil {
		log.Fatalf("Erreur lors de la récupération de la map Hazelcast : %v", err)
	}

	log.Println("✅ Hazelcast connecté avec succès.")
	return client
}

// GetCollegeFromCache récupère un college en JSON et le désérialise
func GetCollegeFromCache(id string) (*models.College, error) {
	raw, err := CacheMap.Get(context.Background(), id)
	if err != nil || raw == nil {
		return nil, err
	}

	jsonBytes, ok := raw.([]byte)
	if !ok {
		return nil, nil
	}

	var college models.College
	if err := json.Unmarshal(jsonBytes, &college); err != nil {
		return nil, err
	}

	return &college, nil
}
4/////package tasks

import (
	"fmt"
	"go/mutitrade/cache"
	"log"
	"time"
)

// TaskSumIQ calcule la somme et la moyenne de l'IQ pour N colleges depuis Hazelcast
func TaskSumIQ(n int) (sum int, avg float64, duration time.Duration) {
	start := time.Now()
	validCount := 0

	for i := 1; i <= n; i++ {
		key := fmt.Sprintf("CLG%04d", i)

		college, err := cache.GetCollegeFromCache(key)
		if err != nil {
			log.Printf("❌ Donnée invalide %s : %v", key, err)
			continue
		}
		if college == nil {
			log.Printf("❌ Donnée vide %s", key)
			continue
		}

		sum += college.IQ
		validCount++
	}

	duration = time.Since(start)
	if validCount > 0 {
		avg = float64(sum) / float64(validCount)
	}
	return
}


elle dit que les tout les doner ne sont pas existe dans hazelcast moi je instaliser les fiche par ce donnee // Initialiser la connexion Hazelcast
func InitHazelcast() {
	config := hazelcast.Config{}
	config.Cluster.Network.SetAddresses("127.0.0.1:5701")
	config.Cluster.Name = "dev-go"

	var err error
	Client, err = hazelcast.StartNewClientWithConfig(Ctx, config)
	if err != nil {
		log.Fatalf(" Connexion Hazelcast échouée : %v", err)
	}

	CacheMap, err = Client.GetMap(Ctx, "college_cache")
	if err != nil {
		log.Fatalf("Impossible d'accéder à la map Hazelcast : %v", err)
	}

	log.Println(" Connexion Hazelcast établie")
}
si tu peut modifer la code pour qui soit de meme connex et map si non je va envoyer la donne pour connection a la base de donner et tu peut fair un ficher qui fair a partir de la base de donner en un 1er temps mais tout data dans hazerlcast et se sa si la 1er idee ne pas marche 

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

et lerr qui je voir  Hazelcast connecté avec succès.
2025/07/18 02:43:52 ❌ Données manquantes dans Hazelcast pour les IDs: [CLG0101 CLG0102 CLG0103 CLG0104 CLG0105 CLG0106 CLG0107 CLG0108 CLG0109 CLG0110 CLG0111 CLG0112 CLG0113 CLG0114 CLG0115 CLG0116 CLG0117 CLG0118 CLG0119 CLG0120 CLG0121 CLG0122 CLG0123 CLG0124 CLG0125 CLG0126 CLG0127 CLG0128 CLG0129 CLG0130 CLG0131 CLG0132 CLG0133 CLG0134 CLG0135 CLG0136 CLG0137 CLG0138 CLG0139 CLG0140 CLG0141 CLG0142 CLG0143 CLG0144 CLG0145 CLG0146 CLG0147 CLG0148 CLG0149 CLG0150 CLG0151 CLG0152 CLG0153 CLG0154 CLG0155 CLG0156 CLG0157 CLG0158 CLG0159 CLG0160 sa pour tout et tu a lacces de modifer dans file de code 