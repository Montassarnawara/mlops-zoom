// ===============================================
// 🏨 HOTEL API SERVICE - HAZELCAST INTEGRATION
// ===============================================
// 
// Service API REST pour la gestion d'hôtels avec Hazelcast
// Version: 2.1.0-fixed
// Auteur: Équipe développement
// Date: 2025-07-27
//
// Ce service charge les données depuis un cluster Hazelcast
// et expose des endpoints REST pour l'accès aux données d'hôtels,
// villes, pays, utilisateurs, rôles et contrats.
//
// Architecture:
// - Connexion Hazelcast avec cluster "dev-go" 
// - Cache en mémoire pour les performances
// - API REST avec Fiber framework
// - Sérialisation GOB pour Hazelcast
//
// ===============================================

package main

import (
	"context"        // Gestion des contextes pour Hazelcast
	"encoding/gob"   // Sérialisation pour Hazelcast
	"fmt"           // Formatage des chaînes
	"log"           // Logging système
	"os"            // Variables d'environnement et signaux
	"os/signal"     // Gestion des signaux système
	"strings"       // Manipulation de chaînes
	"sync"          // Synchronisation (mutex)
	"syscall"       // Appels système
	"time"          // Gestion du temps

	"hotel-sync/models" // Nos modèles de données

	// Framework Fiber pour l'API REST
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	
	// Client Hazelcast
	"github.com/hazelcast/hazelcast-go-client"
)

// ===============================================
// 📋 STRUCTURES DE DONNÉES
// ===============================================

// APIResponse - Structure standardisée pour toutes les réponses API
// Utilise cette structure pour maintenir la cohérence des réponses
type APIResponse struct {
	Success   bool        `json:"success"`   // true/false pour indiquer le succès
	Data      interface{} `json:"data,omitempty"`      // Données de la réponse (optionnel)
	Error     string      `json:"error,omitempty"`     // Message d'erreur (optionnel)
	Count     int         `json:"count,omitempty"`     // Nombre d'éléments retournés (optionnel)
	Message   string      `json:"message,omitempty"`   // Message informatif (optionnel)
	Timestamp time.Time   `json:"timestamp"`           // Horodatage de la réponse
}

// HazelcastService - Service de connexion et gestion Hazelcast
// Encapsule la connexion au cluster Hazelcast avec thread-safety
type HazelcastService struct {
	client  *hazelcast.Client // Client de connexion Hazelcast
	context context.Context   // Contexte pour les opérations Hazelcast
	mutex   sync.RWMutex      // Mutex pour la sécurité thread
}

// DataManager - Gestionnaire principal des données en cache
// Centralise toute la logique de chargement et d'accès aux données
// Les maps utilisent les clés appropriées pour chaque type de données
type DataManager struct {
	// Collections de données en mémoire (cache local)
	hotels    map[string]*models.Hotel    // Hotels indexés par HotelKey (string)
	cities    map[int]*models.City        // Villes indexées par ID
	countries map[int]*models.Country     // Pays indexés par ID  
	users     map[int]*models.User        // Utilisateurs indexés par ID
	roles     map[int]*models.Role        // Rôles indexés par ID
	contracts map[int]*models.Contract    // Contrats indexés par ID
	
	// Gestion de la synchronisation et des services
	mutex     sync.RWMutex        // Mutex pour accès concurrent sécurisé
	service   *HazelcastService   // Service Hazelcast associé
	lastSync  time.Time           // Dernière synchronisation des données
}

// ===============================================
// 🌐 VARIABLES GLOBALES
// ===============================================

var (
	dataManager  *DataManager      // Instance globale du gestionnaire de données
	hazelService *HazelcastService // Instance globale du service Hazelcast
	startTime    = time.Now()      // Heure de démarrage pour calculer l'uptime
)

// ===============================================
// 🔌 SERVICE HAZELCAST
// ===============================================

// NewHazelcastService - Créer une nouvelle connexion Hazelcast
// Configure et établit la connexion au cluster Hazelcast "dev-go"
// 
// Retourne:
//   - *HazelcastService: Service configuré et connecté
//   - error: Erreur de connexion si elle existe
func NewHazelcastService() (*HazelcastService, error) {
	// Configuration du cluster Hazelcast
	config := hazelcast.Config{}
	config.Cluster.Name = "dev-go"                    // Nom du cluster
	config.Cluster.Network.SetAddresses("localhost:5701") // Adresse du cluster
	
	// Tentative de connexion
	client, err := hazelcast.StartNewClientWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("erreur connexion Hazelcast: %w", err)
	}

	return &HazelcastService{
		client:  client,
		context: context.Background(),
	}, nil
}

// GetMap - Récupérer une map Hazelcast par nom
// Méthode thread-safe pour accéder aux collections Hazelcast
//
// Paramètres:
//   - name: Nom de la map Hazelcast à récupérer
//
// Retourne:
//   - *hazelcast.Map: Instance de la map Hazelcast
//   - error: Erreur si la map n'est pas accessible
func (h *HazelcastService) GetMap(name string) (*hazelcast.Map, error) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	if h.client == nil {
		return nil, fmt.Errorf("client Hazelcast non connecté")
	}

	return h.client.GetMap(h.context, name)
}

// Close - Fermer proprement la connexion Hazelcast
// Méthode thread-safe pour fermer la connexion
//
// Retourne:
//   - error: Erreur de fermeture si elle existe
func (h *HazelcastService) Close() error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if h.client != nil {
		return h.client.Shutdown(h.context)
	}
	return nil
}

// ===============================================
// 📊 GESTIONNAIRE DE DONNÉES (DATA MANAGER)
// ===============================================

// NewDataManager - Créer un nouveau gestionnaire de données
// Initialise toutes les collections en mémoire et associe le service Hazelcast
//
// Paramètres:
//   - service: Service Hazelcast configuré et connecté
//
// Retourne:
//   - *DataManager: Instance du gestionnaire de données
func NewDataManager(service *HazelcastService) *DataManager {
	return &DataManager{
		// Initialiser toutes les maps vides
		hotels:    make(map[string]*models.Hotel),
		cities:    make(map[int]*models.City),
		countries: make(map[int]*models.Country),
		users:     make(map[int]*models.User),
		roles:     make(map[int]*models.Role),
		contracts: make(map[int]*models.Contract),
		
		// Associer le service et marquer la synchronisation
		service:   service,
		lastSync:  time.Now(),
	}
}

// LoadAllData - Charger toutes les données depuis Hazelcast
// Fonction principale de synchronisation qui charge toutes les collections
// Utilise un mutex pour garantir la thread-safety pendant le chargement
//
// Retourne:
//   - error: Erreur de chargement si elle existe
func (dm *DataManager) LoadAllData() error {
	dm.mutex.Lock()         // Verrouiller pendant le chargement
	defer dm.mutex.Unlock() // Libérer automatiquement à la fin

	log.Println("🔄 Chargement des données depuis Hazelcast...")

	// Vider toutes les collections existantes pour un rechargement propre
	dm.hotels = make(map[string]*models.Hotel)
	dm.cities = make(map[int]*models.City)
	dm.countries = make(map[int]*models.Country)
	dm.users = make(map[int]*models.User)
	dm.roles = make(map[int]*models.Role)
	dm.contracts = make(map[int]*models.Contract)

	// Charger chaque type de données dans l'ordre approprié
	// L'ordre peut être important pour les dépendances (ex: pays avant villes)
	dm.loadCountries() // 1. Charger les pays d'abord
	dm.loadCities()    // 2. Charger les villes (qui référencent les pays)
	dm.loadRoles()     // 3. Charger les rôles
	dm.loadUsers()     // 4. Charger les utilisateurs (qui référencent les rôles)
	dm.loadHotels()    // 5. Charger les hôtels
	dm.loadContracts() // 6. Charger les contrats (qui référencent les hôtels)

	// Mettre à jour l'horodatage de synchronisation
	dm.lastSync = time.Now()

	// Logger le résumé du chargement
	log.Printf("✅ Données chargées: %d hotels, %d cities, %d countries, %d users, %d roles, %d contracts",
		len(dm.hotels), len(dm.cities), len(dm.countries), len(dm.users), len(dm.roles), len(dm.contracts))

	return nil
}

// ===============================================
// 🗂️ MÉTHODES DE CHARGEMENT DES DONNÉES
// ===============================================

// loadCountries - Charger les pays depuis Hazelcast
// Récupère tous les pays de la collection "countries" et les stocke en mémoire
func (dm *DataManager) loadCountries() {
	// Récupérer la map Hazelcast des pays
	countriesMap, err := dm.service.GetMap("countries")
	if err != nil {
		log.Printf("Erreur map countries: %v", err)
		return
	}

	// Récupérer toutes les entrées de la collection
	entries, err := countriesMap.GetEntrySet(dm.service.context)
	if err != nil {
		log.Printf("Erreur entries countries: %v", err)
		return
	}

	// Parcourir chaque entrée et l'ajouter au cache local
	for _, entry := range entries {
		if country, ok := entry.Value.(models.Country); ok {
			dm.countries[country.ID] = &country // Indexer par ID du pays
		}
	}
}

// loadCities - Charger les villes depuis Hazelcast
// Récupère toutes les villes de la collection "cities" et les stocke en mémoire
func (dm *DataManager) loadCities() {
	// Récupérer la map Hazelcast des villes
	citiesMap, err := dm.service.GetMap("cities")
	if err != nil {
		log.Printf("Erreur map cities: %v", err)
		return
	}

	// Récupérer toutes les entrées de la collection
	entries, err := citiesMap.GetEntrySet(dm.service.context)
	if err != nil {
		log.Printf("Erreur entries cities: %v", err)
		return
	}

	// Parcourir chaque entrée et l'ajouter au cache local
	for _, entry := range entries {
		if city, ok := entry.Value.(models.City); ok {
			dm.cities[city.ID] = &city // Indexer par ID de la ville
		}
	}
}

// loadRoles - Charger les rôles depuis Hazelcast
// Récupère tous les rôles de la collection "roles" et les stocke en mémoire
func (dm *DataManager) loadRoles() {
	// Récupérer la map Hazelcast des rôles
	rolesMap, err := dm.service.GetMap("roles")
	if err != nil {
		log.Printf("Erreur map roles: %v", err)
		return
	}

	// Récupérer toutes les entrées de la collection
	entries, err := rolesMap.GetEntrySet(dm.service.context)
	if err != nil {
		log.Printf("Erreur entries roles: %v", err)
		return
	}

	// Parcourir chaque entrée et l'ajouter au cache local
	for _, entry := range entries {
		if role, ok := entry.Value.(models.Role); ok {
			dm.roles[role.ID] = &role // Indexer par ID du rôle
		}
	}
}

// loadUsers - Charger les utilisateurs depuis Hazelcast
// Récupère tous les utilisateurs de la collection "users" et les stocke en mémoire
func (dm *DataManager) loadUsers() {
	// Récupérer la map Hazelcast des utilisateurs
	usersMap, err := dm.service.GetMap("users")
	if err != nil {
		log.Printf("Erreur map users: %v", err)
		return
	}

	// Récupérer toutes les entrées de la collection
	entries, err := usersMap.GetEntrySet(dm.service.context)
	if err != nil {
		log.Printf("Erreur entries users: %v", err)
		return
	}

	// Parcourir chaque entrée et l'ajouter au cache local
	for _, entry := range entries {
		if user, ok := entry.Value.(models.User); ok {
			dm.users[user.ID] = &user // Indexer par ID de l'utilisateur
		}
	}
}

// loadHotels - Charger les hôtels depuis Hazelcast
// Récupère tous les hôtels de la collection "hotels" et les stocke en mémoire
// IMPORTANT: Les hôtels utilisent HotelKey (string) comme clé, pas l'ID
func (dm *DataManager) loadHotels() {
	// Récupérer la map Hazelcast des hôtels
	hotelsMap, err := dm.service.GetMap("hotels")
	if err != nil {
		log.Printf("Erreur map hotels: %v", err)
		return
	}

	// Récupérer toutes les entrées de la collection
	entries, err := hotelsMap.GetEntrySet(dm.service.context)
	if err != nil {
		log.Printf("Erreur entries hotels: %v", err)
		return
	}

	// Parcourir chaque entrée et l'ajouter au cache local
	for _, entry := range entries {
		if hotel, ok := entry.Value.(models.Hotel); ok {
			dm.hotels[hotel.HotelKey] = &hotel // Indexer par HotelKey (string)
		}
	}
}

// loadContracts - Charger les contrats depuis Hazelcast  
// Récupère tous les contrats de la collection "contracts" et les stocke en mémoire
func (dm *DataManager) loadContracts() {
	// Récupérer la map Hazelcast des contrats
	contractsMap, err := dm.service.GetMap("contracts")
	if err != nil {
		log.Printf("Erreur map contracts: %v", err)
		return
	}

	// Récupérer toutes les entrées de la collection
	entries, err := contractsMap.GetEntrySet(dm.service.context)
	if err != nil {
		log.Printf("Erreur entries contracts: %v", err)
		return
	}

	// Parcourir chaque entrée et l'ajouter au cache local
	for _, entry := range entries {
		if contract, ok := entry.Value.(models.Contract); ok {
			dm.contracts[contract.ID] = &contract // Indexer par ID du contrat
		}
	}
}

// ===============================================
// 🔍 MÉTHODES D'ACCÈS AUX DONNÉES
// ===============================================

// GetAllHotels - Récupérer tous les hôtels en mémoire
// Méthode thread-safe qui retourne une copie de tous les hôtels
//
// Retourne:
//   - []*models.Hotel: Slice de tous les hôtels
func (dm *DataManager) GetAllHotels() []*models.Hotel {
	dm.mutex.RLock()         // Verrouillage en lecture seule
	defer dm.mutex.RUnlock() // Libération automatique

	// Créer un slice avec la capacité appropriée pour les performances
	hotels := make([]*models.Hotel, 0, len(dm.hotels))
	
	// Copier tous les hôtels dans le slice
	for _, hotel := range dm.hotels {
		hotels = append(hotels, hotel)
	}
	return hotels
}

// GetHotelByKey - Récupérer un hôtel par sa clé
// Méthode thread-safe pour accéder à un hôtel spécifique
//
// Paramètres:
//   - hotelKey: Clé unique de l'hôtel (ex: "HOT001")
//
// Retourne:
//   - *models.Hotel: Hôtel trouvé (ou nil)
//   - bool: true si l'hôtel existe, false sinon
func (dm *DataManager) GetHotelByKey(hotelKey string) (*models.Hotel, bool) {
	dm.mutex.RLock()         // Verrouillage en lecture seule
	defer dm.mutex.RUnlock() // Libération automatique
	
	// Rechercher l'hôtel dans la map
	hotel, exists := dm.hotels[hotelKey]
	return hotel, exists
}

// SearchHotels - Rechercher des hôtels par texte
// Effectue une recherche textuelle dans le nom, ville et pays des hôtels
// Recherche insensible à la casse
//
// Paramètres:
//   - query: Texte à rechercher
//
// Retourne:
//   - []*models.Hotel: Slice des hôtels correspondants
func (dm *DataManager) SearchHotels(query string) []*models.Hotel {
	dm.mutex.RLock()         // Verrouillage en lecture seule
	defer dm.mutex.RUnlock() // Libération automatique

	// Convertir la requête en minuscules pour une recherche insensible à la casse
	query = strings.ToLower(query)
	var results []*models.Hotel

	// Parcourir tous les hôtels et chercher dans nom, ville, pays
	for _, hotel := range dm.hotels {
		if strings.Contains(strings.ToLower(hotel.Name), query) ||
			strings.Contains(strings.ToLower(hotel.City), query) ||
			strings.Contains(strings.ToLower(hotel.Country), query) {
			results = append(results, hotel)
		}
	}

	return results
}

// GetStats - Récupérer les statistiques des données chargées
// Fournit un aperçu du nombre d'éléments dans chaque collection
//
// Retourne:
//   - map[string]int: Map avec les compteurs de chaque type de données
func (dm *DataManager) GetStats() map[string]int {
	dm.mutex.RLock()         // Verrouillage en lecture seule
	defer dm.mutex.RUnlock() // Libération automatique

	return map[string]int{
		"hotels":    len(dm.hotels),    // Nombre d'hôtels
		"cities":    len(dm.cities),    // Nombre de villes
		"countries": len(dm.countries), // Nombre de pays
		"users":     len(dm.users),     // Nombre d'utilisateurs
		"roles":     len(dm.roles),     // Nombre de rôles
		"contracts": len(dm.contracts), // Nombre de contrats
	}
}

// ===============================================
// 🌐 HANDLERS HTTP (ENDPOINTS API)
// ===============================================

// healthHandler - Endpoint de vérification de santé de l'API
// GET /health - Retourne l'état de l'API et des statistiques
//
// Réponse JSON avec:
// - Status de l'API
// - Version
// - Uptime (temps de fonctionnement)
// - État de la connexion Hazelcast
// - Statistiques des données chargées
func healthHandler(c *fiber.Ctx) error {
	stats := dataManager.GetStats()

	return c.JSON(APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"status":              "healthy",                         // État de l'API
			"version":             "2.1.0-fixed",                    // Version actuelle
			"uptime":              time.Since(startTime).String(),   // Temps de fonctionnement
			"hazelcast_connected": hazelService.client != nil,       // Connexion Hazelcast
			"last_sync":           dataManager.lastSync,             // Dernière synchronisation
			"stats":               stats,                            // Statistiques des données
		},
		Timestamp: time.Now(),
	})
}

// getAllHotelsHandler - Récupérer tous les hôtels
// GET /api/hotels - Retourne la liste complète des hôtels
//
// Réponse JSON avec:
// - Liste de tous les hôtels
// - Nombre total d'hôtels
// - Message informatif
func getAllHotelsHandler(c *fiber.Ctx) error {
	hotels := dataManager.GetAllHotels()

	return c.JSON(APIResponse{
		Success:   true,
		Data:      hotels,                                               // Liste des hôtels
		Count:     len(hotels),                                          // Nombre d'hôtels
		Message:   fmt.Sprintf("Récupération de %d hôtels", len(hotels)), // Message informatif
		Timestamp: time.Now(),
	})
}

// getHotelByKeyHandler - Récupérer un hôtel par sa clé
// GET /api/hotels/:key - Retourne les détails d'un hôtel spécifique
//
// Paramètres URL:
//   - key: Clé unique de l'hôtel (ex: "HOT001")
//
// Réponse JSON avec:
// - Détails complets de l'hôtel
// - Message avec le nom de l'hôtel
// OU erreur 404 si l'hôtel n'existe pas
func getHotelByKeyHandler(c *fiber.Ctx) error {
	hotelKey := c.Params("key") // Récupérer la clé depuis l'URL

	// Chercher l'hôtel dans le cache
	hotel, exists := dataManager.GetHotelByKey(hotelKey)
	if !exists {
		// Retourner une erreur 404 si l'hôtel n'existe pas
		return c.Status(404).JSON(APIResponse{
			Success:   false,
			Error:     "Hôtel non trouvé",
			Timestamp: time.Now(),
		})
	}

	// Retourner les détails de l'hôtel
	return c.JSON(APIResponse{
		Success:   true,
		Data:      hotel,                                             // Données de l'hôtel
		Message:   fmt.Sprintf("Détails de l'hôtel %s", hotel.Name), // Message avec nom
		Timestamp: time.Now(),
	})
}

// searchHotelsHandler - Rechercher des hôtels par texte
// GET /api/hotels/search?q=terme - Recherche textuelle dans les hôtels
//
// Paramètres Query:
//   - q: Terme de recherche (obligatoire)
//
// Réponse JSON avec:
// - Liste des hôtels correspondants
// - Nombre de résultats
// - Message avec le terme recherché
// OU erreur 400 si le paramètre 'q' est manquant
func searchHotelsHandler(c *fiber.Ctx) error {
	query := c.Query("q") // Récupérer le terme de recherche
	if query == "" {
		// Retourner une erreur si le paramètre est manquant
		return c.Status(400).JSON(APIResponse{
			Success:   false,
			Error:     "Paramètre 'q' requis",
			Timestamp: time.Now(),
		})
	}

	// Effectuer la recherche
	hotels := dataManager.SearchHotels(query)

	// Retourner les résultats
	return c.JSON(APIResponse{
		Success:   true,
		Data:      hotels,                                                        // Résultats de la recherche
		Count:     len(hotels),                                                   // Nombre de résultats
		Message:   fmt.Sprintf("Trouvé %d hôtel(s) pour '%s'", len(hotels), query), // Message informatif
		Timestamp: time.Now(),
	})
}

// reloadDataHandler - Recharger les données depuis Hazelcast
// GET /api/reload OU POST /api/reload - Force le rechargement des données
//
// Réponse JSON avec:
// - Statistiques des données rechargées
// - Message de confirmation
// OU erreur 500 si le rechargement échoue
func reloadDataHandler(c *fiber.Ctx) error {
	log.Println("🔄 Rechargement des données demandé...")

	// Tenter le rechargement des données
	err := dataManager.LoadAllData()
	if err != nil {
		// Retourner une erreur serveur si le rechargement échoue
		return c.Status(500).JSON(APIResponse{
			Success:   false,
			Error:     err.Error(),
			Timestamp: time.Now(),
		})
	}

	// Retourner les nouvelles statistiques
	return c.JSON(APIResponse{
		Success:   true,
		Data:      dataManager.GetStats(),           // Statistiques actualisées
		Message:   "Données rechargées avec succès", // Message de confirmation
		Timestamp: time.Now(),
	})
}

// ===============================================
// ⚙️ CONFIGURATION ET INITIALISATION
// ===============================================

// initGobTypes - Initialiser les types GOB pour Hazelcast
// Enregistre tous les types de données pour la sérialisation GOB
// OBLIGATOIRE: Doit être appelé avant toute interaction avec Hazelcast
func initGobTypes() {
	log.Println("🔧 Enregistrement des types GOB...")

	// Enregistrer tous nos modèles de données pour la sérialisation
	gob.Register(models.Hotel{})    // Hôtels
	gob.Register(models.City{})     // Villes  
	gob.Register(models.Country{})  // Pays
	gob.Register(models.User{})     // Utilisateurs
	gob.Register(models.Role{})     // Rôles
	gob.Register(models.Contract{}) // Contrats
	gob.Register(time.Time{})       // Type Time pour les dates

	log.Println("✅ Types GOB enregistrés")
}

// ===============================================
// 🚀 FONCTION PRINCIPALE
// ===============================================

// main - Point d'entrée de l'application
// Orchestree le démarrage complet du service API
func main() {
	log.Println("🚀 Démarrage Hotel API Service - Version Corrigée")

	// ===== ÉTAPE 1: INITIALISATION GOB =====
	// Enregistrer les types pour la sérialisation Hazelcast
	initGobTypes()

	// ===== ÉTAPE 2: CONNEXION HAZELCAST =====
	log.Println("🔌 Connexion à Hazelcast...")
	var err error
	hazelService, err = NewHazelcastService()
	if err != nil {
		log.Fatalf("❌ Erreur connexion Hazelcast: %v", err)
	}
	log.Println("✅ Connecté à Hazelcast")

	// ===== ÉTAPE 3: INITIALISATION DATA MANAGER =====
	// Créer le gestionnaire de données avec le service Hazelcast
	dataManager = NewDataManager(hazelService)

	// ===== ÉTAPE 4: CHARGEMENT INITIAL DES DONNÉES =====
	// Charger toutes les données depuis Hazelcast au démarrage
	if err := dataManager.LoadAllData(); err != nil {
		log.Printf("⚠️ Erreur chargement initial: %v", err)
	}

	// ===== ÉTAPE 5: CONFIGURATION DE L'API FIBER =====
	// Créer l'application Fiber avec configuration
	app := fiber.New(fiber.Config{
		AppName: "Hotel API Service", // Nom de l'application
	})

	// ===== ÉTAPE 6: MIDDLEWARES =====
	app.Use(cors.New())    // CORS pour les requêtes cross-origin
	app.Use(logger.New())  // Logging des requêtes HTTP
	app.Use(recover.New()) // Récupération automatique des panics

	// ===== ÉTAPE 7: DÉFINITION DES ROUTES =====
	
	// Route racine - Page d'accueil de l'API
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "🏨 Hotel API Service avec vraies données Hazelcast",
			"version": "2.1.0-fixed",
			"time":    time.Now(),
		})
	})

	// Routes de l'API
	app.Get("/health", healthHandler)                      // État de santé de l'API
	app.Get("/api/hotels", getAllHotelsHandler)            // Lister tous les hôtels
	app.Get("/api/hotels/:key", getHotelByKeyHandler)      // Détails d'un hôtel par clé
	app.Get("/api/hotels/search", searchHotelsHandler)     // Recherche d'hôtels
	app.Get("/api/reload", reloadDataHandler)              // GET pour faciliter les tests
	app.Post("/api/reload", reloadDataHandler)             // POST pour les appels programmatiques

	// ===== ÉTAPE 8: ARRÊT GRACIEUX =====
	// Gérer l'arrêt propre de l'application
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM) // Écouter CTRL+C et SIGTERM
		<-c
		log.Println("🛑 Arrêt gracieux...")
		
		// Fermer proprement la connexion Hazelcast
		if hazelService != nil {
			hazelService.Close()
		}
		
		// Arrêter le serveur Fiber
		app.Shutdown()
		os.Exit(0)
	}()

	// ===== ÉTAPE 9: DÉMARRAGE DU SERVEUR =====
	// Déterminer le port d'écoute
	port := "3000"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort // Utiliser PORT si défini dans l'environnement
	}

	// Afficher les informations de démarrage
	log.Printf("🌐 Serveur démarré sur le port %s", port)
	log.Println("📋 Routes disponibles:")
	log.Println("  GET  /health")
	log.Println("  GET  /api/hotels")
	log.Println("  GET  /api/hotels/:key")
	log.Println("  GET  /api/hotels/search?q=...")
	log.Println("  GET  /api/reload")
	log.Println("  POST /api/reload")

	// Démarrer le serveur HTTP
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("❌ Erreur démarrage serveur: %v", err)
	}
}
