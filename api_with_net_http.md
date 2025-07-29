package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

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

var db *sql.DB

func initDB() *sql.DB {
	connStr := "postgres://postgres:123mont@456@localhost:5432/ma_base?sslmode=disable"
	database, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Erreur connexion DB :", err)
	}
	if err := database.Ping(); err != nil {
		log.Fatal("Ping DB échoué :", err)
	}
	return database
}

func logStats(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h(w, r)
		duration := time.Since(start).Microseconds()

		go func() {
			_, err := db.Exec(`INSERT INTO statistiques_api (path, method, status_code, duration_micro_s, created_at)
				VALUES ($1, $2, $3, $4, $5)`,
				r.URL.Path,
				r.Method,
				200, // Tu peux personnaliser si besoin d’un vrai status
				duration,
				time.Now(),
			)
			if err != nil {
				log.Println("Erreur insertion stats:", err)
			}
		}()
	}
}

func getColleges(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`SELECT 
	"Communication_Skills", "Projects_Completed", "Prev_Sem_Result", "CGPA",
	"Academic_Performance", "IQ", "Extra_Curricular_Score", "Placement",
	"Internship_Experience", "College_ID" FROM college`)
	if err != nil {
		http.Error(w, "Erreur récupération", 500)
		return
	}
	defer rows.Close()

	var colleges []College
	for rows.Next() {
		var c College
		if err := rows.Scan(
			&c.CommunicationSkills,
			&c.ProjectsCompleted,
			&c.PrevSemResult,
			&c.CGPA,
			&c.AcademicPerformance,
			&c.IQ,
			&c.ExtraCurricularScore,
			&c.Placement,
			&c.InternshipExperience,
			&c.CollegeID); err != nil {
			http.Error(w, "Erreur scan", 500)
			return
		}
		colleges = append(colleges, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(colleges)
}

func getCollegeByID(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "ID manquant", 400)
		return
	}
	collegeID := parts[3]

	row := db.QueryRow(`SELECT 
	"Communication_Skills", "Projects_Completed", "Prev_Sem_Result", "CGPA",
	"Academic_Performance", "IQ", "Extra_Curricular_Score", "Placement",
	"Internship_Experience", "College_ID" 
	FROM college WHERE "College_ID" = $1`, collegeID)

	var c College
	if err := row.Scan(
		&c.CommunicationSkills,
		&c.ProjectsCompleted,
		&c.PrevSemResult,
		&c.CGPA,
		&c.AcademicPerformance,
		&c.IQ,
		&c.ExtraCurricularScore,
		&c.Placement,
		&c.InternshipExperience,
		&c.CollegeID); err != nil {
		http.Error(w, "Collège non trouvé", 404)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

func createCollege(w http.ResponseWriter, r *http.Request) {
	var c College
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "Données invalides", http.StatusBadRequest)
		return
	}
	_, err := db.Exec(`INSERT INTO college (
		"Communication_Skills", "Projects_Completed", "Prev_Sem_Result", "CGPA",
		"Academic_Performance", "IQ", "Extra_Curricular_Score", "Placement",
		"Internship_Experience", "College_ID") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		c.CommunicationSkills, c.ProjectsCompleted, c.PrevSemResult, c.CGPA,
		c.AcademicPerformance, c.IQ, c.ExtraCurricularScore, c.Placement,
		c.InternshipExperience, c.CollegeID)
	if err != nil {
		http.Error(w, "Erreur insertion", 500)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)
}

func updateCollege(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "ID manquant", 400)
		return
	}
	collegeID := parts[3]
	var c College
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "Données invalides", http.StatusBadRequest)
		return
	}
	_, err := db.Exec(`UPDATE college SET 
		"Communication_Skills"=$1, "Projects_Completed"=$2, "Prev_Sem_Result"=$3, "CGPA"=$4,
		"Academic_Performance"=$5, "IQ"=$6, "Extra_Curricular_Score"=$7, "Placement"=$8,
		"Internship_Experience"=$9 WHERE "College_ID"=$10`,
		c.CommunicationSkills, c.ProjectsCompleted, c.PrevSemResult, c.CGPA,
		c.AcademicPerformance, c.IQ, c.ExtraCurricularScore, c.Placement,
		c.InternshipExperience, collegeID)
	if err != nil {
		http.Error(w, "Erreur mise à jour", 500)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(c)
}

func deleteCollege(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "ID manquant", 400)
		return
	}
	collegeID := parts[3]
	_, err := db.Exec(`DELETE FROM college WHERE "College_ID" = $1`, collegeID)
	if err != nil {
		http.Error(w, "Erreur suppression", 500)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	db = initDB()
	log.Println("API net/http démarrée sur http://localhost:8080")

	http.HandleFunc("/college/net", logStats(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getColleges(w, r)
		case http.MethodPost:
			createCollege(w, r)
		default:
			http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/college/net/", logStats(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getCollegeByID(w, r)
		case http.MethodPut:
			updateCollege(w, r)
		case http.MethodDelete:
			deleteCollege(w, r)
		default:
			http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		}
	}))

	http.ListenAndServe("0.0.0.0:8888", nil)
}