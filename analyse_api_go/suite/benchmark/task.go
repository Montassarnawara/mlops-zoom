package main

import (
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// go get -u gonum.org/v1/gonum/...
// go get -u github.com/gin-gonic/gin
// Génère une matrice carrée n x n d'entiers aléatoires entre 0 et 9
func generateRandomMatrix(n int) [][]int {
	matrix := make([][]int, n)
	for i := 0; i < n; i++ {
		matrix[i] = make([]int, n)
		for j := 0; j < n; j++ {
			matrix[i][j] = rand.Intn(10) // valeurs entre 0 et 9
		}
	}
	return matrix
}

// Multiplie deux matrices carrées n x n

func multiplyMatrices(a, b [][]int) [][]int {
	n := len(a)
	result := make([][]int, n)
	for i := 0; i < n; i++ {
		result[i] = make([]int, n)
		for j := 0; j < n; j++ {
			for k := 0; k < n; k++ {
				result[i][j] += a[i][k] * b[k][j]
			}
		}
	}
	return result
}

// Handler Gin pour générer une matrice aléatoire de dimension n
func getMatrixHandler() gin.HandlerFunc {

	return func(c *gin.Context) {
		start := time.Now() // Début mesure
		n, err := strconv.Atoi(c.Param("n"))
		if err != nil || n <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Dimension invalide"})
			return
		}

		matrix := generateRandomMatrix(n)

		duration := time.Since(start).Microseconds() //  Fin mesure

		c.JSON(http.StatusOK, gin.H{
			"dim":                   n,
			"matrice":               matrix,
			"duration_microseconds": duration, //
		})
	}
}

// Handler Gin pour générer deux matrices aléatoires et leur produit
// Handler Gin pour générer deux matrices aléatoires et leur produit avec mesure de durée
func getMatrixProductHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now() // début chrono
		n, err := strconv.Atoi(c.Param("n"))
		if err != nil || n <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Dimension invalide"})
			return
		}

		matrix1 := generateRandomMatrix(n)
		matrix2 := generateRandomMatrix(n)
		product := multiplyMatrices(matrix1, matrix2)

		duration := time.Since(start).Microseconds() // durée en µs

		c.JSON(http.StatusOK, gin.H{
			"dim":                   n,
			"matrice1":              matrix1,
			"matrice2":              matrix2,
			"matrice_prod":          product,
			"duration_microseconds": duration,
		})
	}
}

// Multiplie deux matrices gonum de type mat.Dense
// Utilise gonum pour la multiplication de matrices
// Génère une matrice carrée n x n d'entiers aléatoires entre 0 et 9
/*func multiplyGonumMatrices(a, b *mat.Dense) *mat.Dense {
	var result mat.Dense
	result.Mul(a, b)
	return &result
}

func generateRandomGonumMatrix(n int) *mat.Dense {
	data := make([]float64, n*n)
	for i := range data {
		data[i] = float64(rand.Intn(10))
	}
	return mat.NewDense(n, n, data)
}

func convertMatrixTo2DInt(m *mat.Dense) [][]int {
	r, c := m.Dims()
	result := make([][]int, r)
	for i := 0; i < r; i++ {
		row := make([]int, c)
		for j := 0; j < c; j++ {
			row[j] = int(m.At(i, j))
		}
		result[i] = row
	}
	return result
}

func getMatrixHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		n, err := strconv.Atoi(c.Param("n"))
		if err != nil || n <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Dimension invalide"})
			return
		}

		matrix := generateRandomGonumMatrix(n)
		duration := time.Since(start).Microseconds()

		c.JSON(http.StatusOK, gin.H{
			"dim":                   n,
			"matrice":               convertMatrixTo2DInt(matrix),
			"duration_microseconds": duration,
		})
	}
}

func getMatrixProductHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		n, err := strconv.Atoi(c.Param("n"))
		if err != nil || n <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Dimension invalide"})
			return
		}

		a := generateRandomGonumMatrix(n)
		b := generateRandomGonumMatrix(n)
		result := multiplyGonumMatrices(a, b)

		duration := time.Since(start).Microseconds()

		c.JSON(http.StatusOK, gin.H{
			"dim":                   n,
			"matrix_a":              convertMatrixTo2DInt(a),
			"matrix_b":              convertMatrixTo2DInt(b),
			"result":                convertMatrixTo2DInt(result),
			"duration_microseconds": duration,
		})
	}
}
*/
// Génère un vecteur de n réels aléatoires entre 0.0 et 100.0
func generateRandomVector(n int) []float64 {
	vec := make([]float64, n)
	for i := 0; i < n; i++ {
		vec[i] = rand.Float64() * 100
	}
	return vec
}

// Calcule la somme, la moyenne et le produit (protégé) d'un vecteur
func analyzeVector(vec []float64) (sum float64, mean float64, prod float64, valid bool) {
	sum = 0
	prod = 1
	valid = true

	for _, val := range vec {
		sum += val
		prod *= val

		// Vérifie si le produit devient trop grand ou anormal
		if math.IsInf(prod, 0) || math.IsNaN(prod) {
			valid = false
			break
		}
	}

	if len(vec) > 0 {
		mean = sum / float64(len(vec))
	}

	return
}

// Handler pour générer un vecteur avec analyse (somme, moyenne, produit)
func getVectorAnalysisHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		n, err := strconv.Atoi(c.Param("n"))
		if err != nil || n <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Paramètre invalide"})
			return
		}

		vec := generateRandomVector(n)
		sum, mean, prod, valid := analyzeVector(vec)

		duration := time.Since(start).Microseconds()

		if !valid {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":                 "Produit trop grand ou non calculable (overflow)",
				"dim":                   n,
				"somme":                 sum,
				"moyenne":               mean,
				"produit":               "invalid",
				"duration_microseconds": duration,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"dim":                   n,
			"somme":                 sum,
			"moyenne":               mean,
			"produit":               prod,
			"duration_microseconds": duration,
		})
	}
}

// (Facultatif) Handler pour juste générer un vecteur sans analyse
func getVectorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		n, err := strconv.Atoi(c.Param("n"))
		if err != nil || n <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Paramètre invalide"})
			return
		}

		vec := generateRandomVector(n)

		duration := time.Since(start).Microseconds()

		c.JSON(http.StatusOK, gin.H{
			"dim":                   n,
			"vecteur":               vec,
			"duration_microseconds": duration,
		})
	}
}

// Génère un vecteur binaire (0 ou 1) de taille n
func generateBinaryVector(n int) []int {
	vec := make([]int, n)
	for i := 0; i < n; i++ {
		vec[i] = rand.Intn(2) // génère 0 ou 1
	}
	return vec
}

// Teste si une offre existe à l'index donné dans le vecteur binaire
func testOfferFromVector(vec []int, index int) string {
	if index < 0 || index >= len(vec) {
		return "Index hors limite"
	}
	if vec[index-1] != 0 {
		return "Offre existante"
	}
	return "Offre non existante"
}

// Handler Gin pour tester la présence d'une offre à un index donné dans un vecteur binaire
func promoHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now() // Début du chronomètre
		index, err := strconv.Atoi(c.Param("index"))
		if err != nil || index < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Index invalide"})
			return
		}

		const vectorSize = 200 // Taille du vecteur binaire

		vec := generateBinaryVector(vectorSize)
		result := testOfferFromVector(vec, index)

		duration := time.Since(start).Microseconds() // Durée en µs

		c.JSON(http.StatusOK, gin.H{
			"vecteur":               vec,
			"index":                 index,
			"résultat":              result,
			"duration_microseconds": duration,
		})
	}
}

func main() {
	rand.Seed(time.Now().UnixNano()) // Initialisation aléatoire

	r := gin.Default()
	r.GET("/matrice/:n", getMatrixHandler())                 // Génère une matrice aléatoire
	r.GET("/pro_matrice/:n", getMatrixProductHandler())      // Génère deux matrices et leur produit
	r.GET("/vecteur/:n", getVectorHandler())                 // Génère un vecteur aléatoire
	r.GET("/analyse_vecteur/:n", getVectorAnalysisHandler()) // Analyse un vecteur aléatoire
	r.GET("/condition/:index", promoHandler())               // Teste la présence d'une offre

	r.Run(":8080")
}

/*
// curl http://localhost:8080/vecteur/5
# {
#   "dim": 5,
#   "vecteur": [34.5, 12.7, 98.1, 43.2, 67.3]
# }

// curl http://localhost:8080/analyse_vecteur/5
# {
#   "dim": 5,
#   "vecteur": [12.3, 89.1, 54.6, 78.0, 44.2],
#   "somme": 278.2,
#   "moyenne": 55.64,
#   "produit": 18439824.6
# }

# curl http://localhost:8080/matrice/3
# {
#   "dim": 3,
#   "matrice": [[3,1,7],[4,0,2],[9,5,6]]
# }

# curl http://localhost:8080/pro_matrice/3
# {
#   "dim": 3,
#   "matrice1": [...],
#   "matrice2": [...],
#   "matrice_prod": [...]
# }


curl http://localhost:8080/vecteur/10

curl http://localhost:8080/condition/1
*/
