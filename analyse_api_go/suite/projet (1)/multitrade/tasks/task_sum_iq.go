package tasks

import (
	"fmt"
	"go/mutitrade/cache"
	"time"
)

// TaskSumIQ calcule la somme et la moyenne de l'IQ pour N colleges depuis Hazelcast
func TaskSumIQ(n int) (sum int, avg float64, duration time.Duration) {
	// Mesure le temps sur toute l'opération (récupération + calcul)
	// Prépare la liste des clés à lire
	keys := make([]string, n)
	for i := 1; i <= n; i++ {
		keys[i-1] = fmt.Sprintf("CLG%04d", i)
	}
	start := time.Now()
	colleges, err := cache.GetCollegesFromCacheBulk(keys)
	duration = time.Since(start)
	if err != nil {
		// Handle the error as appropriate, here we just return zero values
		return 0, 0, duration
	}
	sum = 0
	validCount := 0
	for _, college := range colleges {
		if college == nil {
			continue
		}
		sum += college.IQ
		validCount++
	}
	if validCount > 0 {
		avg = float64(sum) / float64(validCount)
	}
	return
}

// TaskSumField calcule la somme et la moyenne d'un champ numérique pour N colleges depuis Hazelcast
func TaskSumField(n int, field string) (sum float64, avg float64, duration time.Duration) {
	   keys := make([]string, n)
	   for i := 1; i <= n; i++ {
			   keys[i-1] = fmt.Sprintf("CLG%04d", i)
	   }
	   start := time.Now()
	   colleges, err := cache.GetCollegesFromCacheBulk(keys)
	   duration = time.Since(start)
	   if err != nil {
			   return 0, 0, duration
	   }
	   sum = 0
	   validCount := 0
	   for _, college := range colleges {
			   if college == nil {
					   continue
			   }
			   switch field {
			   case "IQ":
					   sum += float64(college.IQ)
			   case "CGPA":
					   sum += college.CGPA
			   case "AcademicPerformance":
					   sum += float64(college.AcademicPerformance)
			   case "PrevSemResult":
					   sum += college.PrevSemResult
			   case "ExtraCurricularScore":
					   sum += float64(college.ExtraCurricularScore)
			   case "CommunicationSkills":
					   sum += float64(college.CommunicationSkills)
			   case "ProjectsCompleted":
					   sum += float64(college.ProjectsCompleted)
			   }
			   validCount++
	   }
	   if validCount > 0 {
			   avg = sum / float64(validCount)
	   }
	   return
}
