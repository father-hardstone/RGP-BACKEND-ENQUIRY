package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"syedibrahimshah067/RGP-BACKEND-ENQUIRY/main/models"
	"syedibrahimshah067/RGP-BACKEND-ENQUIRY/main/services"
)

// EmailController handles HTTP requests for email operations
type EmailController struct {
	emailService *services.EmailService
}

// NewEmailController creates a new instance of EmailController
func NewEmailController(emailService *services.EmailService) *EmailController {
	return &EmailController{
		emailService: emailService,
	}
}

// SendTestEmail handles GET requests to send a test email
// This endpoint sends a test email to the configured Gmail address
func (c *EmailController) SendTestEmail(w http.ResponseWriter, r *http.Request) {
	// Send test email (always goes to syedibrahimshah067@gmail.com)
	response, err := c.emailService.SendTestEmail("")
	if err != nil {
		errorResp := models.CreateErrorResponse(
			http.StatusInternalServerError,
			"Failed to send test email",
			err.Error(),
		)
		models.SendJSONResponse(w, http.StatusInternalServerError, errorResp)
		return
	}

	// Create success response
	successResp := models.Response{
		StatusCode: http.StatusOK,
		Status:     "success",
		Message:    "Test email sent successfully to syedibrahimshah067@gmail.com",
		Data:       response,
		Timestamp:  response.SentAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	models.SendJSONResponse(w, http.StatusOK, successResp)
}

// SendEmail handles POST requests to send custom emails
func (c *EmailController) SendEmail(w http.ResponseWriter, r *http.Request) {
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
	var emailReq models.EmailRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&emailReq); err != nil {
		errorResp := models.CreateErrorResponse(
			http.StatusBadRequest,
			"Invalid JSON format",
			err.Error(),
		)
		models.SendJSONResponse(w, http.StatusBadRequest, errorResp)
		return
	}

	// Validate required fields
	if emailReq.To == "" || emailReq.Subject == "" || emailReq.Body == "" {
		errorResp := models.CreateErrorResponse(
			http.StatusBadRequest,
			"Missing required fields",
			"to, subject, and body are required",
		)
		models.SendJSONResponse(w, http.StatusBadRequest, errorResp)
		return
	}

	// Send email
	response, err := c.emailService.SendEmail(&emailReq)
	if err != nil {
		errorResp := models.CreateErrorResponse(
			http.StatusInternalServerError,
			"Failed to send email",
			err.Error(),
		)
		models.SendJSONResponse(w, http.StatusInternalServerError, errorResp)
		return
	}

	// Create success response
	successResp := models.Response{
		StatusCode: http.StatusOK,
		Status:     "success",
		Message:    "Email sent successfully",
		Data:       response,
		Timestamp:  response.SentAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	models.SendJSONResponse(w, http.StatusOK, successResp)
}
