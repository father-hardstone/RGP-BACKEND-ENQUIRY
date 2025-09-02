package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"syedibrahimshah067/RGP-BACKEND-ENQUIRY/main/models"
)

// JWTService handles JWT token operations
type JWTService struct {
	secretKey []byte
}

// Claims represents the JWT claims structure
type Claims struct {
	UserID   string          `json:"user_id"`
	Email    string          `json:"email"`
	Username string          `json:"username"`
	Role     models.UserRole `json:"role"`
	jwt.RegisteredClaims
}

// NewJWTService creates a new JWT service instance
func NewJWTService(secretKey string) *JWTService {
	return &JWTService{
		secretKey: []byte(secretKey),
	}
}

// GenerateToken creates a JWT token for a user
// Different expiration times based on user role
func (s *JWTService) GenerateToken(user *models.User) (string, error) {
	// Set expiration based on user role
	var expiration time.Time
	switch user.Role {
	case models.RoleAdmin:
		// Admin: 2 days expiration
		expiration = time.Now().Add(48 * time.Hour)
	case models.RoleSuperAdmin:
		// Super-Admin: 30 days expiration (long but not infinite)
		expiration = time.Now().Add(30 * 24 * time.Hour)
	default:
		return "", errors.New("invalid user role")
	}

	// Create claims
	claims := &Claims{
		UserID:   user.ID.Hex(),
		Email:    user.Email,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "rgp-backend-enquiry",
			Subject:   user.ID.Hex(),
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Check if token is valid
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Extract claims
	if claims, ok := token.Claims.(*Claims); ok {
		return claims, nil
	}

	return nil, errors.New("invalid token claims")
}

// RefreshToken creates a new token with extended expiration
func (s *JWTService) RefreshToken(tokenString string) (string, error) {
	// Validate existing token
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	// Create new token with extended expiration
	var expiration time.Time
	switch claims.Role {
	case models.RoleAdmin:
		expiration = time.Now().Add(48 * time.Hour)
	case models.RoleSuperAdmin:
		expiration = time.Now().Add(30 * 24 * time.Hour)
	default:
		return "", errors.New("invalid user role")
	}

	// Update expiration
	claims.ExpiresAt = jwt.NewNumericDate(expiration)
	claims.IssuedAt = jwt.NewNumericDate(time.Now())
	claims.NotBefore = jwt.NewNumericDate(time.Now())

	// Create new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	newTokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", err
	}

	return newTokenString, nil
}

// GetTokenExpiration returns the expiration time for a given role
func (s *JWTService) GetTokenExpiration(role models.UserRole) time.Duration {
	switch role {
	case models.RoleAdmin:
		return 48 * time.Hour
	case models.RoleSuperAdmin:
		return 30 * 24 * time.Hour
	default:
		return 24 * time.Hour // Default fallback
	}
}

