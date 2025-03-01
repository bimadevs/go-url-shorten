/*
Dibuat oleh Bimadev
Source Code ini open source untuk orang.
Jangan hapus ini untuk menghargai Developer

Happy Coding :_)
*/

package main

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB instance
var db *gorm.DB

// Struktur tabel URL
type URL struct {
	ID          uint   `gorm:"primaryKey"`
	ShortCode   string `gorm:"uniqueIndex"`
	OriginalURL string
	ClickCount  uint `gorm:"default:0"`
	APIKey      string
}

// Struktur tabel User API Key
type User struct {
	ID     uint   `gorm:"primaryKey"`
	APIKey string `gorm:"uniqueIndex"`
}

// Regex validasi alias dan URL
var (
	aliasRegex = regexp.MustCompile(`^[a-zA-Z0-9]{1,15}$`)
	urlRegex   = regexp.MustCompile(`^(https?://)?([\w\-]+\.)+[\w\-]+(/[^\s]*)?$`)
)

// Inisialisasi database
func initDB() {
	var err error
	db, err = gorm.Open(sqlite.Open("shortener.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	db.AutoMigrate(&URL{}, &User{})
}

// Generate random API Key
func generateAPIKey() (string, error) {
	b := make([]byte, 10)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	randomPart := base64.RawURLEncoding.EncodeToString(b)
	return "Bimadev" + randomPart, nil
}

// Generate random short URL
func generateShortURL() (string, error) {
	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b)[:6], nil
}

// CORS Middleware untuk penggunaaan FE seperti react
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-API-Key")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

// Function untuk generate API Key
func generateAPIKeyHandler(c *gin.Context) {
	apiKey, err := generateAPIKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate API key"})
		return
	}

	user := User{APIKey: apiKey}
	db.Create(&user)

	c.JSON(http.StatusOK, gin.H{"api_key": apiKey})
}

// Function Shorten URL
func shortenURL(c *gin.Context) {
	apiKey := c.GetHeader("X-API-Key")
	if apiKey == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "API key is required"})
		return
	}

	// Cek apakah API Key valid
	var user User
	if err := db.Where("api_key = ?", apiKey).First(&user).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid API key"})
		return
	}

	var req struct {
		OriginalURL string `json:"original_url"`
		CustomAlias string `json:"custom_alias"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Validasi URL
	if !urlRegex.MatchString(req.OriginalURL) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL format"})
		return
	}

	// Validasi alias jika diberikan
	if req.CustomAlias != "" {
		if !aliasRegex.MatchString(req.CustomAlias) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid alias. Hanya huruf dan angka yang diizinkan (max 15 chars)."})
			return
		}

		// Cek apakah alias sudah dipakai
		var existing URL
		if err := db.Where("short_code = ?", req.CustomAlias).First(&existing).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Custom alias already taken"})
			return
		}
	}

	// Jika tidak ada custom alias, buat yang random
	shortCode := req.CustomAlias
	if shortCode == "" {
		var err error
		shortCode, err = generateShortURL()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate short URL"})
			return
		}
	}

	url := URL{ShortCode: shortCode, OriginalURL: req.OriginalURL, ClickCount: 0, APIKey: apiKey}
	db.Create(&url)

	c.JSON(http.StatusOK, gin.H{"short_url": "http://localhost:8080/" + shortCode})
}

// Redirect handler + Tambah jumlah klik
func redirectURL(c *gin.Context) {
	shortCode := c.Param("short")

	var url URL
	result := db.Where("short_code = ?", shortCode).First(&url)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	db.Model(&url).Update("ClickCount", gorm.Expr("ClickCount + ?", 1))

	c.Redirect(http.StatusFound, url.OriginalURL)
}

// Statistik klik
func getStats(c *gin.Context) {
	shortCode := c.Param("short")

	var url URL
	result := db.Where("short_code = ?", shortCode).First(&url)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"short_url":    "http://localhost:8080/" + url.ShortCode,
		"original_url": url.OriginalURL,
		"click_count":  url.ClickCount,
	})
}

func main() {
	// Init database
	initDB()

	r := gin.Default()
	r.Use(CORSMiddleware()) // CORS untuk FE

	r.POST("/generate-key", generateAPIKeyHandler)
	r.POST("/shorten", shortenURL)
	r.GET("/:short", redirectURL)
	r.GET("/stats/:short", getStats)

	r.Run(":8080")
}
