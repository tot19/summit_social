package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/tot19/summit_social/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Post model
type Post struct {
	ID            uint   `gorm:"primaryKey;autoIncrement"`
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

// @title Summit Social API
// @version 1.0
// @description This is a social media API server.
// @host localhost:8080
// @BasePath /
func main() {
	// Initialize the database connection
	initDB()

	// Create a new Gin router
	router := gin.Default()

	// Serve the Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Define the API routes
	router.POST("/post", createPost)
	router.GET("/posts", getPosts)
	router.GET("/post/:id", getPost)
	router.PUT("/post/:id", updatePost)
	router.DELETE("/post/:id", deletePost)

	// Start the server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// Create a new post
// @Summary Create a new post
// @Description Create a new post with the input payload
// @Tags posts
// @Accept json
// @Produce json
// @Param post body Post true "Create post"
// @Success 201 {object} Post
// @Failure 400 {object} map[string]interface{}
// @Router /post [post]
func createPost(c *gin.Context) {
	var post Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Create(&post)
	c.JSON(http.StatusCreated, post)
}

// Get all posts
// @Summary Get all posts
// @Description Get a list of all posts
// @Tags posts
// @Produce json
// @Success 200 {array} Post
// @Router /posts [get]
func getPosts(c *gin.Context) {
	var posts []Post
	db.Find(&posts)
	c.JSON(http.StatusOK, posts)
}

// Get a single post by ID
// @Summary Get a post by ID
// @Description Get a post by its ID
// @Tags posts
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} Post
// @Failure 400 {object} map[string]interface{}
// @Router /post/{id} [get]
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
// @Summary Update a post
// @Description Update an existing post by its ID
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Param post body Post true "Update post"
// @Success 200 {object} Post
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object}  map[string]interface{}
// @Router /post/{id} [put]
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
// @Summary Delete a post by ID
// @Description Delete a post by its ID
// @Tags posts
// @Param id path int true "Post ID"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Router /post/{id} [delete]
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
