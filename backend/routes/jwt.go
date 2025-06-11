package routes

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kubestellar/ui/auth"
	"github.com/kubestellar/ui/middleware"
	"github.com/kubestellar/ui/models"
	"github.com/kubestellar/ui/telemetry"
	"github.com/kubestellar/ui/utils"
)

// SetupRoutes initializes all application routes
func setupAuthRoutes(router *gin.Engine) {
	// Authentication routes
	router.POST("/login", LoginHandler)

	// API group for all endpoints
	api := router.Group("/api")

	// Protected API endpoints requiring authentication
	protected := api.Group("/")
	protected.Use(middleware.AuthenticateMiddleware())
	{
		protected.GET("/me", CurrentUserHandler)

		// Read-only endpoints
		read := protected.Group("/")
		read.Use(middleware.RequirePermission("read"))
		{
			read.GET("/resources", GetResourcesHandler)
		}

		// Write-requiring endpoints
		write := protected.Group("/auth")
		write.Use(middleware.RequirePermission("write"))
		{
			write.POST("/auth/resources", CreateResourceHandler)
			write.PUT("/auth/resources/:id", UpdateResourceHandler)
			write.DELETE("/auth/resources/:id", DeleteResourceHandler)
		}

		// Admin-only endpoints
		admin := protected.Group("/admin")
		admin.Use(middleware.RequireAdmin())
		{
			admin.GET("/users", ListUsersHandler)
			admin.POST("/users", CreateUserHandler)
			admin.PUT("/users/:username", UpdateUserHandler)
			admin.DELETE("/users/:username", DeleteUserHandler)
		}
	}

	// Setup other route groups as needed
	setupAdditionalRoutes(router)
}

// setupAdditionalRoutes adds any additional route groups
func setupAdditionalRoutes(router *gin.Engine) {
	// Add additional routes here as needed
}

// LoginHandler verifies user credentials and issues JWT
func LoginHandler(c *gin.Context) {
	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		telemetry.HTTPErrorCounter.WithLabelValues("POST", "/login", "400").Inc()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	loginData.Username = strings.TrimSpace(loginData.Username)
	loginData.Password = strings.TrimSpace(loginData.Password)

	user, err := models.AuthenticateUser(loginData.Username, loginData.Password)
	if user == nil || err != nil {
		telemetry.HTTPErrorCounter.WithLabelValues("POST", "/login", "401").Inc()
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Fixed: Pass both username and permissions to GenerateToken
	token, err := utils.GenerateToken(loginData.Username, user.Permissions)
	if err != nil {
		telemetry.HTTPErrorCounter.WithLabelValues("POST", "/login", "500").Inc()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}
	telemetry.TotalHTTPRequests.WithLabelValues("POST", "/login", "200").Inc()
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"username":    user.Username,
			"permissions": user.Permissions,
		},
	})
}

// CurrentUserHandler returns the current user's information
func CurrentUserHandler(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		telemetry.HTTPErrorCounter.WithLabelValues("GET", "/api/me", "401").Inc()
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	permissions, exists := c.Get("permissions")
	if !exists {
		telemetry.HTTPErrorCounter.WithLabelValues("GET", "/api/me", "500").Inc()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Permissions not found"})
		return
	}
	telemetry.TotalHTTPRequests.WithLabelValues("GET", "/api/me", "200").Inc()
	c.JSON(http.StatusOK, gin.H{
		"username":    username,
		"permissions": permissions,
	})
}

// ListUsersHandler returns a list of all users (admin only)
func ListUsersHandler(c *gin.Context) {
	users, err := auth.ListUsersWithPermissions()
	if err != nil {
		telemetry.HTTPErrorCounter.WithLabelValues("GET", "/api/admin/users", "500").Inc()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}
	telemetry.TotalHTTPRequests.WithLabelValues("GET", "/api/admin/users", "200").Inc()
	c.JSON(http.StatusOK, gin.H{"users": users})
}

// CreateUserHandler creates a new user (admin only)
func CreateUserHandler(c *gin.Context) {
	var userData struct {
		Username    string   `json:"username" binding:"required"`
		Password    string   `json:"password" binding:"required"`
		Permissions []string `json:"permissions" binding:"required"`
	}

	if err := c.ShouldBindJSON(&userData); err != nil {
		telemetry.HTTPErrorCounter.WithLabelValues("POST", "/api/admin/users", "400").Inc()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	err := auth.AddOrUpdateUser(userData.Username, userData.Password, userData.Permissions)
	if err != nil {
		telemetry.HTTPErrorCounter.WithLabelValues("POST", "/api/admin/users", "500").Inc()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create user",
			"details": err.Error(),
		})
		return
	}
	telemetry.TotalHTTPRequests.WithLabelValues("POST", "/api/admin/users", "201").Inc()
	c.JSON(http.StatusCreated, gin.H{
		"message":  "User created successfully",
		"username": userData.Username,
	})
}

// UpdateUserHandler updates an existing user (admin only)
func UpdateUserHandler(c *gin.Context) {
	username := c.Param("username")

	var userData struct {
		Password    string   `json:"password"`
		Permissions []string `json:"permissions"`
	}

	if err := c.ShouldBindJSON(&userData); err != nil {
		telemetry.HTTPErrorCounter.WithLabelValues("PUT", "/api/admin/users/"+username, "400").Inc()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Get existing user data
	userConfig, exists, err := auth.GetUserByUsername(username)
	if err != nil {
		telemetry.HTTPErrorCounter.WithLabelValues("PUT", "/api/admin/users/"+username, "500").Inc()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}

	if !exists {
		telemetry.HTTPErrorCounter.WithLabelValues("PUT", "/api/admin/users/"+username, "404").Inc()
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Update only provided fields
	if userData.Password != "" {
		// Update password if provided
		userConfig.Password = userData.Password
	}

	// Update permissions if provided
	if userData.Permissions != nil {
		userConfig.Permissions = userData.Permissions
	}

	// Save updated user
	err = auth.AddOrUpdateUser(username, userConfig.Password, userConfig.Permissions)
	if err != nil {
		telemetry.HTTPErrorCounter.WithLabelValues("PUT", "/api/admin/users/"+username, "500").Inc()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update user",
			"details": err.Error(),
		})
		return
	}
	telemetry.TotalHTTPRequests.WithLabelValues("PUT", "/api/admin/users/"+username, "200").Inc()
	c.JSON(http.StatusOK, gin.H{
		"message":  "User updated successfully",
		"username": username,
	})
}

// DeleteUserHandler deletes a user (admin only)
func DeleteUserHandler(c *gin.Context) {
	username := c.Param("username")

	err := auth.RemoveUser(username)
	if err != nil {
		telemetry.HTTPErrorCounter.WithLabelValues("DELETE", "/api/admin/users/"+username, "500").Inc()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete user",
			"details": err.Error(),
		})
		return
	}
	telemetry.TotalHTTPRequests.WithLabelValues("DELETE", "/api/admin/users/"+username, "200").Inc()
	c.JSON(http.StatusOK, gin.H{
		"message":  "User deleted successfully",
		"username": username,
	})
}

// Example handlers for resource endpoints (implement as needed)
func GetResourcesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Resources retrieved successfully"})
}

func CreateResourceHandler(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "Resource created successfully"})
}

func UpdateResourceHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Resource updated successfully"})
}

func DeleteResourceHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Resource deleted successfully"})
}
