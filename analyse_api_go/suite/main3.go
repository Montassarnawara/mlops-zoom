package main

import (
	"context"
	"encoding/json"
	"log"
	"math"
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hazelcast/hazelcast-go-client"
)

type College struct {
	CollegeID            string  `json:"college_id"`
	IQ                   int64   `json:"iq"`
	PrevSemResult        float64 `json:"prev_sem_result"`
	CGPA                 float64 `json:"cgpa"`
	AcademicPerformance  int64   `json:"academic_performance"`
	InternshipExperience string  `json:"internship_experience"`
	ExtraCurricularScore int64   `json:"extra_curricular_score"`
	CommunicationSkills  int64   `json:"communication_skills"`
	ProjectsCompleted    int64   `json:"projects_completed"`
	Placement            string  `json:"placement"`
}

var (
	Ctx      = context.Background()
	Client   *hazelcast.Client
	CacheMap *hazelcast.Map
)

func InitHazelcast() {
	config := hazelcast.Config{}
	config.Cluster.Network.SetAddresses("127.0.0.1:5701")
	config.Cluster.Name = "dev-go"

	var err error
	Client, err = hazelcast.StartNewClientWithConfig(Ctx, config)
	if err != nil {
		log.Fatalf(" Hazelcast Connexion Failed: %v", err)
	}
	CacheMap, err = Client.GetMap(Ctx, "college_cache")
	if err != nil {
		log.Fatalf(" GetMap Failed: %v", err)
	}
	log.Println(" Hazelcast Connected")
}

func AnalyseCollegeHandler(c *gin.Context) {
	id := c.Param("id")
	college, err := GetFromCache(id)
	if err != nil || college == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "College not found"})
		return
	}

	all, err := CacheMap.GetEntrySet(Ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load entries"})
		return
	}

	var allIQs []float64
	for _, kv := range all {
		var col College
		if b, ok := kv.Value.([]byte); ok {
			if err := json.Unmarshal(b, &col); err == nil {
				if strings.EqualFold(col.CollegeID, id) {
					continue // exclude current ID
				}
				allIQs = append(allIQs, float64(col.IQ)) // Convert IQ to float64
			}
		}
	}
	sort.Float64s(allIQs)

	iq := float64(college.IQ) // ← convert once here

	lowIQs := padWithZeros(filterLower(allIQs, iq), 8)
	highIQs := padWithZeros(filterHigher(allIQs, iq), 8)

	A := buildMatrix(highIQs, iq)
	B := buildMatrix(lowIQs, iq)
	detA := Determinant3x3(A)
	detB := Determinant3x3(B)
	diff := math.Abs(detA - detB)

	status := "OK"
	if diff < iq {
		status = "ALERTE: IQ plus fort que différence"
	}

	c.JSON(http.StatusOK, gin.H{
		"College":    college.CollegeID,
		"IQ":         college.IQ,
		"Det_A":      detA,
		"Det_B":      detB,
		"Difference": diff,
		"Status":     status,
	})
}

func GetFromCache(id string) (*College, error) {
	raw, err := CacheMap.Get(Ctx, id)
	if err != nil || raw == nil {
		return nil, err
	}

	var col College
	switch v := raw.(type) {
	case []byte:
		err = json.Unmarshal(v, &col)
		return &col, err
	case College:
		return &v, nil
	default:
		return nil, nil
	}
}

func filterLower(data []float64, ref float64) []float64 {
	var out []float64
	for _, v := range data {
		if v < ref {
			out = append(out, v)
		}
	}
	return out
}

func filterHigher(data []float64, ref float64) []float64 {
	var out []float64
	for _, v := range data {
		if v > ref {
			out = append(out, v)
		}
	}
	return out
}

func padWithZeros(arr []float64, n int) []float64 {
	for len(arr) < n {
		arr = append(arr, 0.0)
	}
	return arr[:n]
}

func buildMatrix(vals []float64, center float64) [][]float64 {
	return [][]float64{
		{vals[0], vals[1], vals[2]},
		{vals[3], center, vals[4]},
		{vals[5], vals[6], vals[7]},
	}
}

func Determinant3x3(m [][]float64) float64 {
	return m[0][0]*(m[1][1]*m[2][2]-m[1][2]*m[2][1]) -
		m[0][1]*(m[1][0]*m[2][2]-m[1][2]*m[2][0]) +
		m[0][2]*(m[1][0]*m[2][1]-m[1][1]*m[2][0])
}

func main3() {
	InitHazelcast()
	r := gin.Default()
	r.GET("/college/analyse/:id", AnalyseCollegeHandler)
	r.Run(":8080")
}
