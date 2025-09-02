package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Query represents an enquiry/query submitted by a user
// This struct defines the data structure for storing enquiry information
// in the MongoDB database and handling JSON requests/responses
type Query struct {
	QueryID     primitive.ObjectID `json:"queryid" bson:"_id,omitempty"`
	FirstName   string             `json:"first_name" bson:"first_name"`
	LastName    string             `json:"last_name" bson:"last_name"`
	Email       string             `json:"email" bson:"email"`
	PhoneNumber string             `json:"phone_number" bson:"phone_number"`
	CompanyName string             `json:"company_name" bson:"company_name"`
	EnquiryType string             `json:"enquiry_type" bson:"enquiry_type"`
	Message     string             `json:"message" bson:"message"`
	CreatedAt   primitive.DateTime `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt   primitive.DateTime `json:"updated_at" bson:"updated_at,omitempty"`
}

// Validate checks if the Query struct has all required fields
// Returns true if valid, false otherwise
func (q *Query) Validate() bool {
	return q.FirstName != "" && q.LastName != "" && q.Email != "" && q.Message != ""
}

// ValidateEmail performs basic email format validation
// Returns true if email format is valid, false otherwise
func (q *Query) ValidateEmail() bool {
	return len(q.Email) > 0 && len(q.Email) < 255
}

// PaginationInfo contains pagination metadata for list responses
type PaginationInfo struct {
	CurrentPage  int64 `json:"current_page"`
	TotalPages   int64 `json:"total_pages"`
	TotalCount   int64 `json:"total_count"`
	Limit        int64 `json:"limit"`
	HasNext      bool  `json:"has_next"`
	HasPrevious  bool  `json:"has_previous"`
	NextPage     int64 `json:"next_page"`
	PreviousPage int64 `json:"previous_page"`
}

// EnquiryFilters contains filter parameters for enquiry queries
type EnquiryFilters struct {
	EnquiryType string `json:"enquiry_type,omitempty"`
	Date        string `json:"date,omitempty"` // Format: YYYY-MM-DD
}

// EnquiryListResponse contains paginated enquiries with metadata
type EnquiryListResponse struct {
	Enquiries  []Query        `json:"enquiries"`
	Pagination PaginationInfo `json:"pagination"`
	Filters    EnquiryFilters `json:"filters"`
}
