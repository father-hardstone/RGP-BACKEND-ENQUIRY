package app

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"syedibrahimshah067/RGP-BACKEND-ENQUIRY/main/config"
	"syedibrahimshah067/RGP-BACKEND-ENQUIRY/main/controllers"
	"syedibrahimshah067/RGP-BACKEND-ENQUIRY/main/middleware"
	"syedibrahimshah067/RGP-BACKEND-ENQUIRY/main/routes"
	"syedibrahimshah067/RGP-BACKEND-ENQUIRY/main/services"
	"syedibrahimshah067/RGP-BACKEND-ENQUIRY/main/utils"
)

// App represents the main application instance
// Contains all the necessary components to run the server
type App struct {
	Router *http.Server
	DB     *config.Database
	Logger *utils.Logger
}

// NewApp creates and configures a new application instance
// Initializes all components: database, services, controllers, and routes
func NewApp() (*App, error) {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load and validate database configuration
	dbConfig := config.LoadDatabaseConfig()
	if !dbConfig.ValidateConfig() {
		log.Fatal("Database configuration is incomplete")
	}

	// Connect to database
	db, err := dbConfig.Connect()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize logger
	logger := utils.NewLogger()

	// Initialize JWT service with secret key from environment
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	if jwtSecretKey == "" {
		jwtSecretKey = "default-secret-key-change-in-production"
	}
	jwtService := services.NewJWTService(jwtSecretKey)

	// Initialize services
	enquiryService := services.NewEnquiryService(db)
	userService := services.NewUserService(db, jwtService)
	emailService := services.NewEmailService()

	// Initialize controllers
	rootController := controllers.NewRootController()
	enquiryController := controllers.NewEnquiryController(enquiryService)
	userController := controllers.NewUserController(userService)
	emailController := controllers.NewEmailController(emailService)

	// Setup routes
	router := routes.SetupRoutes(rootController, enquiryController, userController, emailController, jwtService)

	// Apply middleware in correct order
	router.Use(middleware.CorsMiddleware)            // CORS first
	router.Use(middleware.LoggingMiddleware(logger)) // Logging second
	router.Use(middleware.SecurityMiddleware)        // Security third

	// Determine port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default for local dev
	}

	// Create HTTP server
	server := &http.Server{
		Addr:    "0.0.0.0:" + port,
		Handler: router,
	}

	return &App{
		Router: server,
		DB:     db,
		Logger: logger,
	}, nil
}

// Start begins the HTTP server
// Blocks until the server is stopped or encounters an error
func (a *App) Start() error {
	a.Logger.LogStartup(os.Getenv("PORT"))
	a.Logger.LogDatabaseConnection("MongoDB", true)
	return a.Router.ListenAndServe()
}

// Shutdown gracefully stops the server and closes database connections
// Should be called when the application is shutting down
func (a *App) Shutdown() {
	a.Logger.LogShutdown()
	if a.DB != nil {
		a.DB.Disconnect()
	}
}
