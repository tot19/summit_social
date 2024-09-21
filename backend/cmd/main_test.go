package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Test the ping route
func TestPingRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"message\":\"pong\"}", w.Body.String())
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
	json.Unmarshal(w.Body.Bytes(), &createdPost)
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
	db.Create(&post)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/posts/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var fetchedPost Post
	json.Unmarshal(w.Body.Bytes(), &fetchedPost)
	assert.Equal(t, post.Content, fetchedPost.Content)
	assert.Equal(t, post.Image, fetchedPost.Image)
	assert.Equal(t, post.Poster, fetchedPost.Poster)
}

// Test updating a post
func TestUpdatePost(t *testing.T) {
	router := setupRouter()

	// Create a test post first
	post := Post{
		Content: "Test Post for Update",
		Image:   "/images/test-update.png",
		Poster:  "tester-update",
	}
	db.Create(&post)

	updatedPost := Post{
		Content: "Updated Content",
		Image:   "/images/updated.png",
		Poster:  "tester-update",
	}

	jsonValue, _ := json.Marshal(updatedPost)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/posts/1", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var fetchedPost Post
	json.Unmarshal(w.Body.Bytes(), &fetchedPost)
	assert.Equal(t, updatedPost.Content, fetchedPost.Content)
	assert.Equal(t, updatedPost.Image, fetchedPost.Image)
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
	db.Create(&post)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/posts/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 204, w.Code)
}

// Utility function to set up the Gin router with routes for testing
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
