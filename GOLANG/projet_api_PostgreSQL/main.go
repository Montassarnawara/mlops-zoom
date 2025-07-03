package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var db *sql.DB

func initDB() {
	connStr := "host=localhost port=5433 user=myuser password=mypassword dbname=products_db sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Erreur de connexion à la base :", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Impossible de ping la base :", err)
	}

	fmt.Println("Connexion à PostgreSQL réussie !")
}

type Product struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

func getProducts(c *gin.Context) {
	rows, err := db.Query("SELECT id, name, price, created_at, updated_at FROM products")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	products := []Product{}

	for rows.Next() {
		var p Product
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		products = append(products, p)
	}

	c.JSON(http.StatusOK, products)
}

func getProductByID(c *gin.Context) {
	id := c.Param("id")
	var p Product

	query := "SELECT id, name, price, created_at, updated_at FROM products WHERE id = $1"
	err := db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Price, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, p)
}

func addProduct(c *gin.Context) {
	var p Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	query := "INSERT INTO products (name, price, created_at, updated_at) VALUES ($1, $2, NOW(), NOW()) RETURNING id"
	err := db.QueryRow(query, p.Name, p.Price).Scan(&p.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, p)
}

func modifyProduct(c *gin.Context) {
	var p Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	query := "UPDATE products SET name = $1, price = $2, updated_at = NOW() WHERE id = $3"
	_, err := db.Exec(query, p.Name, p.Price, p.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
}

func delproductperID(c *gin.Context) {
	id := c.Param("id")

	query := "DELETE FROM products WHERE id = $1"
	result, err := db.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
func main() {
	initDB()
	defer db.Close()
	router := gin.Default()
	router.GET("/products", getProducts)
	router.GET("/products/:id", getProductByID)
	router.POST("/products", addProduct)
	router.PUT("/products/:id", modifyProduct)
	router.DELETE("/products/:id", delproductperID)
	router.Run(":8081")
}
