package db

import (
	"database/sql"
	"log"

	"go/hazelcast/api/models" //   chemin d'import

	_ "github.com/lib/pq"
)

var DB *sql.DB

// Initialiser la connexion à la base PostgreSQL
func InitDB() {
	//  `%40` est le code URL de `@`, donc ton mot de passe est bien : 123mont@456
	connStr := "postgres://montassar:123mont%40456@localhost:5432/ma_base?sslmode=disable"

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Connexion PostgreSQL échouée : %v", err)
	}

	if err := DB.Ping(); err != nil {
		log.Fatalf("PostgreSQL ne répond pas : %v", err)
	}

	log.Println("Connexion PostgreSQL établie")
}

// Récupérer un college par ID depuis PostgreSQL
func GetCollegeByID(id string) (*models.College, error) {
	query := `
		SELECT College_ID, IQ, Prev_Sem_Result, CGPA, Academic_Performance,
		       Internship_Experience, Extra_Curricular_Score,
		       Communication_Skills, Projects_Completed, Placement
		FROM college
		WHERE College_ID = $1
	`

	row := DB.QueryRow(query, id)

	var college models.College
	err := row.Scan(
		&college.CollegeID,
		&college.IQ,
		&college.PrevSemResult,
		&college.CGPA,
		&college.AcademicPerformance,
		&college.InternshipExperience,
		&college.ExtraCurricularScore,
		&college.CommunicationSkills,
		&college.ProjectsCompleted,
		&college.Placement,
	)

	if err == sql.ErrNoRows {
		return nil, nil // Pas trouvé
	}
	if err != nil {
		return nil, err
	}

	return &college, nil
}
