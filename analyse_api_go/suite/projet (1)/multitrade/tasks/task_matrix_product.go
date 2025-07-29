package tasks

import (
	"fmt"
	"go/mutitrade/cache"
	"log"
	"time"
)

// BuildFeatureMatrix construit une matrice n x 3 à partir des n premiers colleges dans Hazelcast
// Chaque ligne = [IQ, PrevSemResult, CGPA]
func BuildFeatureMatrix(n int) ([][]float64, time.Duration) {
	start := time.Now()
	keys := make([]string, n)
	for i := 1; i <= n; i++ {
		keys[i-1] = fmt.Sprintf("CLG%04d", i)
	}
	colleges, err := cache.GetCollegesFromCacheBulk(keys)
	if err != nil {
		log.Printf("Erreur GetCollegesFromCacheBulk: %v", err)
		return nil, 0
	}
	matrix := make([][]float64, 0, n)
	for _, col := range colleges {
		if col == nil {
			matrix = append(matrix, []float64{0, 0, 0})
			continue
		}
		row := []float64{float64(col.IQ), col.PrevSemResult, col.CGPA}
		matrix = append(matrix, row)
	}
	duration := time.Since(start)
	return matrix, duration
}

// BuildFilterMatrix construit une matrice n x 3 de 0/1 selon un masque booléen sur les colonnes
// mask: [useIQ, usePrevSemResult, useCGPA]
func BuildFilterMatrix(n int, mask [3]bool) [][]int {
	mat := make([][]int, n)
	for i := 0; i < n; i++ {
		row := make([]int, 3)
		for j := 0; j < 3; j++ {
			if mask[j] {
				row[j] = 1
			} else {
				row[j] = 0
			}
		}
		mat[i] = row
	}
	return mat
}

// TaskMatrixProduct étend la logique pour retourner la matrice, la matrice de filtrage, et la durée
func TaskMatrixProduct(n int, mask [3]bool) ([][]float64, [][]int, time.Duration) {
	matrix, duration := BuildFeatureMatrix(n)
	filter := BuildFilterMatrix(n, mask)
	return matrix, filter, duration
}
