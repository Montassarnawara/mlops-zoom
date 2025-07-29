package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

// Structure College
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

// Connexion PostgreSQL
func initDB() (*sql.DB, error) {
	connStr := "postgres://postgres:123mont%40456@localhost:5432/ma_base?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

// Middleware de journalisation Echo
func LoggerMiddleware(db *sql.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			duration := time.Since(start)

			_, errInsert := db.Exec(`INSERT INTO statistiques_api (path, method, status_code, duration_micro_s, created_at)
				VALUES ($1, $2, $3, $4, NOW())`,
				c.Request().URL.Path,
				c.Request().Method,
				c.Response().Status,
				duration.Microseconds())

			if errInsert != nil {
				log.Printf("Erreur insertion stats: %v", errInsert)
			}
			return err
		}
	}
}

// GET /college
func getColleges(c echo.Context) error {
	db, err := initDB()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erreur DB"})
	}
	defer db.Close()

	rows, err := db.Query(`SELECT "Communication_Skills", "Projects_Completed", "Prev_Sem_Result", "CGPA",
		"Academic_Performance", "IQ", "Extra_Curricular_Score", "Placement",
		"Internship_Experience", "College_ID" FROM college`)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erreur SELECT"})
	}
	defer rows.Close()

	var colleges []College
	for rows.Next() {
		var clg College
		err := rows.Scan(&clg.CommunicationSkills, &clg.ProjectsCompleted, &clg.PrevSemResult,
			&clg.CGPA, &clg.AcademicPerformance, &clg.IQ,
			&clg.ExtraCurricularScore, &clg.Placement,
			&clg.InternshipExperience, &clg.CollegeID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erreur scan"})
		}
		colleges = append(colleges, clg)
	}
	return c.JSON(http.StatusOK, colleges)
}

// GET /college/:college_id
func getCollegeByID(c echo.Context) error {
	collegeID := c.Param("college_id")
	db, err := initDB()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erreur DB"})
	}
	defer db.Close()

	row := db.QueryRow(`SELECT "Communication_Skills", "Projects_Completed", "Prev_Sem_Result", "CGPA",
		"Academic_Performance", "IQ", "Extra_Curricular_Score", "Placement",
		"Internship_Experience", "College_ID" FROM college WHERE "College_ID" = $1`, collegeID)

	var clg College
	err = row.Scan(&clg.CommunicationSkills, &clg.ProjectsCompleted, &clg.PrevSemResult,
		&clg.CGPA, &clg.AcademicPerformance, &clg.IQ,
		&clg.ExtraCurricularScore, &clg.Placement,
		&clg.InternshipExperience, &clg.CollegeID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erreur lecture par ID"})
	}
	return c.JSON(http.StatusOK, clg)
}
// Ajouter un nouveau collège (POST /college/echo)
func createCollege(c echo.Context) error {
	var clg College
	if err := c.Bind(&clg); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "JSON invalide"})
	}

	db, err := initDB()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erreur DB"})
	}
	defer db.Close()

	_, err = db.Exec(`INSERT INTO college (
		"Communication_Skills", "Projects_Completed", "Prev_Sem_Result", "CGPA",
		"Academic_Performance", "IQ", "Extra_Curricular_Score", "Placement",
		"Internship_Experience", "College_ID")
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		clg.CommunicationSkills, clg.ProjectsCompleted, clg.PrevSemResult, clg.CGPA,
		clg.AcademicPerformance, clg.IQ, clg.ExtraCurricularScore, clg.Placement,
		clg.InternshipExperience, clg.CollegeID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erreur insertion"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Collège ajouté"})
}
// Incrémenter Projects_Completed (PUT /college/echo/:college_id/increment)
func incrementProjects(c echo.Context) error {
	id := c.Param("college_id")
	db, err := initDB()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erreur DB"})
	}
	defer db.Close()

	res, err := db.Exec(`UPDATE college SET "Projects_Completed" = "Projects_Completed" + 1 WHERE "College_ID" = $1`, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erreur mise à jour"})
	}

	count, _ := res.RowsAffected()
	if count == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "ID non trouvé"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Projects_Completed incrémenté"})
}
//Inverser Placement (PUT /college/echo/:college_id/toggle-placement)

func togglePlacement(c echo.Context) error {
	id := c.Param("college_id")
	db, err := initDB()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erreur DB"})
	}
	defer db.Close()

	_, err = db.Exec(`UPDATE college SET "Placement" = CASE WHEN "Placement" = 'yes' THEN 'no' ELSE 'yes' END WHERE "College_ID" = $1`, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erreur update Placement"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Placement mis à jour"})
}

//Supprimer un collège (DELETE /college/echo/:college_id)

func deleteCollege(c echo.Context) error {
	id := c.Param("college_id")
	db, err := initDB()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erreur DB"})
	}
	defer db.Close()

	res, err := db.Exec(`DELETE FROM college WHERE "College_ID" = $1`, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Erreur suppression"})
	}

	count, _ := res.RowsAffected()
	if count == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "ID non trouvé"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Collège supprimé"})
}


func main() {
	db, err := initDB()
	if err != nil {
		log.Fatal("Erreur connexion DB:", err)
	}

	e := echo.New()
	e.Use(LoggerMiddleware(db))

	e.GET("/college/echo", getColleges)
	e.GET("/college/echo/:college_id", getCollegeByID)
	e.POST("/college/echo", createCollege)
	e.PUT("/college/echo/:college_id/increment", incrementProjects)
	e.PUT("/college/echo/:college_id/toggle-placement", togglePlacement)
	e.DELETE("/college/echo/:college_id", deleteCollege)


	log.Println("API Echo lancée sur http://localhost:8080")
	e.Logger.Fatal(e.Start("0.0.0.0:8888"))
}
