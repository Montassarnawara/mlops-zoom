package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type Product struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	CreatedAt string  `json:"created_at"` // ici string
	UpdatedAt string  `json:"updated_at"` // ici string
}

func getDBConnection() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./product.db")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func getProducts(c *gin.Context) {
	db, err := getDBConnection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, name, price, created_at, updated_at FROM product")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Query execution failed"})
		return
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Row scan failed"})
			return
		}
		products = append(products, p)
	}

	c.JSON(http.StatusOK, products)
}

func addProduct(c *gin.Context) {
	var p Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Ignore p.ID envoyé par le client, car c'est autoincrement

	// Si CreatedAt vide, on met la date actuelle
	if p.CreatedAt == "" {
		p.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	}
	if p.UpdatedAt == "" {
		p.UpdatedAt = p.CreatedAt
	}

	db, err := getDBConnection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO product(name, price, created_at, updated_at) VALUES(?, ?, ?, ?)")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Statement preparation failed"})
		return
	}
	defer stmt.Close()
	res, err := stmt.Exec(p.Name, p.Price, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Insert operation failed"})
		return
	}

	// Récupérer le dernier ID inséré
	lastID, err := res.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve last insert id"})
		return
	}

	p.ID = int(lastID) // mettre à jour l’ID généré dans la réponse

	c.JSON(http.StatusCreated, p) // retourner le produit inséré avec ID
}

func getproductperid(c *gin.Context) {
	id := c.Param("id")

	db, err := getDBConnection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
		return
	}
	defer db.Close()

	var p Product
	err = db.QueryRow("SELECT id, name, price, created_at, updated_at FROM product WHERE id = ?", id).Scan(&p.ID, &p.Name, &p.Price, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Query execution failed"})
		}
		return
	}

	c.JSON(http.StatusOK, p)
}

func getproductperName(c *gin.Context) {
	name := c.Param("name")

	db, err := getDBConnection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
		return
	}
	defer db.Close()

	var p Product
	err = db.QueryRow("SELECT id, name, price, created_at, updated_at FROM product WHERE name = ?", name).Scan(&p.ID, &p.Name, &p.Price, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Query execution failed"})
		}
		return
	}

	c.JSON(http.StatusOK, p)
}

func delproductperID(c *gin.Context) {
	id := c.Param("id")

	db, err := getDBConnection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM product WHERE id = ?")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Statement preparation failed"})
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Delete operation failed"})
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve rows affected"})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

func main() {
	router := gin.Default()
	router.GET("/products", getProducts)
	router.GET("/products/:id", getproductperid)
	router.GET("/products/name/:name", getproductperName)
	router.DELETE("/products/:id", delproductperID)
	router.POST("/products", addProduct)
	router.Run(":8081")
}
