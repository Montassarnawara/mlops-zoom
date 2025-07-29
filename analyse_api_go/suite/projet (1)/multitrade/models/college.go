// File: projet/multitrade/models/college.go
package models

type College struct {
	CollegeID            string  `json:"College_ID"`
	IQ                   int     `json:"IQ"`
	PrevSemResult        float64 `json:"Prev_Sem_Result"`
	CGPA                 float64 `json:"CGPA"`
	AcademicPerformance  int     `json:"Academic_Performance"`
	InternshipExperience string  `json:"Internship_Experience"`
	ExtraCurricularScore int     `json:"Extra_Curricular_Score"`
	CommunicationSkills  int     `json:"Communication_Skills"`
	ProjectsCompleted    int     `json:"Projects_Completed"`
	Placement            string  `json:"Placement"`
}
