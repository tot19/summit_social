package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Post model
type Post struct {
	ID            uint   `gorm:"primaryKey"`
	Date          string `gorm:"default:CURRENT_TIMESTAMP"`
	Content       string
	Image         string
	Likes         int `gorm:"default:0"`
	Poster        string
	CommentsCount int `gorm:"default:0"`
}

var db *gorm.DB

func initDB() {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	// Migrate the schema
	if err := db.AutoMigrate(&Post{}); err != nil {
		log.Fatal("Failed to auto migrate database schema:", err)
	}
}

func main() {
	// Initialize the database connection
	initDB()

	// Create a new Gin router
	router := gin.Default()

	// Define the API routes
	router.POST("/posts", createPost)
	router.GET("/posts/:id", getPost)
	router.PUT("/posts/:id", updatePost)
	router.DELETE("/posts/:id", deletePost)

	// Start the server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// Create a new post
func createPost(c *gin.Context) {
	var post Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Create(&post)
	c.JSON(http.StatusCreated, post)
}

// Get a single post by ID
func getPost(c *gin.Context) {
	var post Post
	id := c.Param("id")

	if err := db.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(http.StatusOK, post)
}

// Update an existing post
func updatePost(c *gin.Context) {
	var post Post
	id := c.Param("id")
	if err := db.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Model(&post).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}

	c.JSON(http.StatusOK, post)
}

// Delete a post by ID
func deletePost(c *gin.Context) {
	var post Post
	id := c.Param("id")

	if err := db.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	db.Delete(&post)
	c.JSON(http.StatusNoContent, nil)
}
