package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type user struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Validation bool   `json:"validation"`
}

var users = []user{
	{ID: "1", Name: "ali", Validation: true},
	{ID: "2", Name: "amin", Validation: false},
	{ID: "3", Name: "chadya", Validation: true},
}

func getUsers(Context *gin.Context) {
	Context.IndentedJSON(http.StatusOK, users)

}

func getuserId(id string) (*user, bool) {
	for i, u := range users {
		if u.ID == id {
			return &users[i], true
		}
	}
	return nil, false
}

func getuserperId(context *gin.Context) {
	id := context.Param("i")
	userId, found := getuserId(id)
	if !found {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	}
	context.IndentedJSON(http.StatusOK, userId)
}

func adduser(context *gin.Context) {
	var newuser user
	if err := context.BindJSON(&newuser); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	users = append(users, newuser)
	context.IndentedJSON(http.StatusCreated, newuser)

}

func patchuserperID(context *gin.Context) {
	id := context.Param("i")
	userId, found := getuserId(id)
	if !found {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	}
	var updateuser user
	if err := context.BindJSON(&updateuser); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	userId.Name = updateuser.Name
	userId.Validation = updateuser.Validation
	context.IndentedJSON(http.StatusOK, userId)

}
func delfirstuser(context *gin.Context) {
	users = append(users[:0], users[1:]...)
}

func checkHeaders(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	if auth == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
		c.Abort()
		return
	}
	// Répond juste avec OK, sans corps
	c.Status(http.StatusOK)
}

func main() {
	router := gin.Default()
	router.GET("/users", getUsers)
	router.GET("/users/:i", getuserperId)
	router.POST("/users", adduser)
	router.PATCH("/users/:i", patchuserperID)
	router.DELETE("/users", delfirstuser)
	router.HEAD("/users", checkHeaders)
	router.Run("localhost:8888")

}
