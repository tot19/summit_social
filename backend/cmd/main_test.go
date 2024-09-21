package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// SetupRouter initializes the Gin router with routes
func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/posts", createPost)
	router.GET("/posts/:id", getPost)
	router.PUT("/posts/:id", updatePost)
	router.DELETE("/posts/:id", deletePost)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	return router
}

// TestMain sets up the environment before running tests
func TestMain(m *testing.M) {
	// Initialize the database
	initDB()

	// Run the tests
	code := m.Run()

	// Exit with the appropriate code
	os.Exit(code)
}

// Test the ping route
func TestPingRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, `{"message":"pong"}`, w.Body.String())
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
	req, _ := http.NewRequest("POST", "/posts", bytes.NewBuffer(jsonValue))
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
	req, _ := http.NewRequest("GET", fmt.Sprintf("/posts/%d", post.ID), nil)
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

	req, _ := http.NewRequest("PUT", fmt.Sprintf("/posts/%d", post.ID), bytes.NewBuffer(updatedPostJSON))
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
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/posts/%d", post.ID), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 204, w.Code)

	// Verify deletion
	var deletedPost Post
	err := db.First(&deletedPost, post.ID).Error
	assert.Error(t, err)
}
