package config

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DatabaseConfig holds the database configuration parameters
// Loaded from environment variables for flexibility across different environments
type DatabaseConfig struct {
	URI                 string
	DatabaseName        string
	CollectionName      string
	UsersCollectionName string
}

// Database holds the MongoDB client and database references
// Provides methods for database operations
type Database struct {
	Client          *mongo.Client
	Database        *mongo.Database
	Collection      *mongo.Collection
	UsersCollection *mongo.Collection
}

// LoadDatabaseConfig loads database configuration from environment variables
// Returns a DatabaseConfig struct with the loaded values
func LoadDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		URI:                 os.Getenv("NEW_MONGO_URI"),
		DatabaseName:        os.Getenv("NEW_DB_NAME"),
		CollectionName:      os.Getenv("NEW_COLLECTION_NAME"),
		UsersCollectionName: os.Getenv("USERS_COLLECTION_NAME"),
	}
}

// ValidateConfig checks if all required database configuration is present
// Returns true if valid, false otherwise
func (c *DatabaseConfig) ValidateConfig() bool {
	return c.URI != "" && c.DatabaseName != "" && c.CollectionName != ""
}

// Connect establishes a connection to MongoDB using the provided configuration
// Returns a Database struct with client, database, and collection references
func (c *DatabaseConfig) Connect() (*Database, error) {
	// Create MongoDB client
	client, err := mongo.NewClient(options.Client().ApplyURI(c.URI))
	if err != nil {
		return nil, err
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	// Ping the database to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Get database and collection references
	database := client.Database(c.DatabaseName)
	collection := database.Collection(c.CollectionName)
	usersCollection := database.Collection(c.UsersCollectionName)

	log.Printf("Successfully connected to MongoDB database: %s", c.DatabaseName)

	return &Database{
		Client:          client,
		Database:        database,
		Collection:      collection,
		UsersCollection: usersCollection,
	}, nil
}

// Disconnect closes the MongoDB client connection
// Should be called when the application shuts down
func (db *Database) Disconnect() {
	if db.Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		db.Client.Disconnect(ctx)
		log.Println("MongoDB connection closed")
	}
}
