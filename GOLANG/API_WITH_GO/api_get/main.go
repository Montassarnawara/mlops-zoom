package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type data struct {
	Id         string `json:"Id"`
	Name       string `json:"Name"`
	Validation bool   `json:"Validation"`
}

func getData(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, datalist)
}

var datalist = []data{
	{Id: "1", Name: "ali", Validation: true},
	{Id: "2", Name: "ahmed", Validation: false},
	{Id: "3", Name: "mohamed", Validation: true},
	{Id: "4", Name: "sadok", Validation: false},
}

func getdatabyid(id string) (*data, error) {
	for i, t := range datalist {
		if t.Id == id {
			return &datalist[i], nil
		}
	}
	return nil, errors.New("data not found")
}
func getDataById(context *gin.Context) {
	id := context.Param("id")
	data, err := getdatabyid(id)
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "data not found"})
		return
	}
	context.IndentedJSON(http.StatusOK, data)
}

func main() {
	router := gin.Default()
	router.GET("/go_test", getData)
	router.GET("/go_test/:id", getDataById)
	router.Run("localhost:9090")
}
