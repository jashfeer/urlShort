package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"sync"

	"github.com/gin-gonic/gin"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

var urlStore = make(map[string]string)
var mu sync.Mutex

// Generate a random string of fixed length
func generateShortURL() string {
	b := make([]rune, 6)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// Check if URL is valid
func isValidURL(testURL string) bool {
	_, err := url.ParseRequestURI(testURL)
	return err == nil
}

func main() {
	// Create Gin router
	r := gin.Default()

	// POST endpoint to shorten a long URL
	r.POST("/shorten", func(c *gin.Context) {
		var requestBody struct {
			URL string `json:"url"`
		}

		// Bind request JSON to struct
		if err := c.ShouldBindJSON(&requestBody); err != nil || !isValidURL(requestBody.URL) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL"})
			return
		}

		// Generate a short URL
		shortURL := generateShortURL()

		// Ensure thread-safe access to the map
		mu.Lock()
		urlStore[shortURL] = requestBody.URL
		fmt.Println("Stored URLs after POST:", urlStore) 
		mu.Unlock()

		// Respond with the short URL
		c.JSON(http.StatusOK, gin.H{
			"short_url": fmt.Sprintf("http://localhost:8080/%s", shortURL),
		})
	})

	// GET endpoint to redirect from short URL to original URL
	r.GET("/:shortURL", func(c *gin.Context) {
		shortURL := c.Param("shortURL")
		fmt.Println("Requested short URL:", shortURL)       
		fmt.Println("Stored URLs at GET request:", urlStore) 

		mu.Lock()
		originalURL, exists := urlStore[shortURL]
		mu.Unlock()

		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
			return
		}

		c.Redirect(http.StatusFound, originalURL)
	})

	// Start the server on port 8080
	r.Run(":8080")
}
