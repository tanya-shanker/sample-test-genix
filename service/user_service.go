package service

import (
	"errors"
	"fmt"

	"github.com/sample-user-service/models"
	"github.com/sample-user-service/repository"
)

var (
	ErrInvalidInput = errors.New("invalid input")
)

// UserService handles business logic for users
type UserService struct {
	repo *repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

// CreateUser creates a new user with validation
func (s *UserService) CreateUser(req *models.CreateUserRequest) (*models.User, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}

	// Check if username already exists
	if _, err := s.repo.GetByUsername(req.Username); err == nil {
		return nil, fmt.Errorf("username '%s' already exists", req.Username)
	}

	user := &models.User{
		ID:       generateID(),
		Username: req.Username,
		Email:    req.Email,
		FullName: req.FullName,
		IsActive: true,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(id string) (*models.User, error) {
	if id == "" {
		return nil, fmt.Errorf("%w: user ID is required", ErrInvalidInput)
	}

	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetAllUsers retrieves all users
func (s *UserService) GetAllUsers() ([]*models.User, error) {
	return s.repo.GetAll()
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(id string, req *models.UpdateUserRequest) (*models.User, error) {
	if id == "" {
		return nil, fmt.Errorf("%w: user ID is required", ErrInvalidInput)
	}

	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}

	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := s.repo.Update(id, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

// DeleteUser deletes a user by ID
func (s *UserService) DeleteUser(id string) error {
	if id == "" {
		return fmt.Errorf("%w: user ID is required", ErrInvalidInput)
	}

	if err := s.repo.Delete(id); err != nil {
		return err
	}

	return nil
}

// GetUserStats returns statistics about users
func (s *UserService) GetUserStats() (map[string]interface{}, error) {
	users, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	activeCount := 0
	for _, user := range users {
		if user.IsActive {
			activeCount++
		}
	}

	stats := map[string]interface{}{
		"total_users":    len(users),
		"active_users":   activeCount,
		"inactive_users": len(users) - activeCount,
	}

	return stats, nil
}

// SearchUsersByName searches users by name (partial match)
func (s *UserService) SearchUsersByName(query string) ([]*models.User, error) {
	if query == "" {
		return nil, fmt.Errorf("%w: search query is required", ErrInvalidInput)
	}

	users, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var results []*models.User
	for _, user := range users {
		if contains(user.FullName, query) || contains(user.Username, query) {
			results = append(results, user)
		}
	}

	return results, nil
}

// contains checks if a string contains a substring (case-insensitive)
func contains(str, substr string) bool {
	return len(str) >= len(substr) &&
		(str == substr || len(substr) == 0 ||
			findSubstring(str, substr))
}

// findSubstring performs case-insensitive substring search
func findSubstring(str, substr string) bool {
	for i := 0; i <= len(str)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if toLower(str[i+j]) != toLower(substr[j]) {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

// toLower converts a byte to lowercase
func toLower(b byte) byte {
	if b >= 'A' && b <= 'Z' {
		return b + 32
	}
	return b
}

// generateID generates a simple ID (in production, use UUID)
func generateID() string {
	return fmt.Sprintf("user_%d", generateRandomNumber())
}

// generateRandomNumber generates a random number (simplified for example)
func generateRandomNumber() int {
	// In production, use crypto/rand or a proper ID generator
	return 1500000
}

// Made with Bob
