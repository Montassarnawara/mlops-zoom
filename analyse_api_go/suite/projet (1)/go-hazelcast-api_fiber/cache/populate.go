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
