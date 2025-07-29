package models

// College represents the structure of a college entity.
// It includes fields such as CollegeID, IQ, CGPA, MinMarks, Attendance,
// Backlogs, Communication, Logical, Projects, and Certifications.
type College struct {
	CollegeID      string
	IQ             int
	CGPA           float64
	MinMarks       float64
	Attendance     int
	Backlogs       string
	Communication  int
	Logical        int
	Projects       int
	Certifications string
}
