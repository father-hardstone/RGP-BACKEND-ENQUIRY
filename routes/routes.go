package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"syedibrahimshah067/RGP-BACKEND-ENQUIRY/main/controllers"
	"syedibrahimshah067/RGP-BACKEND-ENQUIRY/main/middleware"
	"syedibrahimshah067/RGP-BACKEND-ENQUIRY/main/services"
)

// SetupRoutes configures all the application routes
// Creates a new router instance and registers all endpoints
func SetupRoutes(
	rootController *controllers.RootController,
	enquiryController *controllers.EnquiryController,
	userController *controllers.UserController,
	emailController *controllers.EmailController,
	jwtService *services.JWTService,
) *mux.Router {
	// Create a new router instance
	router := mux.NewRouter()

	// Handle OPTIONS requests globally for CORS preflight
	router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Accept, Origin")
		w.Header().Set("Access-Control-Max-Age", "86400")
		w.WriteHeader(http.StatusOK)
	})

	// Root endpoint - API information
	router.HandleFunc("/", rootController.GetRoot).Methods("GET")

	// Enquiry endpoints
	router.HandleFunc("/enquiry", enquiryController.CreateEnquiry).Methods("POST")

	// User authentication endpoints (NO authentication required)
	router.HandleFunc("/create-user", userController.CreateUser).Methods("POST")
	router.HandleFunc("/auth/signin", userController.SignIn).Methods("POST")
	router.HandleFunc("/auth/login", userController.AuthenticateUser).Methods("POST")

	// Email endpoints (NO authentication required for testing)
	router.HandleFunc("/email/test", emailController.SendTestEmail).Methods("GET")
	router.HandleFunc("/email/send", emailController.SendEmail).Methods("POST")

	// Protected enquiry endpoints (require authentication)
	router.HandleFunc("/enquiries", enquiryController.GetAllEnquiries).Methods("GET")
	router.HandleFunc("/enquiries/{id}", enquiryController.GetEnquiryByID).Methods("GET")

	// Protected user management endpoints (require authentication)
	router.HandleFunc("/users", userController.GetAllUsers).Methods("GET")
	router.HandleFunc("/users/{id}", userController.GetUser).Methods("GET")
	router.HandleFunc("/users/{id}", userController.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", userController.DeleteUser).Methods("DELETE")

	// Apply authentication middleware to protected routes
	// This is applied after all routes are registered
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip authentication for OPTIONS requests (CORS preflight)
			if r.Method == "OPTIONS" {
				next.ServeHTTP(w, r)
				return
			}

			// Skip authentication for public routes
			if r.URL.Path == "/" ||
				r.URL.Path == "/enquiry" ||
				r.URL.Path == "/create-user" ||
				r.URL.Path == "/auth/signin" ||
				r.URL.Path == "/auth/login" ||
				r.URL.Path == "/health" {
				next.ServeHTTP(w, r)
				return
			}

			// Apply authentication for protected routes
			// Create the middleware chain: auth first, then role check
			// The order matters: innermost middleware runs first
			roleHandler := middleware.AdminOrSuperAdminMiddleware()(next)
			authHandler := middleware.AuthMiddleware(jwtService)(roleHandler)
			authHandler.ServeHTTP(w, r)
		})
	})

	// Note: User endpoints are now handled above in their respective sections

	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","message":"Service is running"}`))
	}).Methods("GET")

	// 404 handler for undefined routes
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"status_code":404,"status":"error","message":"Endpoint not found","error":"The requested endpoint does not exist"}`))
	})

	return router
}
