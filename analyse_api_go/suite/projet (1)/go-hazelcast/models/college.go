package models

// College représente un enregistrement de la table "college"
type College struct {
	CollegeID             string  `json:"college_id"`
	IQ                    int64   `json:"iq"`
	PrevSemResult         float64 `json:"prev_sem_result"`
	CGPA                  float64 `json:"cgpa"`
	AcademicPerformance   int64   `json:"academic_performance"`
	InternshipExperience  string  `json:"internship_experience"`
	ExtraCurricularScore  int64   `json:"extra_curricular_score"`
	CommunicationSkills   int64   `json:"communication_skills"`
	ProjectsCompleted     int64   `json:"projects_completed"`
	Placement             string  `json:"placement"`
}
