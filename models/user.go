package models

import (
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// UserRole represents the role of a user in the system
type UserRole string

const (
	RoleAdmin      UserRole = "admin"
	RoleSuperAdmin UserRole = "super-admin"
)

// User represents an admin or super-admin user in the system
// This struct defines the data structure for storing user information
// in the MongoDB database and handling JSON requests/responses
type User struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username    string             `json:"username" bson:"username"`
	FirstName   string             `json:"first_name" bson:"first_name"`
	LastName    string             `json:"last_name" bson:"last_name"`
	Email       string             `json:"email" bson:"email"`
	Password    string             `json:"password,omitempty" bson:"password"`
	ProfilePic  *string            `json:"profile_pic" bson:"profile_pic"`
	Role        UserRole           `json:"role" bson:"role"`
	CompanyName *string            `json:"company_name" bson:"company_name"`
	IsActive    bool               `json:"is_active" bson:"is_active"`
	LastLogin   *time.Time         `json:"last_login,omitempty" bson:"last_login,omitempty"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// CreateUserRequest represents the request structure for creating a new user
// Username is automatically generated from email (part before @)
// Password is required in the request but will be hashed before storage
type CreateUserRequest struct {
	FirstName   string   `json:"first_name" validate:"required,min=2,max=50"`
	LastName    string   `json:"last_name" validate:"required,min=2,max=50"`
	Email       string   `json:"email" validate:"required,email"`
	Password    string   `json:"password" validate:"required,min=8"`
	ProfilePic  *string  `json:"profile_pic"`
	Role        UserRole `json:"role" validate:"required,oneof=admin super-admin"`
	CompanyName *string  `json:"company_name"`
}

// UpdateUserRequest represents the request structure for updating a user
// All fields are optional for updates
type UpdateUserRequest struct {
	Username    *string   `json:"username,omitempty"`
	FirstName   *string   `json:"first_name,omitempty"`
	LastName    *string   `json:"last_name,omitempty"`
	Email       *string   `json:"email,omitempty"`
	ProfilePic  *string   `json:"profile_pic,omitempty"`
	Role        *UserRole `json:"role,omitempty"`
	CompanyName *string   `json:"company_name,omitempty"`
	IsActive    *bool     `json:"is_active,omitempty"`
}

// UserResponse represents the response structure for user data
// Password is always omitted from responses
type UserResponse struct {
	ID          primitive.ObjectID `json:"id"`
	Username    string             `json:"username"`
	FirstName   string             `json:"first_name"`
	LastName    string             `json:"last_name"`
	Email       string             `json:"email"`
	ProfilePic  *string            `json:"profile_pic"`
	Role        UserRole           `json:"role"`
	CompanyName *string            `json:"company_name"`
	IsActive    bool               `json:"is_active"`
	LastLogin   *time.Time         `json:"last_login,omitempty"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

// Validate checks if the User struct has all required fields
// Returns true if valid, false otherwise
func (u *User) Validate() bool {
	return u.FirstName != "" && u.LastName != "" && u.Email != "" && u.Password != "" && u.Role != ""
}

// ValidateEmail performs basic email format validation
// Returns true if email format is valid, false otherwise
func (u *User) ValidateEmail() bool {
	return len(u.Email) > 0 && len(u.Email) < 255
}

// ValidateRole checks if the user role is valid
// Returns true if role is valid, false otherwise
func (u *User) ValidateRole() bool {
	return u.Role == RoleAdmin || u.Role == RoleSuperAdmin
}

// HashPassword hashes the user's password using bcrypt
// Should be called before storing the user in the database
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword compares the provided password with the stored hash
// Returns true if passwords match, false otherwise
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// ToResponse converts a User to UserResponse (omitting sensitive data)
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:          u.ID,
		Username:    u.Username,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		Email:       u.Email,
		ProfilePic:  u.ProfilePic,
		Role:        u.Role,
		CompanyName: u.CompanyName,
		IsActive:    u.IsActive,
		LastLogin:   u.LastLogin,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}

// GetFullName returns the user's full name
func (u *User) GetFullName() string {
	return u.FirstName + " " + u.LastName
}

// IsSuperAdmin checks if the user is a super admin
func (u *User) IsSuperAdmin() bool {
	return u.Role == RoleSuperAdmin
}

// IsAdmin checks if the user is an admin (including super admin)
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin || u.Role == RoleSuperAdmin
}

// GenerateUsernameFromEmail creates a username from the email address
// Takes the part before @ and ensures it's valid
func (u *User) GenerateUsernameFromEmail() string {
	if u.Email == "" {
		return ""
	}

	// Split email by @ and take the first part
	parts := strings.Split(u.Email, "@")
	if len(parts) > 0 {
		username := parts[0]

		// Remove any special characters that might cause issues
		// Keep only alphanumeric characters and underscores
		var cleanUsername strings.Builder
		for _, char := range username {
			if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') ||
				(char >= '0' && char <= '9') || char == '_' {
				cleanUsername.WriteRune(char)
			}
		}

		// Ensure minimum length and convert to lowercase
		result := strings.ToLower(cleanUsername.String())
		if len(result) < 3 {
			result = result + "user"
		}
		if len(result) > 30 {
			result = result[:30]
		}

		return result
	}

	return "user"
}

// SignInRequest represents the request structure for user sign-in
// Supports both email and username (email) for authentication
type SignInRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// SignInResponse represents the response structure for successful sign-in
// Returns user data without sensitive information
type SignInResponse struct {
	User      UserResponse `json:"user"`
	Message   string       `json:"message"`
	LoginTime time.Time    `json:"login_time"`
	Token     string       `json:"token"`
	ExpiresAt time.Time    `json:"expires_at"`
	Role      UserRole     `json:"role"`
}

// Custom error types for better error handling
type AuthError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Details string `json:"details"`
}

func (e *AuthError) Error() string {
	return e.Message
}

// Predefined authentication errors
var (
	ErrUserNotFound       = &AuthError{Type: "user_not_found", Message: "User not found", Details: "No user exists with this email address"}
	ErrInvalidPassword    = &AuthError{Type: "invalid_password", Message: "Invalid password", Details: "The password you entered is incorrect"}
	ErrAccountDeactivated = &AuthError{Type: "account_deactivated", Message: "Account deactivated", Details: "Your account has been deactivated. Please contact support"}
	ErrInvalidCredentials = &AuthError{Type: "invalid_credentials", Message: "Invalid credentials", Details: "Email or password is incorrect"}
)

// UserListResponse represents the response structure for user list data
// Contains only essential user information for listing purposes
type UserListResponse struct {
	ID         primitive.ObjectID `json:"id"`
	FirstName  string             `json:"first_name"`
	LastName   string             `json:"last_name"`
	Email      string             `json:"email"`
	ProfilePic *string            `json:"profile_pic"`
	Username   string             `json:"username"`
	Role       UserRole           `json:"role"`
}

// ToListResponse converts a User to UserListResponse (only essential fields)
func (u *User) ToListResponse() UserListResponse {
	return UserListResponse{
		ID:         u.ID,
		FirstName:  u.FirstName,
		LastName:   u.LastName,
		Email:      u.Email,
		ProfilePic: u.ProfilePic,
		Username:   u.Username,
		Role:       u.Role,
	}
}
