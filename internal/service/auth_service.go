package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mihaamiharu/test-case-management-be/internal/config"
	"github.com/mihaamiharu/test-case-management-be/internal/models"
	"github.com/mihaamiharu/test-case-management-be/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo  repository.UserRepositoryInterface
	jwtSecret string
	jwtExpiry int
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo repository.UserRepositoryInterface, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: cfg.JWTSecret,
		jwtExpiry: cfg.JWTExpiryHours,
	}
}

// Register registers a new user
func (s *AuthService) Register(userCreate *models.UserCreate) (*models.User, error) {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userCreate.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Set default role if not provided
	role := userCreate.Role
	if role == "" {
		role = models.RoleUser
	}

	user := &models.User{
		Username:     userCreate.Username,
		Email:        userCreate.Email,
		PasswordHash: string(hashedPassword),
		Role:         role,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := s.generateToken(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

// GetUserByEmail retrieves a user by email
func (s *AuthService) GetUserByEmail(email string) (*models.User, error) {
	return s.userRepo.GetByEmail(email)
}

// generateToken creates a new JWT token for a user
func (s *AuthService) generateToken(user *models.User) (string, error) {
	// Create the Claims
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * time.Duration(s.jwtExpiry)).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token
func (s *AuthService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.jwtSecret), nil
	})
}

// GetUserFromToken extracts the user from a validated token
func (s *AuthService) GetUserFromToken(token *jwt.Token) (*models.User, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errors.New("invalid user_id in token")
	}

	return s.userRepo.GetByID(int64(userID))
}
