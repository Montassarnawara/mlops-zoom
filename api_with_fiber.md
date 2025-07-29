package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

// Structure représentant la table college
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

// Connexion à PostgreSQL
func initDB() (*sql.DB, error) {
	connStr := "postgres://postgres:123mont@456@localhost:5432/ma_base?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

// Middleware pour journaliser chaque requête
func logStats(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Appelle le prochain handler
		err := c.Next()

		// Sécurité : évite panic si nil
		statusCode := 0
		if c.Response() != nil {
			statusCode = c.Response().StatusCode()
		}

		duration := time.Since(start).Microseconds()

		go func() {
			_, err := db.Exec(`INSERT INTO statistiques_api (path, method, status_code, duration_micro_s, created_at) 
				VALUES ($1, $2, $3, $4, $5)`,
				c.Path(),
				c.Method(),
				statusCode,
				duration,
				time.Now(),
			)
			if err != nil {
				log.Println("Erreur insertion stats:", err)
			}
		}()

		return err
	}
}


// Endpoint pour récupérer tous les colleges
func getColleges(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
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
			return c.Status(500).JSON(fiber.Map{"error": "Erreur récupération"})
		}
		defer rows.Close()

		var colleges []College
		for rows.Next() {
			var col College
			err := rows.Scan(
				&col.CommunicationSkills,
				&col.ProjectsCompleted,
				&col.PrevSemResult,
				&col.CGPA,
				&col.AcademicPerformance,
				&col.IQ,
				&col.ExtraCurricularScore,
				&col.Placement,
				&col.InternshipExperience,
				&col.CollegeID,
			)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": "Erreur scan"})
			}
			colleges = append(colleges, col)
		}
		return c.JSON(colleges)
	}
}

// Endpoint par College_ID
func getCollegeByID(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("college_id")
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
			FROM college WHERE "College_ID" = $1`, id)

		var col College
		err := row.Scan(
			&col.CommunicationSkills,
			&col.ProjectsCompleted,
			&col.PrevSemResult,
			&col.CGPA,
			&col.AcademicPerformance,
			&col.IQ,
			&col.ExtraCurricularScore,
			&col.Placement,
			&col.InternshipExperience,
			&col.CollegeID,
		)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Non trouvé"})
		}
		return c.JSON(col)
	}
}
// POST /college/fiber – Ajouter un collège
func createCollege(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var col College
		if err := c.BodyParser(&col); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Format JSON invalide"})
		}

		_, err := db.Exec(`INSERT INTO college (
			"Communication_Skills", "Projects_Completed", "Prev_Sem_Result", "CGPA",
			"Academic_Performance", "IQ", "Extra_Curricular_Score", "Placement",
			"Internship_Experience", "College_ID")
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
			col.CommunicationSkills, col.ProjectsCompleted, col.PrevSemResult, col.CGPA,
			col.AcademicPerformance, col.IQ, col.ExtraCurricularScore,
			col.Placement, col.InternshipExperience, col.CollegeID)

		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Erreur insertion"})
		}
		return c.JSON(fiber.Map{"message": "Collège ajouté"})
	}
}

// PUT /college/fiber/:college_id – Mettre à jour un collège entier

func updateCollege(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("college_id")
		var col College
		if err := c.BodyParser(&col); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Format JSON invalide"})
		}

		_, err := db.Exec(`UPDATE college SET 
			"Communication_Skills"=$1, "Projects_Completed"=$2, "Prev_Sem_Result"=$3, "CGPA"=$4,
			"Academic_Performance"=$5, "IQ"=$6, "Extra_Curricular_Score"=$7,
			"Placement"=$8, "Internship_Experience"=$9 WHERE "College_ID"=$10`,
			col.CommunicationSkills, col.ProjectsCompleted, col.PrevSemResult, col.CGPA,
			col.AcademicPerformance, col.IQ, col.ExtraCurricularScore,
			col.Placement, col.InternshipExperience, id)

		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Erreur update"})
		}
		return c.JSON(fiber.Map{"message": "Collège mis à jour"})
	}
}

// DELETE /college/fiber/:college_id – Supprimer un collège

func deleteCollege(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("college_id")

		_, err := db.Exec(`DELETE FROM college WHERE "College_ID" = $1`, id)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Erreur suppression"})
		}
		return c.JSON(fiber.Map{"message": "Collège supprimé"})
	}
}



 //PATCH /college/fiber/:college_id/placement – Modifier uniquement le champ Placement

 func updatePlacement(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("college_id")
		type Req struct {
			Placement string `json:"placement"`
		}
		var req Req
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Format JSON invalide"})
		}

		_, err := db.Exec(`UPDATE college SET "Placement" = $1 WHERE "College_ID" = $2`, req.Placement, id)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Erreur update placement"})
		}
		return c.JSON(fiber.Map{"message": "Placement mis à jour"})
	}
}


func main() {
	db, err := initDB()
	if err != nil {
		log.Fatal("Erreur DB:", err)
	}

	app := fiber.New()

	// Middleware de log personnalisé
	app.Use(logStats(db))

	// Routes
	app.Get("/college/fiber", getColleges(db))
	app.Get("/college/fiber/:college_id", getCollegeByID(db))
	app.Post("/college/fiber", createCollege(db))
	app.Put("/college/fiber/:college_id", updateCollege(db))
	app.Delete("/college/fiber/:college_id", deleteCollege(db))
	app.Patch("/college/fiber/:college_id/placement", updatePlacement(db))

	log.Println("API Fiber lancée sur http://localhost:8080")
	app.Listen("0.0.0.0:8888")
}
