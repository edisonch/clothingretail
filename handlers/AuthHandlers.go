package handlers

import (
	"clothingretail/db"
	"clothingretail/models"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// LoginRequest represents the login request payload
type LoginRequest struct {
	Username string `json:"username" binding:"required,max=32"`
	Pin      string `json:"pin" binding:"required,len=6"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	UserID  int    `json:"user_id,omitempty"`
}

// Login handles user authentication
func Login(c *gin.Context) {
	var req LoginRequest

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	// Validate PIN is numeric
	pin, err := strconv.Atoi(req.Pin)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Failed to authenticate",
		})
		return
	}

	// Query user from database
	var user models.ClothingUser
	query := `SELECT id, username, pin, user_status, created_at, updated_at 
	          FROM clothing_users 
	          WHERE username = ? AND pin = ? AND user_status = 1`

	err = db.DB.QueryRow(query, req.Username, pin).Scan(
		&user.ID,
		&user.Username,
		&user.Pin,
		&user.UserStatus,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			// User not found or credentials don't match
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Failed to authenticate",
			})
			return
		}
		// Database error
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to authenticate",
		})
		return
	}

	// Check if user is active
	if user.UserStatus != 1 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Failed to authenticate",
		})
		return
	}

	// Authentication successful
	// Set session cookie (simple implementation)
	c.SetCookie(
		"user_session",
		strconv.Itoa(user.ID),
		3600,  // Max age in seconds (1 hour)
		"/",   // Path
		"",    // Domain
		false, // Secure
		true,  // HttpOnly
	)

	c.JSON(http.StatusOK, LoginResponse{
		Success: true,
		Message: "Login successful",
		UserID:  user.ID,
	})
}

// Logout handles user logout
func Logout(c *gin.Context) {
	// Clear the session cookie
	c.SetCookie(
		"user_session",
		"",
		-1, // Max age -1 to delete
		"/",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logout successful",
	})
}

// AuthMiddleware checks if user is authenticated
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("user_session")
		if err != nil || cookie == "" {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		// Validate session (check if user exists and is active)
		userID, err := strconv.Atoi(cookie)
		if err != nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		var userStatus int
		err = db.DB.QueryRow("SELECT user_status FROM clothing_users WHERE id = ?", userID).Scan(&userStatus)
		if err != nil || userStatus != 1 {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		// Set user ID in context for use in handlers
		c.Set("user_id", userID)
		c.Next()
	}
}

// CreateDefaultUser creates a default admin user if no users exist
func CreateDefaultUser() error {
	var count int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM clothing_users").Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		// Create default admin user
		query := `INSERT INTO clothing_users (username, pin, user_status, created_at, updated_at)
		          VALUES (?, ?, ?, ?, ?)`

		now := time.Now()
		_, err = db.DB.Exec(query, "admin", 123456, 1, now, now)
		if err != nil {
			return err
		}
	}

	return nil
}
