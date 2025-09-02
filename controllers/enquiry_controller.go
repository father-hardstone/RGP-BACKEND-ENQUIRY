package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"syedibrahimshah067/RGP-BACKEND-ENQUIRY/main/models"
	"syedibrahimshah067/RGP-BACKEND-ENQUIRY/main/services"
)

// EnquiryController handles HTTP requests related to enquiries
// Responsible for request validation, calling business logic, and formatting responses
type EnquiryController struct {
	enquiryService *services.EnquiryService
}

// NewEnquiryController creates a new instance of EnquiryController
// enquiryService: Service layer instance for business logic
func NewEnquiryController(enquiryService *services.EnquiryService) *EnquiryController {
	return &EnquiryController{
		enquiryService: enquiryService,
	}
}

// CreateEnquiry handles POST requests to create new enquiries
// Validates the request, processes the enquiry, and returns appropriate responses
func (c *EnquiryController) CreateEnquiry(w http.ResponseWriter, r *http.Request) {
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

	// Check if content type is JSON (allow charset parameter)
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
	var query models.Query
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&query); err != nil {
		errorResp := models.CreateErrorResponse(
			http.StatusBadRequest,
			"Invalid JSON format",
			err.Error(),
		)
		models.SendJSONResponse(w, http.StatusBadRequest, errorResp)
		return
	}

	// Validate required fields
	if !query.Validate() {
		errorResp := models.CreateErrorResponse(
			http.StatusBadRequest,
			"Missing required fields",
			"first_name, last_name, email, and message are required",
		)
		models.SendJSONResponse(w, http.StatusBadRequest, errorResp)
		return
	}

	// Validate email format
	if !query.ValidateEmail() {
		errorResp := models.CreateErrorResponse(
			http.StatusBadRequest,
			"Invalid email format",
			"Please provide a valid email address",
		)
		models.SendJSONResponse(w, http.StatusBadRequest, errorResp)
		return
	}

	// Additional email format validation
	if !strings.Contains(query.Email, "@") || !strings.Contains(query.Email, ".") {
		errorResp := models.CreateErrorResponse(
			http.StatusBadRequest,
			"Invalid email format",
			"Email must contain @ and domain",
		)
		models.SendJSONResponse(w, http.StatusBadRequest, errorResp)
		return
	}

	// Create the enquiry using the service layer
	createdQuery, err := c.enquiryService.CreateEnquiry(&query)
	if err != nil {
		errorResp := models.CreateErrorResponse(
			http.StatusInternalServerError,
			"Failed to create enquiry",
			err.Error(),
		)
		models.SendJSONResponse(w, http.StatusInternalServerError, errorResp)
		return
	}

	// Create success response
	response := models.CreateSuccessResponse(
		http.StatusCreated,
		"Enquiry submitted successfully. We will get back to you soon.",
		map[string]interface{}{
			"enquiry_id":   createdQuery.QueryID,
			"submitted_at": createdQuery.CreatedAt,
		},
	)

	models.SendJSONResponse(w, http.StatusCreated, response)
}

// GetAllEnquiries handles GET requests to retrieve all enquiries with pagination and filtering
// Supports query parameters: page, limit, enquiry_type
func (c *EnquiryController) GetAllEnquiries(w http.ResponseWriter, r *http.Request) {
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

	// Parse query parameters
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	enquiryType := r.URL.Query().Get("enquiry_type")
	date := r.URL.Query().Get("date")

	// Set default values
	page := int64(1)
	limit := int64(10)

	// Parse page parameter
	if pageStr != "" {
		if parsedPage, err := strconv.ParseInt(pageStr, 10, 64); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	// Parse limit parameter (max 100 per page)
	if limitStr != "" {
		if parsedLimit, err := strconv.ParseInt(limitStr, 10, 64); err == nil && parsedLimit > 0 {
			if parsedLimit > 100 {
				parsedLimit = 100
			}
			limit = parsedLimit
		}
	}

	// Get enquiries from service
	enquiries, err := c.enquiryService.GetAllEnquiries(page, limit, enquiryType, date)
	if err != nil {
		errorResp := models.CreateErrorResponse(
			http.StatusInternalServerError,
			"Failed to retrieve enquiries",
			err.Error(),
		)
		models.SendJSONResponse(w, http.StatusInternalServerError, errorResp)
		return
	}

	// Create success response
	response := models.CreateSuccessResponse(
		http.StatusOK,
		"Enquiries retrieved successfully",
		enquiries,
	)

	models.SendJSONResponse(w, http.StatusOK, response)
}

// GetEnquiryByID handles GET requests to retrieve a specific enquiry by ID
func (c *EnquiryController) GetEnquiryByID(w http.ResponseWriter, r *http.Request) {
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

	// Extract ID from URL path
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Parse ObjectID
	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		errorResp := models.CreateErrorResponse(
			http.StatusBadRequest,
			"Invalid enquiry ID",
			"Please provide a valid enquiry ID",
		)
		models.SendJSONResponse(w, http.StatusBadRequest, errorResp)
		return
	}

	// Get enquiry from service
	enquiry, err := c.enquiryService.GetEnquiryByID(objectID)
	if err != nil {
		errorResp := models.CreateErrorResponse(
			http.StatusInternalServerError,
			"Failed to retrieve enquiry",
			err.Error(),
		)
		models.SendJSONResponse(w, http.StatusInternalServerError, errorResp)
		return
	}

	if enquiry == nil {
		errorResp := models.CreateErrorResponse(
			http.StatusNotFound,
			"Enquiry not found",
			"No enquiry found with the provided ID",
		)
		models.SendJSONResponse(w, http.StatusNotFound, errorResp)
		return
	}

	// Create success response
	response := models.CreateSuccessResponse(
		http.StatusOK,
		"Enquiry retrieved successfully",
		enquiry,
	)

	models.SendJSONResponse(w, http.StatusOK, response)
}
