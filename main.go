package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	"os"
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
	r.Use(corsMiddleware())
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
	log.Fatal(r.Run("0.0.0.0:8080"))
}

// EnquiryHandler handles POST requests to store the enquiry data in MongoDB.
func EnquiryHandler(c *gin.Context) {
	var q Query
	clientIP := c.ClientIP()

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

	c.JSON(http.StatusCreated, gin.H{"message": "Thanks for reaching out. we will get back to you."})
	logToFile("Client IP Address: " + clientIP)
}

//middle function...
func corsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(http.StatusOK)
        } else {
            c.Next()
        }
    }
}
func logToFile(message string) {
    // Open the file in append mode. Create it if it doesn't exist.
    file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Println("Error opening log file:", err)
        return
    }
    defer file.Close()

    // Write the message to the file
    _, err = file.WriteString(message + "\n")
    if err != nil {
        log.Println("Error writing to log file:", err)
    }
}