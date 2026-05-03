package userSys

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserService defines the interface for user service operations.
type UserService interface {
	RegisterUser(username, password string) (*User, error)
	LoginUser(username, password string) (*User, error)
	GetUserDetails(username string) (*User, *UserProfile, error)
	ChangePassword(userID uint, oldPassword, newPassword string) error
	UpdateUserProfile(userID uint, updates map[string]interface{}) error
	DeactivateUser(userID uint) error
}

// userService provides methods for user-related operations.
type userService struct {
	db *gorm.DB
}

// NewUserService creates a new UserServiceInterface instance.
func NewUserService(db *gorm.DB) UserService {
	return &userService{db: db}
}

// RegisterUser registers a new user with username and password.
func (s *userService) RegisterUser(username, password string) (*User, error) {
	// Check if user already exists
	var existingUser User
	if err := s.db.Where("username = ?", username).First(&existingUser).Error; err == nil {
		return nil, errors.New("username already exists")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := User{
		Username: username,
		Password: string(hashedPassword),
	}
	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}

	// Create default profile
	profile := UserProfile{
		UserID: user.ID,
	}
	if err := s.db.Create(&profile).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// LoginUser authenticates a user and returns the user details.
func (s *userService) LoginUser(username, password string) (*User, error) {
	var user User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("username not found")
		}
		return nil, errors.New("invalid username or password")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid username or password")
	}

	return &user, nil
}

// GetUserDetails retrieves detailed user information including profile.
func (s *userService) GetUserDetails(username string) (*User, *UserProfile, error) {
	var user User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, nil, err
	}

	var profile UserProfile
	if err := s.db.Where("user_id = ?", user.ID).First(&profile).Error; err != nil {
		return nil, nil, err
	}

	return &user, &profile, nil
}

// ChangePassword updates the user's password after verifying the old password.
func (s *userService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	var user User
	if err := s.db.First(&user, userID).Error; err != nil {
		return err
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.New("invalid old password")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update password
	user.Password = string(hashedPassword)
	return s.db.Save(&user).Error
}

// UpdateUserProfile updates the user's profile information.
func (s *userService) UpdateUserProfile(userID uint, updates map[string]interface{}) error {
	return s.db.Model(&UserProfile{}).Where("user_id = ?", userID).Updates(updates).Error
}

// DeactivateUser soft deletes the user (logical deletion).
func (s *userService) DeactivateUser(userID uint) error {
	return s.db.Delete(&User{}, userID).Error
}
