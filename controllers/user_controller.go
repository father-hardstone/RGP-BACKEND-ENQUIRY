package controllers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"syedibrahimshah067/RGP-BACKEND-ENQUIRY/main/models"
	"syedibrahimshah067/RGP-BACKEND-ENQUIRY/main/services"
)

// UserController handles HTTP requests related to users
// Responsible for request validation, calling business logic, and formatting responses
type UserController struct {
	userService *services.UserService
}

// NewUserController creates a new instance of UserController
// userService: Service layer instance for business logic
func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// CreateUser handles POST requests to create new users
// Validates the request, processes the user creation, and returns appropriate responses
func (c *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Check if request method is POST
	if r.Method != "POST" {
		errorResp := models.CreateErrorResponse(
			http.StatusMethodNotAllowed,
			"Method not allowed",
			"Only POST method is allowed for this endpoint",
		)
		models.SendJSONResponse(w, http.StatusMethodNotAllowed, errorResp)
		return
	}

	// Check if content type is JSON
	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		errorResp := models.CreateErrorResponse(
			http.StatusUnsupportedMediaType,
			"Unsupported media type",
			"Content-Type must be application/json",
		)
		models.SendJSONResponse(w, http.StatusUnsupportedMediaType, errorResp)
		return
	}

	// Parse JSON request body
	var createReq models.CreateUserRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&createReq); err != nil {
		errorResp := models.CreateErrorResponse(
			http.StatusBadRequest,
			"Invalid JSON format",
			err.Error(),
		)
		models.SendJSONResponse(w, http.StatusBadRequest, errorResp)
		return
	}

	// Validate required fields
	if createReq.FirstName == "" || createReq.LastName == "" || createReq.Email == "" || createReq.Password == "" {
		errorResp := models.CreateErrorResponse(
			http.StatusBadRequest,
			"Missing required fields",
			"first_name, last_name, email, and password are required",
		)
		models.SendJSONResponse(w, http.StatusBadRequest, errorResp)
		return
	}

	// Validate email format
	if !strings.Contains(createReq.Email, "@") || !strings.Contains(createReq.Email, ".") {
		errorResp := models.CreateErrorResponse(
			http.StatusBadRequest,
			"Invalid email format",
			"Email must contain @ and domain",
		)
		models.SendJSONResponse(w, http.StatusBadRequest, errorResp)
		return
	}

	// Validate password length
	if len(createReq.Password) < 8 {
		errorResp := models.CreateErrorResponse(
			http.StatusBadRequest,
			"Password too short",
			"Password must be at least 8 characters long",
		)
		models.SendJSONResponse(w, http.StatusBadRequest, errorResp)
		return
	}

	// Validate role
	if createReq.Role != models.RoleAdmin && createReq.Role != models.RoleSuperAdmin {
		errorResp := models.CreateErrorResponse(
			http.StatusBadRequest,
			"Invalid role",
			"Role must be either 'admin' or 'super-admin'",
		)
		models.SendJSONResponse(w, http.StatusBadRequest, errorResp)
		return
	}

	// Create user model from request
	user := &models.User{
		FirstName:   createReq.FirstName,
		LastName:    createReq.LastName,
		Email:       createReq.Email,
		Password:    createReq.Password,
		ProfilePic:  createReq.ProfilePic, // This will be nil if not provided
		Role:        createReq.Role,
		CompanyName: createReq.CompanyName, // This will be nil if not provided
		// Username will be auto-generated in the service layer
	}

	// Ensure ProfilePic and CompanyName are explicitly set (even if empty)
	// This guarantees they appear in the database
	if user.ProfilePic == nil {
		emptyString := ""
		user.ProfilePic = &emptyString
	}
	if user.CompanyName == nil {
		emptyString := ""
		user.CompanyName = &emptyString
	}

	// Create the user using the service layer
	createdUser, err := c.userService.CreateUser(user)
	if err != nil {
		// Check for duplicate email error
		if strings.Contains(err.Error(), "email already exists") {
			errorResp := models.CreateErrorResponse(
				http.StatusConflict,
				"User already exists",
				"Email address is already registered",
			)
			models.SendJSONResponse(w, http.StatusConflict, errorResp)
			return
		}

		errorResp := models.CreateErrorResponse(
			http.StatusInternalServerError,
			"Failed to create user",
			err.Error(),
		)
		models.SendJSONResponse(w, http.StatusInternalServerError, errorResp)
		return
	}

	// Create success response
	response := models.CreateSuccessResponse(
		http.StatusCreated,
		"User created successfully",
		createdUser.ToResponse(),
	)

	models.SendJSONResponse(w, http.StatusCreated, response)
}

// GetUser handles GET requests to retrieve a specific user by ID
// Extracts the ID from the URL path and returns the user data
func (c *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	// This would typically extract the ID from the URL path
	// For now, we'll return a method not allowed response
	errorResp := models.CreateErrorResponse(
		http.StatusMethodNotAllowed,
		"Method not allowed",
		"GET method not implemented for this endpoint",
	)
	models.SendJSONResponse(w, http.StatusMethodNotAllowed, errorResp)
}

// GetAllUsers handles GET requests to retrieve all users
// Returns a list of all users with only essential information
func (c *UserController) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	// Get all users from the service
	users, err := c.userService.GetAllUsersList(0) // 0 means no limit
	if err != nil {
		errorResp := models.CreateErrorResponse(
			http.StatusInternalServerError,
			"Failed to fetch users",
			"An error occurred while retrieving users from the database",
		)
		models.SendJSONResponse(w, http.StatusInternalServerError, errorResp)
		return
	}

	// Create success response
	response := models.Response{
		StatusCode: http.StatusOK,
		Status:     "success",
		Message:    "Users retrieved successfully",
		Data:       users,
		Timestamp:  time.Now().Format(time.RFC3339),
	}

	models.SendJSONResponse(w, http.StatusOK, response)
}

// UpdateUser handles PUT requests to update an existing user
// Validates the request and updates the user data
func (c *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// This would typically extract the ID from the URL path
	// For now, we'll return a method not allowed response
	errorResp := models.CreateErrorResponse(
		http.StatusMethodNotAllowed,
		"Method not allowed",
		"PUT method not implemented for this endpoint",
	)
	models.SendJSONResponse(w, http.StatusMethodNotAllowed, errorResp)
}

// DeleteUser handles DELETE requests to remove a user
// Removes the user from the database
func (c *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// This would typically extract the ID from the URL path
	// For now, we'll return a method not allowed response
	errorResp := models.CreateErrorResponse(
		http.StatusMethodNotAllowed,
		"Method not allowed",
		"DELETE method not implemented for this endpoint",
	)
	models.SendJSONResponse(w, http.StatusMethodNotAllowed, errorResp)
}

// AuthenticateUser handles POST requests for user authentication
// Validates credentials and returns user data if successful
func (c *UserController) AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	// This would handle login requests
	// For now, we'll return a method not allowed response
	errorResp := models.CreateErrorResponse(
		http.StatusMethodNotAllowed,
		"Method not allowed",
		"Authentication endpoint not implemented yet",
	)
	models.SendJSONResponse(w, http.StatusMethodNotAllowed, errorResp)
}

// SignIn handles POST requests for user sign-in
// Validates credentials and returns user data if successful
func (c *UserController) SignIn(w http.ResponseWriter, r *http.Request) {
	// Check if request method is POST
	if r.Method != "POST" {
		errorResp := models.CreateErrorResponse(
			http.StatusMethodNotAllowed,
			"Method not allowed",
			"Only POST method is allowed for this endpoint",
		)
		models.SendJSONResponse(w, http.StatusMethodNotAllowed, errorResp)
		return
	}

	// Check if content type is JSON
	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		errorResp := models.CreateErrorResponse(
			http.StatusUnsupportedMediaType,
			"Unsupported media type",
			"Content-Type must be application/json",
		)
		models.SendJSONResponse(w, http.StatusUnsupportedMediaType, errorResp)
		return
	}

	// Parse JSON request body
	var signInReq models.SignInRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&signInReq); err != nil {
		errorResp := models.CreateErrorResponse(
			http.StatusBadRequest,
			"Invalid JSON format",
			err.Error(),
		)
		models.SendJSONResponse(w, http.StatusBadRequest, errorResp)
		return
	}

	// Validate required fields
	if signInReq.Email == "" || signInReq.Password == "" {
		errorResp := models.CreateErrorResponse(
			http.StatusBadRequest,
			"Missing required fields",
			"email and password are required",
		)
		models.SendJSONResponse(w, http.StatusBadRequest, errorResp)
		return
	}

	// Validate email format
	if !strings.Contains(signInReq.Email, "@") || !strings.Contains(signInReq.Email, ".") {
		errorResp := models.CreateErrorResponse(
			http.StatusBadRequest,
			"Invalid email format",
			"Email must contain @ and domain",
		)
		models.SendJSONResponse(w, http.StatusBadRequest, errorResp)
		return
	}

	// Attempt to sign in the user
	signInResponse, err := c.userService.SignInUser(signInReq.Email, signInReq.Password)
	if err != nil {
		// Check for specific authentication error types
		if authErr, ok := err.(*models.AuthError); ok {
			switch authErr.Type {
			case "user_not_found":
				errorResp := models.CreateErrorResponse(
					http.StatusUnauthorized,
					"Email not found",
					authErr.Details,
				)
				models.SendJSONResponse(w, http.StatusUnauthorized, errorResp)
				return

			case "invalid_password":
				errorResp := models.CreateErrorResponse(
					http.StatusUnauthorized,
					"Wrong password",
					authErr.Details,
				)
				models.SendJSONResponse(w, http.StatusUnauthorized, errorResp)
				return

			case "account_deactivated":
				errorResp := models.CreateErrorResponse(
					http.StatusForbidden,
					"Account deactivated",
					authErr.Details,
				)
				models.SendJSONResponse(w, http.StatusForbidden, errorResp)
				return
			}
		}

		// Fallback for unexpected errors
		errorResp := models.CreateErrorResponse(
			http.StatusInternalServerError,
			"Sign-in failed",
			"An unexpected error occurred. Please try again later.",
		)
		models.SendJSONResponse(w, http.StatusInternalServerError, errorResp)
		return
	}

	// Create success response
	response := models.CreateSuccessResponse(
		http.StatusOK,
		"Sign-in successful",
		signInResponse,
	)

	models.SendJSONResponse(w, http.StatusOK, response)
}
