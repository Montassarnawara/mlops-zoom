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
		fmt.Printf("\n Benchmark pour N = %d\n", n)
		for _, field := range fields {
			sum, avg, dur := tasks.TaskSumField(n, field)
			writer.Write([]string{
				strconv.Itoa(n), field,
				fmt.Sprintf("%.0f", dur.Seconds()*1_000_000),
				fmt.Sprintf("%.2f", sum),
				fmt.Sprintf("%.2f", avg),
			})
		}

		// Benchmark de création de matrice features (IQ, PrevSemResult, CGPA)
		mask := [3]bool{true, true, true} // On prend toutes les colonnes
		_, _, matDur := tasks.TaskMatrixProduct(n, mask)
		writer.Write([]string{
			strconv.Itoa(n), "MatrixProduct(IQ,PrevSemResult,CGPA)",
			fmt.Sprintf("%.0f", matDur.Seconds()*1_000_000),
			"-", // Pas de somme
			"-", // Pas de moyenne
		})
	}
	fmt.Println(" Benchmarks terminés. Résultats : results/benchmark.csv")

	cache.HazelcastClient.Shutdown(context.Background()) // Fermeture propre du client
}
