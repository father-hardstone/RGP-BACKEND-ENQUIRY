package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"syedibrahimshah067/RGP-BACKEND-ENQUIRY/main/config"
	"syedibrahimshah067/RGP-BACKEND-ENQUIRY/main/models"
)

// UserService handles all business logic related to users
// Acts as an intermediary between controllers and the database layer
type UserService struct {
	db         *config.Database
	jwtService *JWTService
}

// NewUserService creates a new instance of UserService
// db: Database connection instance
func NewUserService(db *config.Database, jwtService *JWTService) *UserService {
	return &UserService{
		db:         db,
		jwtService: jwtService,
	}
}

// CreateUser creates a new user in the database
// user: The user data to be stored
// Returns the created user with generated ID and timestamps
func (s *UserService) CreateUser(user *models.User) (*models.User, error) {
	// Check if user with email already exists
	existingUser, err := s.GetUserByEmail(user.Email)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Always generate username from email (since it's never provided in the request)
	user.Username = user.GenerateUsernameFromEmail()

	// Ensure username is not empty
	if user.Username == "" {
		return nil, errors.New("failed to generate username from email")
	}

	// ProfilePic and CompanyName are already handled in the controller

	// Debug: Log the user object before database insertion
	profilePicValue := ""
	if user.ProfilePic != nil {
		profilePicValue = *user.ProfilePic
	}
	companyNameValue := ""
	if user.CompanyName != nil {
		companyNameValue = *user.CompanyName
	}

	fmt.Printf("DEBUG: User before DB insertion - Username: '%s', ProfilePic: '%s', CompanyName: '%s'\n",
		user.Username, profilePicValue, companyNameValue)

	// Debug: Log the exact values and their types
	fmt.Printf("DEBUG: Username type: %T, value: '%v', length: %d\n", user.Username, user.Username, len(user.Username))
	fmt.Printf("DEBUG: ProfilePic type: %T, value: '%v', length: %d\n", user.ProfilePic, profilePicValue, len(profilePicValue))
	fmt.Printf("DEBUG: CompanyName type: %T, value: '%v', length: %d\n", user.CompanyName, companyNameValue, len(companyNameValue))

	// Handle username conflicts by adding a number suffix
	originalUsername := user.Username
	counter := 1
	for {
		existingUserByUsername, err := s.GetUserByUsername(user.Username)
		if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
			return nil, err
		}
		if existingUserByUsername == nil {
			break // Username is available
		}

		// Username exists, try with number suffix
		user.Username = fmt.Sprintf("%s%d", originalUsername, counter)
		counter++

		// Prevent infinite loop (max 999 attempts)
		if counter > 999 {
			return nil, errors.New("unable to generate unique username after multiple attempts")
		}
	}

	// Set creation and update timestamps
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	user.IsActive = true

	// Hash the password before storing
	if err := user.HashPassword(); err != nil {
		return nil, err
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Insert the user into the database
	result, err := s.db.UsersCollection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	// Set the generated ID
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		user.ID = oid
	}

	return user, nil
}

// GetUserByID retrieves a user by their ID
// id: The ObjectID of the user to retrieve
// Returns the user if found, nil otherwise
func (s *UserService) GetUserByID(id primitive.ObjectID) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := s.db.UsersCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // User not found
		}
		return nil, err
	}

	return &user, nil
}

// GetUserByEmail retrieves a user by their email address
// email: The email address of the user to retrieve
// Returns the user if found, nil otherwise
func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := s.db.UsersCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // User not found
		}
		return nil, err
	}

	return &user, nil
}

// GetUserByUsername retrieves a user by their username
// username: The username of the user to retrieve
// Returns the user if found, nil otherwise
func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := s.db.UsersCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // User not found
		}
		return nil, err
	}

	return &user, nil
}

// GetAllUsers retrieves all users from the database
// limit: Maximum number of users to return (0 for no limit)
// Returns a slice of users
func (s *UserService) GetAllUsers(limit int64) ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var users []models.User

	// Set up find options
	findOptions := options.Find()
	if limit > 0 {
		findOptions.SetLimit(limit)
	}
	findOptions.SetSort(bson.M{"created_at": -1}) // Sort by creation date, newest first

	// Execute find operation
	cursor, err := s.db.UsersCollection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Decode results
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

// GetAllUsersList retrieves all users from the database and returns simplified user data
// limit: Maximum number of users to return (0 for no limit)
// Returns a slice of UserListResponse with only essential fields
func (s *UserService) GetAllUsersList(limit int64) ([]models.UserListResponse, error) {
	users, err := s.GetAllUsers(limit)
	if err != nil {
		return nil, err
	}

	// Convert to simplified response
	var userList []models.UserListResponse
	for _, user := range users {
		userList = append(userList, user.ToListResponse())
	}

	return userList, nil
}

// UpdateUser updates an existing user in the database
// id: The ObjectID of the user to update
// updates: Map of fields to update
// Returns the updated user
func (s *UserService) UpdateUser(id primitive.ObjectID, updates map[string]interface{}) (*models.User, error) {
	// Add update timestamp
	updates["updated_at"] = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Update the user
	_, err := s.db.UsersCollection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": updates},
	)
	if err != nil {
		return nil, err
	}

	// Return the updated user
	return s.GetUserByID(id)
}

// DeleteUser removes a user from the database
// id: The ObjectID of the user to delete
// Returns true if deleted, false otherwise
func (s *UserService) DeleteUser(id primitive.ObjectID) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := s.db.UsersCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}

	return result.DeletedCount > 0, nil
}

// AuthenticateUser authenticates a user with email and password
// email: User's email address
// password: User's password (plain text)
// Returns the user if authentication successful, nil otherwise
func (s *UserService) AuthenticateUser(email, password string) (*models.User, error) {
	user, err := s.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, models.ErrUserNotFound
	}

	if !user.IsActive {
		return nil, models.ErrAccountDeactivated
	}

	if !user.CheckPassword(password) {
		return nil, models.ErrInvalidPassword
	}

	// Update last login time
	now := time.Now()
	_, err = s.UpdateUser(user.ID, map[string]interface{}{
		"last_login": now,
	})
	if err != nil {
		// Log the error but don't fail authentication
		// You might want to add proper logging here
	}

	return user, nil
}

// SignInUser handles the complete sign-in process
// email: User's email address
// password: User's password (plain text)
// Returns SignInResponse with user data, login information, and JWT token
func (s *UserService) SignInUser(email, password string) (*models.SignInResponse, error) {
	// Authenticate the user
	user, err := s.AuthenticateUser(email, password)
	if err != nil {
		return nil, err
	}

	// Generate JWT token
	token, err := s.jwtService.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	// Get token expiration
	expiration := s.jwtService.GetTokenExpiration(user.Role)

	// Create sign-in response
	response := &models.SignInResponse{
		User:      user.ToResponse(),
		Message:   "Sign-in successful",
		LoginTime: time.Now(),
		Token:     token,
		ExpiresAt: time.Now().Add(expiration),
		Role:      user.Role,
	}

	return response, nil
}

// GetUsersByRole retrieves all users with a specific role
// role: The role to filter by
// Returns a slice of users with the specified role
func (s *UserService) GetUsersByRole(role models.UserRole) ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var users []models.User

	// Execute find operation with role filter
	cursor, err := s.db.UsersCollection.Find(ctx, bson.M{"role": role})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Decode results
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

// CountUsers returns the total number of users in the database
// Returns the count of users
func (s *UserService) CountUsers() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err := s.db.UsersCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, err
	}

	return count, nil
}
