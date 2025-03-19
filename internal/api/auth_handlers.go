package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mihaamiharu/test-case-management-be/internal/models"
	"github.com/mihaamiharu/test-case-management-be/internal/repository"
	"github.com/mihaamiharu/test-case-management-be/internal/service"
)

// AuthHandler handles authentication-related API endpoints
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var userCreate models.UserCreate
	if err := c.ShouldBindJSON(&userCreate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.Register(&userCreate)
	if err != nil {
		if err == repository.ErrUserExists {
			c.JSON(http.StatusConflict, gin.H{"error": "User with this email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	// Generate token for the newly registered user
	token, err := h.authService.Login(userCreate.Email, userCreate.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User registered but failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"token": token,
		"user":  user.ToResponse(),
	})
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var userLogin models.UserLogin
	if err := c.ShouldBindJSON(&userLogin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authService.Login(userLogin.Email, userLogin.Password)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login"})
		return
	}

	// Get user information
	user, err := h.authService.GetUserByEmail(userLogin.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user information"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user.ToResponse(),
	})
}

// Me returns the current authenticated user
func (h *AuthHandler) Me(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	user, ok := userInterface.(*models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user from context"})
		return
	}

	c.JSON(http.StatusOK, user.ToResponse())
}

// AuthMiddleware authenticates requests
func (h *AuthHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header must start with 'Bearer '"})
			c.Abort()
			return
		}

		// Remove "Bearer " prefix
		tokenString := authHeader[7:]

		// Validate token
		token, err := h.authService.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Get user from token
		user, err := h.authService.GetUserFromToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Set user and userID in context
		c.Set("user", user)
		c.Set("userID", user.ID)
		c.Next()
	}
}
