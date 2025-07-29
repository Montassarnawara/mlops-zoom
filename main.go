package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// Structure qui représente les colonnes de la table "college"
type College struct {
	CommunicationSkills  int64   `json:"communication_skills"`
	ProjectsCompleted    int64   `json:"projects_completed"`
	PrevSemResult        float64 `json:"prev_sem_result"`
	CGPA                 float64 `json:"cgpa"`
	AcademicPerformance  int64   `json:"academic_performance"`
	IQ                   int64   `json:"iq"`
	ExtraCurricularScore int64   `json:"extra_curricular_score"`
	Placement            string  `json:"placement"`
	InternshipExperience string  `json:"internship_experience"`
	CollegeID            string  `json:"college_id"`
}

// Fonction pour initialiser la connexion à la base PostgreSQL
func initDB() (*sql.DB, error) {
	connStr := "postgres://postgres:123mont@456@localhost:5432/ma_base?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	// Vérification de la connexion
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

// Handler pour récupérer tous les collèges
func getColleges(c *gin.Context) {
	db, err := initDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur de connexion à la base de données"})
		return
	}
	defer db.Close()

	rows, err := db.Query(`SELECT 
	"Communication_Skills", 
	"Projects_Completed", 
	"Prev_Sem_Result", 
	"CGPA",
	"Academic_Performance",
	"IQ",
	"Extra_Curricular_Score",
	"Placement",
	"Internship_Experience",
	"College_ID"
	FROM college`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des données"})
		return
	}
	defer rows.Close()

	var colleges []College
	for rows.Next() {
		var college College
		// Remplissage de la structure à partir du résultat SQL
		if err := rows.Scan(&college.CommunicationSkills, &college.ProjectsCompleted, &college.PrevSemResult,
			&college.CGPA, &college.AcademicPerformance, &college.IQ,
			&college.ExtraCurricularScore, &college.Placement,
			&college.InternshipExperience, &college.CollegeID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du scan des données"})
			return
		}
		colleges = append(colleges, college)
	}

	// Retour du tableau JSON des collèges
	c.JSON(http.StatusOK, colleges)
}

// Handler pour récupérer un collège via son ID
func getCollegesperCollegeID(c *gin.Context) {
	collegeID := c.Param("college_id")

	db, err := initDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur de connexion à la base de données"})
		return
	}
	defer db.Close()

	row := db.QueryRow(`SELECT 
	"Communication_Skills", 
	"Projects_Completed", 
	"Prev_Sem_Result", 
	"CGPA",
	"Academic_Performance",
	"IQ",
	"Extra_Curricular_Score",
	"Placement",
	"Internship_Experience",
	"College_ID"
	FROM college WHERE "College_ID" = $1`, collegeID)

	var college College
	if err := row.Scan(&college.CommunicationSkills, &college.ProjectsCompleted, &college.PrevSemResult,
		&college.CGPA, &college.AcademicPerformance, &college.IQ,
		&college.ExtraCurricularScore, &college.Placement,
		&college.InternshipExperience, &college.CollegeID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Collège introuvable"})
		return
	}

	c.JSON(http.StatusOK, college)
}

// Middleware pour mesurer le temps d'exécution réel d'une route
func LoggerMiddleware(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now() // Marquer le début

		c.Next() // Continuer vers le handler principal

		duration := time.Since(start) // Mesurer la durée d'exécution

		// Insertion asynchrone pour ne pas bloquer la réponse
		go func() {
			_, err := db.Exec(`
            INSERT INTO statistiques_api (path, method, status_code, duration_micro_s,created_at)
            VALUES ($1, $2, $3, $4, NOW())`,
				c.FullPath(),
				c.Request.Method,
				c.Writer.Status(),
				duration.Microseconds(),
			)

			if err != nil {
				log.Println("Erreur insertion statistique:", err)
			}
		}()
	}
}

// func pour modiefer la structure de la base de données
// 1. Ajouter un nouveau collège
func createCollege(c *gin.Context) {
	var newCollege College

	if err := c.ShouldBindJSON(&newCollege); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format JSON invalide"})
		return
	}

	db, err := initDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur de connexion DB"})
		return
	}
	defer db.Close()

	_, err = db.Exec(`INSERT INTO college (
		"Communication_Skills", "Projects_Completed", "Prev_Sem_Result", "CGPA",
		"Academic_Performance", "IQ", "Extra_Curricular_Score", "Placement",
		"Internship_Experience", "College_ID")
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		newCollege.CommunicationSkills,
		newCollege.ProjectsCompleted,
		newCollege.PrevSemResult,
		newCollege.CGPA,
		newCollege.AcademicPerformance,
		newCollege.IQ,
		newCollege.ExtraCurricularScore,
		newCollege.Placement,
		newCollege.InternshipExperience,
		newCollege.CollegeID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur insertion"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Collège ajouté avec succès"})
}

// 2. Incrémenter Projects_Completed de +1
func incrementProjectsCompleted(c *gin.Context) {
	id := c.Param("college_id")
	db, err := initDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur DB"})
		return
	}
	defer db.Close()

	res, err := db.Exec(`UPDATE college SET "Projects_Completed" = "Projects_Completed" + 1 WHERE "College_ID" = $1`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur update"})
		return
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "ID introuvable"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Projects_Completed incrémenté"})
}

// 3. Toggle "Placement" (yes <-> no)
func togglePlacement(c *gin.Context) {
	id := c.Param("college_id")
	db, err := initDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur DB"})
		return
	}
	defer db.Close()

	_, err = db.Exec(`
		UPDATE college 
		SET "Placement" = CASE 
			WHEN "Placement" = 'yes' THEN 'no'
			ELSE 'yes' 
		END
		WHERE "College_ID" = $1`, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur modification placement"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Placement modifié"})
}

// 4. Supprimer un collège par ID
func deleteCollege(c *gin.Context) {
	id := c.Param("college_id")
	db, err := initDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur DB"})
		return
	}
	defer db.Close()

	res, err := db.Exec(`DELETE FROM college WHERE "College_ID" = $1`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur suppression"})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "ID introuvable"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Collège supprimé"})
}

func main() {
	// Connexion à la base au lancement
	db, err := initDB()
	if err != nil {
		log.Fatal("Erreur connexion DB :", err)
	}

	// Initialisation de Gin
	router := gin.Default()

	// Ajout du middleware Logger
	router.Use(LoggerMiddleware(db))

	// Définition des routes
	router.GET("/college/gin", getColleges)
	router.GET("/college/gin/:college_id", getCollegesperCollegeID)
	router.POST("/college/gin", createCollege)
	router.PUT("/college/gin/:college_id/increment", incrementProjectsCompleted)
	router.PUT("/college/gin/:college_id/toggle-placement", togglePlacement)
	router.DELETE("/college/gin/:college_id", deleteCollege)

	// Lancement du serveur
	log.Println("API démarrée sur http://localhost:8080")
	router.Run(":8080")
}
