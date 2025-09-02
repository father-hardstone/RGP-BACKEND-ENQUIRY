package models

import "time"

// EmailRequest represents the structure for sending emails
type EmailRequest struct {
	To      string `json:"to" validate:"required,email"`
	Subject string `json:"subject" validate:"required"`
	Body    string `json:"body" validate:"required"`
}

// EmailResponse represents the response after sending an email
type EmailResponse struct {
	MessageID string    `json:"message_id"`
	To        string    `json:"to"`
	Subject   string    `json:"subject"`
	SentAt    time.Time `json:"sent_at"`
	Status    string    `json:"status"`
}

// AdminWelcomeEmail represents the structure for admin welcome emails
type AdminWelcomeEmail struct {
	To          string `json:"to" validate:"required,email"`
	FirstName   string `json:"first_name" validate:"required"`
	Username    string `json:"username" validate:"required"`
	Role        string `json:"role" validate:"required"`
	CompanyName string `json:"company_name,omitempty"`
}

