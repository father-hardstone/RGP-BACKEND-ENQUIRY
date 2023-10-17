package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
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
	mongoURI       = "mongodb+srv://rgpitglobal:blS6uA4ZnB9OnGRV@rgpitglobal.up1wnpe.mongodb.net/?retryWrites=true&w=majority"
	dbName         = "RGP"
	collectionName = "Enquiries"
)

func main() {
	r := mux.NewRouter()
	r.Use(CorsMiddleware)
	// Add custom logging middleware
	r.Use(loggingMiddleware)
	r.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusOK)
	})
	// Define API routes
	r.HandleFunc("/enquiry", EnquiryHandler).Methods("POST")
	r.HandleFunc("/", RootHandler).Methods("GET")
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	fmt.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", r))
}
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
func RootHandler(w http.ResponseWriter, r *http.Request) {
	// Response message
	message := "You have reached the end of the line...\nState your wish!!!"

	// Write the response to the client
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, message)
}
func EnquiryHandler(w http.ResponseWriter, r *http.Request) {
	var q Query

	// Parse JSON request body into the Query struct
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&q); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode JSON: %s", err.Error()), http.StatusBadRequest)
		return
	}

	// Create a MongoDB client
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to connect to MongoDB: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	err = client.Connect(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to connect to MongoDB: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(ctx)

	// Get a handle to the collection
	collection := client.Database(dbName).Collection(collectionName)

	// Insert the enquiry data into MongoDB
	_, err = collection.InsertOne(ctx, q)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to insert data into MongoDB: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Send a JSON response for success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Thanks for reaching out. We will get back to you.",
	})

	// //Alternatively, send a JSON response for failure
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusInternalServerError)
	// json.NewEncoder(w).Encode(map[string]interface{}{
	// 	"status":  "error",
	// 	"message": "Internal Server Error. Please try again later.",
	// })
}

// loggingMiddleware is a custom middleware function for logging requests.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// Create a response writer that captures the status code and response size
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		// Call the next handler in the chain
		next.ServeHTTP(lrw, r)
		// Log the request details and status code
		fmt.Printf("[%s]'   '%s'   '%s'   '%s'   '%v'   '-'   '%d\n", r.Method, r.RemoteAddr, r.URL.Path, r.Proto, time.Since(start), lrw.statusCode)
	})
}

// loggingResponseWriter is a custom ResponseWriter that captures the status code.
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}