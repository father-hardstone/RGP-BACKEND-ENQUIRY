package services

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"gopkg.in/mail.v2"

	"syedibrahimshah067/RGP-BACKEND-ENQUIRY/main/models"
)

// EmailService handles all email-related operations
type EmailService struct {
	host     string
	port     int
	username string
	password string
	fromName string
}

// NewEmailService creates a new instance of EmailService
func NewEmailService() *EmailService {
	// Default to 587 for STARTTLS (Gmail standard)
	port := 587
	if envPort := os.Getenv("EMAIL_PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			port = p
		}
	}

	// Debug logs to confirm env variables
	fmt.Println("=== Email Service Configuration ===")
	fmt.Printf("Host: %s\n", os.Getenv("EMAIL_HOST"))
	fmt.Printf("Port: %d\n", port)
	fmt.Printf("Username: %s\n", os.Getenv("EMAIL_USERNAME"))
	fmt.Printf("App Password (hidden): %v\n", len(os.Getenv("EMAIL_APP_PASSWORD")) > 0)
	fmt.Printf("From Name: %s\n", os.Getenv("EMAIL_FROM_NAME"))
	fmt.Println("===================================")

	return &EmailService{
		host:     os.Getenv("EMAIL_HOST"),
		port:     port,
		username: os.Getenv("EMAIL_USERNAME"),
		password: os.Getenv("EMAIL_APP_PASSWORD"), // Use app password instead of regular password
		fromName: os.Getenv("EMAIL_FROM_NAME"),
	}
}

// SendEmail sends a basic email
func (s *EmailService) SendEmail(req *models.EmailRequest) (*models.EmailResponse, error) {
	fmt.Println("Preparing to send email...")
	fmt.Printf("To: %s\n", req.To)
	fmt.Printf("Subject: %s\n", req.Subject)

	// Create new message
	m := mail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", s.fromName, s.username))
	m.SetHeader("To", req.To)
	m.SetHeader("Subject", req.Subject)
	m.SetBody("text/html", req.Body)

	// Create dialer with TLS enabled for Gmail
	d := mail.NewDialer(s.host, s.port, s.username, s.password)
	d.StartTLSPolicy = mail.MandatoryStartTLS // Important for Gmail (587 STARTTLS)

	fmt.Printf("Connecting to Gmail SMTP server %s:%d using STARTTLS...\n", s.host, s.port)

	// Send email
	if err := d.DialAndSend(m); err != nil {
		fmt.Printf("Error while sending email: %v\n", err)
		return nil, fmt.Errorf("failed to send email: %v", err)
	}

	fmt.Println("Email sent successfully!")

	// Create response
	response := &models.EmailResponse{
		MessageID: fmt.Sprintf("msg_%d", time.Now().Unix()),
		To:        req.To,
		Subject:   req.Subject,
		SentAt:    time.Now(),
		Status:    "sent",
	}

	return response, nil
}

// SendAdminWelcomeEmail sends a welcome email to new admin users
func (s *EmailService) SendAdminWelcomeEmail(req *models.AdminWelcomeEmail) (*models.EmailResponse, error) {
	subject := fmt.Sprintf("Welcome to RGP Backend - %s Role", req.Role)

	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>Welcome to RGP Backend, %s!</h2>
			<p>Your account has been successfully created with the following details:</p>
			<ul>
				<li><strong>Username:</strong> %s</li>
				<li><strong>Role:</strong> %s</li>
				<li><strong>Company:</strong> %s</li>
			</ul>
			<p>You can now sign in to access the admin panel.</p>
			<br>
			<p>Best regards,<br>RGP Backend Team</p>
		</body>
		</html>
	`, req.FirstName, req.Username, req.Role, req.CompanyName)

	emailReq := &models.EmailRequest{
		To:      req.To,
		Subject: subject,
		Body:    body,
	}

	return s.SendEmail(emailReq)
}

// SendTestEmail sends a test email (for testing purposes)
func (s *EmailService) SendTestEmail(to string) (*models.EmailResponse, error) {
	fixedRecipient := "khanbahaduribrahim@outlook.com"

	req := &models.EmailRequest{
		To:      fixedRecipient,
		Subject: "Test Email from RGP Backend",
		Body: `
			<html>
			<body>
				<h2>Test Email</h2>
				<p>This is a test email from your RGP Backend service.</p>
				<p>If you received this, your email service is working correctly!</p>
				<br>
				<p>Sent at: ` + time.Now().Format("2006-01-02 15:04:05") + `</p>
				<p><strong>Note:</strong> This email was sent via Gmail SMTP to your Outlook address.</p>
			</body>
			</html>
		`,
	}

	return s.SendEmail(req)
}
