package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/mihaamiharu/test-case-management-be/internal/models"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("user already exists")
)

// UserRepositoryInterface defines the interface for user repository operations
type UserRepositoryInterface interface {
	Create(user *models.User) error
	GetByEmail(email string) (*models.User, error)
	GetByID(id int64) (*models.User, error)
}

// UserRepository handles database operations for users
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create adds a new user to the database
func (r *UserRepository) Create(user *models.User) error {
	// Check if user with email already exists
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", user.Email).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrUserExists
	}

	// Insert new user
	query := `
		INSERT INTO users (username, email, password_hash, role, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	now := time.Now()
	result, err := r.db.Exec(query, user.Username, user.Email, user.PasswordHash, user.Role, now, now)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = id
	user.CreatedAt = now
	user.UpdatedAt = now
	return nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, role, created_at, updated_at
		FROM users
		WHERE email = ?
	`
	user := &models.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id int64) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, role, created_at, updated_at
		FROM users
		WHERE id = ?
	`
	user := &models.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}
