package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func initTestDB() {
	var err error
	db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the in-memory database:", err)
	}

	// Migrate the schema
	if err := db.AutoMigrate(&Post{}); err != nil {
		log.Fatal("Failed to auto migrate database schema:", err)
	}
}

// SetupRouter initializes the Gin router with routes
func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/posts", getPosts)
	router.POST("/post", createPost)
	router.GET("/post/:id", getPost)
	router.PUT("/post/:id", updatePost)
	router.DELETE("/post/:id", deletePost)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	return router
}

// TestMain sets up the environment before running tests
func TestMain(m *testing.M) {
	// Initialize the test database
	initTestDB()

	seedTestData()

	// Run tests
	code := m.Run()

	os.Exit(code)
}

func seedTestData() {
	// Create common test data here
	db.Create(&Post{
		Content: "Initial Test Post",
		Image:   "/images/init.png",
		Poster:  "init-tester",
	})
}

// Test creating a post
func TestCreatePost(t *testing.T) {
	router := setupRouter()

	post := Post{
		Content: "Test Post",
		Image:   "/images/test.png",
		Poster:  "tester",
	}

	jsonValue, _ := json.Marshal(post)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/post", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)

	var createdPost Post
	err := json.Unmarshal(w.Body.Bytes(), &createdPost)
	assert.NoError(t, err)
	assert.Equal(t, post.Content, createdPost.Content)
	assert.Equal(t, post.Image, createdPost.Image)
	assert.Equal(t, post.Poster, createdPost.Poster)
	assert.Equal(t, 0, createdPost.Likes)
	assert.Equal(t, 0, createdPost.CommentsCount)
}

// Test getting all posts
func TestGetPosts(t *testing.T) {
	router := setupRouter()

	// Optionally, create a post to ensure the database is not empty
	post := Post{
		Content: "Sample Post",
		Image:   "/images/sample.png",
		Poster:  "tester",
	}
	result := db.Create(&post)
	assert.NoError(t, result.Error)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/posts", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var posts []Post
	err := json.Unmarshal(w.Body.Bytes(), &posts)
	assert.NoError(t, err)
	assert.NotEmpty(t, posts)
}

// Test getting a post by ID
func TestGetPost(t *testing.T) {
	router := setupRouter()

	// Create a test post first
	post := Post{
		Content: "Test Post for Get",
		Image:   "/images/test-get.png",
		Poster:  "tester-get",
	}
	result := db.Create(&post)
	assert.NoError(t, result.Error)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/post/%d", post.ID), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var fetchedPost Post
	err := json.Unmarshal(w.Body.Bytes(), &fetchedPost)
	assert.NoError(t, err)
	assert.Equal(t, post.Content, fetchedPost.Content)
	assert.Equal(t, post.Image, fetchedPost.Image)
	assert.Equal(t, post.Poster, fetchedPost.Poster)
}

// Test updating a post
func TestUpdatePost(t *testing.T) {
	router := setupRouter()
	post := Post{
		Content: "Initial Content",
		Image:   "/images/test.png",
		Poster:  "tester",
	}
	result := db.Create(&post)
	assert.NoError(t, result.Error)

	updatedPost := map[string]interface{}{
		"likes": 5,
	}

	updatedPostJSON, _ := json.Marshal(updatedPost)

	req, _ := http.NewRequest("PUT", fmt.Sprintf("/post/%d", post.ID), bytes.NewBuffer(updatedPostJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var updatedPostFromDB Post
	err := db.First(&updatedPostFromDB, post.ID).Error
	assert.NoError(t, err)

	assert.Equal(t, 5, updatedPostFromDB.Likes)
}

// Test deleting a post
func TestDeletePost(t *testing.T) {
	router := setupRouter()

	// Create a test post first
	post := Post{
		Content: "Test Post for Delete",
		Image:   "/images/test-delete.png",
		Poster:  "tester-delete",
	}
	result := db.Create(&post)
	assert.NoError(t, result.Error)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/post/%d", post.ID), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 204, w.Code)

	// Verify deletion
	var deletedPost Post
	err := db.First(&deletedPost, post.ID).Error
	assert.Error(t, err)
}
