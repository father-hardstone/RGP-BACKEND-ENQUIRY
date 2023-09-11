package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Query struct to represent the data.
type Query struct {
	QueryID     primitive.ObjectID `json:"queryid" bson:"_id,omitempty"`
	FirstName   string             `json:"first_name"`
	LastName    string             `json:"last_name"`
	Email       string             `json:"email"`
	PhoneNumber string             `json:"phone_number"`
	CompanyName string             `json:"company_name"`
	EnquiryType string             `json:"enquiry_type"`
	Message     string             `json:"message"`
}

// MongoDB configuration
const (
	mongoURI       = "mongodb+srv://Faridi:tAZXwJjYgiKHo6WB@zeto.vxe0b.mongodb.net/?retryWrites=true&w=majority"
	dbName         = "RGP"
	collectionName = "Enquiries"
)

func main() {
	// Initialize Gin-Gonic router
	r := gin.Default()

	// Define API routes
	r.POST("/enquiry", EnquiryHandler)

	// Create a MongoDB client
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	// Start the HTTP server
	fmt.Println("Server is running on :8080")
	log.Fatal(r.Run(":8080"))
}

// EnquiryHandler handles POST requests to store the enquiry data in MongoDB.
func EnquiryHandler(c *gin.Context) {
	var q Query

	// Generate a unique query ID
	q.QueryID = primitive.NewObjectID()

	// Parse JSON request body into the Query struct
	if err := c.ShouldBindJSON(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create a MongoDB client
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	err = client.Connect(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer client.Disconnect(ctx)

	// Get a handle to the collection
	collection := client.Database(dbName).Collection(collectionName)

	// Insert the enquiry data into MongoDB
	_, err = collection.InsertOne(ctx, q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Enquiry successfully submitted"})
}