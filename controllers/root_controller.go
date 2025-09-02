package controllers

import (
	"net/http"

	"syedibrahimshah067/RGP-BACKEND-ENQUIRY/main/models"
)

// RootController handles requests to the root endpoint
// Provides API information and available endpoints
type RootController struct{}

// NewRootController creates a new instance of RootController
func NewRootController() *RootController {
	return &RootController{}
}

// GetRoot handles GET requests to the root endpoint
// Returns API information, available endpoints, and system status
func (c *RootController) GetRoot(w http.ResponseWriter, r *http.Request) {
	// Check if request method is GET
	if r.Method != "GET" {
		errorResp := models.CreateErrorResponse(
			http.StatusMethodNotAllowed,
			"Method not allowed",
			"Only GET method is allowed for this endpoint",
		)
		models.SendJSONResponse(w, http.StatusMethodNotAllowed, errorResp)
		return
	}

	// Define available endpoints
	endpoints := map[string]string{
		"GET  /":               "API information and available endpoints",
		"POST /enquiry":        "Create a new enquiry",
		"POST /create-user":    "Create a new admin/super-admin user",
		"POST /auth/signin":    "User sign-in with JWT authentication",
		"POST /auth/login":     "User login (legacy endpoint)",
		"GET  /email/test":     "Send a test email to configured address",
		"POST /email/send":     "Send a custom email",
		"GET  /enquiries":      "Get all enquiries with pagination, filtering by enquiry_type and date (protected - requires auth)",
		"GET  /enquiries/{id}": "Get enquiry by ID (protected - requires auth)",
		"GET  /users":          "Get all users (protected - requires auth)",
		"GET  /users/{id}":     "Get user by ID (protected - requires auth)",
		"PUT  /users/{id}":     "Update user (protected - requires auth)",
		"DELETE /users/{id}":   "Delete user (protected - requires auth)",
		"GET  /health":         "Health check endpoint",
	}

	// Create success response with API information
	response := models.CreateSuccessResponse(
		http.StatusOK,
		"Welcome to RGP Backend Enquiry API",
		map[string]interface{}{
			"endpoints":   endpoints,
			"version":     "1.0.0",
			"description": "RGP Backend Enquiry Management System",
			"status":      "running",
		},
	)

	models.SendJSONResponse(w, http.StatusOK, response)
}
