package services

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"syedibrahimshah067/RGP-BACKEND-ENQUIRY/main/config"
	"syedibrahimshah067/RGP-BACKEND-ENQUIRY/main/models"
)

// EnquiryService handles business logic for enquiry creation
// Acts as an intermediary between controllers and the database layer
type EnquiryService struct {
	db *config.Database
}

// NewEnquiryService creates a new instance of EnquiryService
// db: Database connection instance
func NewEnquiryService(db *config.Database) *EnquiryService {
	return &EnquiryService{
		db: db,
	}
}

// CreateEnquiry creates a new enquiry in the database
// query: The enquiry data to be stored
// Returns the created enquiry with generated ID and timestamps
func (s *EnquiryService) CreateEnquiry(query *models.Query) (*models.Query, error) {
	// Set creation and update timestamps
	now := primitive.NewDateTimeFromTime(time.Now())
	query.CreatedAt = now
	query.UpdatedAt = now

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Insert the enquiry into the database
	result, err := s.db.Collection.InsertOne(ctx, query)
	if err != nil {
		return nil, err
	}

	// Set the generated ID
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		query.QueryID = oid
	}

	return query, nil
}

// GetAllEnquiries retrieves all enquiries with pagination and filtering
// page: Page number (1-based)
// limit: Number of enquiries per page
// enquiryType: Filter by enquiry type (optional)
// date: Filter by creation date in YYYY-MM-DD format (optional)
// Returns paginated enquiries and total count
func (s *EnquiryService) GetAllEnquiries(page, limit int64, enquiryType, date string) (*models.EnquiryListResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Build filter
	filter := bson.M{}
	if enquiryType != "" {
		filter["enquiry_type"] = enquiryType
	}

	// Add date filter if provided
	if date != "" {
		fmt.Printf("DEBUG: Processing date filter: %s\n", date)

		// Parse the date string to create start and end of day
		parsedDate, err := time.Parse("2006-01-02", date)
		if err != nil {
			fmt.Printf("DEBUG: Failed to parse date '%s': %v\n", date, err)
		} else {
			// Create start of day (00:00:00) in UTC timezone
			startOfDay := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, time.UTC)
			// Create end of day (23:59:59.999999999) in UTC timezone
			endOfDay := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 23, 59, 59, 999999999, time.UTC)

			// Convert to primitive.DateTime for MongoDB
			startDateTime := primitive.NewDateTimeFromTime(startOfDay)
			endDateTime := primitive.NewDateTimeFromTime(endOfDay)

			fmt.Printf("DEBUG: Date range - Start: %v, End: %v\n", startOfDay, endOfDay)
			fmt.Printf("DEBUG: MongoDB range - Start: %v, End: %v\n", startDateTime, endDateTime)

			filter["created_at"] = bson.M{
				"$gte": startDateTime,
				"$lte": endDateTime,
			}
		}
	}

	// Debug: Print final filter
	fmt.Printf("DEBUG: Final MongoDB filter: %+v\n", filter)

	// Calculate skip value for pagination
	skip := (page - 1) * limit

	// Set up find options
	findOptions := options.Find()
	findOptions.SetLimit(limit)
	findOptions.SetSkip(skip)
	findOptions.SetSort(bson.M{"created_at": -1}) // Sort by creation date, newest first

	// Execute find operation
	fmt.Printf("DEBUG: Executing MongoDB query with filter: %+v\n", filter)
	cursor, err := s.db.Collection.Find(ctx, filter, findOptions)
	if err != nil {
		fmt.Printf("DEBUG: MongoDB query error: %v\n", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	// Decode results
	var enquiries []models.Query
	if err = cursor.All(ctx, &enquiries); err != nil {
		return nil, err
	}

	// Get total count for pagination
	totalCount, err := s.db.Collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Calculate pagination info
	totalPages := (totalCount + limit - 1) / limit
	hasNext := page < totalPages
	hasPrev := page > 1

	// Create response
	response := &models.EnquiryListResponse{
		Enquiries: enquiries,
		Pagination: models.PaginationInfo{
			CurrentPage:  page,
			TotalPages:   totalPages,
			TotalCount:   totalCount,
			Limit:        limit,
			HasNext:      hasNext,
			HasPrevious:  hasPrev,
			NextPage:     page + 1,
			PreviousPage: page - 1,
		},
		Filters: models.EnquiryFilters{
			EnquiryType: enquiryType,
			Date:        date,
		},
	}

	return response, nil
}

// GetEnquiryByID retrieves a specific enquiry by ID
// id: The ObjectID of the enquiry to retrieve
// Returns the enquiry if found, nil otherwise
func (s *EnquiryService) GetEnquiryByID(id primitive.ObjectID) (*models.Query, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var query models.Query
	err := s.db.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&query)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Enquiry not found
		}
		return nil, err
	}

	return &query, nil
}
