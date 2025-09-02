package middleware

import (
	"context"
	"net/http"
	"strings"

	"syedibrahimshah067/RGP-BACKEND-ENQUIRY/main/models"
	"syedibrahimshah067/RGP-BACKEND-ENQUIRY/main/services"
)

// AuthMiddleware validates JWT tokens and protects routes
func AuthMiddleware(jwtService *services.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				models.SendJSONResponse(w, http.StatusUnauthorized, models.CreateErrorResponse(
					http.StatusUnauthorized,
					"Authorization header missing",
					"Bearer token is required",
				))
				return
			}

			// Check if it's a Bearer token
			if !strings.HasPrefix(authHeader, "Bearer ") {
				models.SendJSONResponse(w, http.StatusUnauthorized, models.CreateErrorResponse(
					http.StatusUnauthorized,
					"Invalid authorization format",
					"Authorization header must be 'Bearer <token>'",
				))
				return
			}

			// Extract token
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			// Validate token
			claims, err := jwtService.ValidateToken(tokenString)
			if err != nil {
				models.SendJSONResponse(w, http.StatusUnauthorized, models.CreateErrorResponse(
					http.StatusUnauthorized,
					"Invalid token",
					"Token is expired or invalid",
				))
				return
			}

			// Add user info to request context
			ctx := r.Context()
			ctx = context.WithValue(ctx, "user_id", claims.UserID)
			ctx = context.WithValue(ctx, "user_email", claims.Email)
			ctx = context.WithValue(ctx, "user_role", claims.Role)
			ctx = context.WithValue(ctx, "user_username", claims.Username)

			// Update request with new context
			r = r.WithContext(ctx)

			// Continue to next handler
			next.ServeHTTP(w, r)
		})
	}
}

// RoleMiddleware checks if user has required role
func RoleMiddleware(requiredRole models.UserRole) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get user role from context
			userRole, ok := r.Context().Value("user_role").(models.UserRole)
			if !ok {
				models.SendJSONResponse(w, http.StatusInternalServerError, models.CreateErrorResponse(
					http.StatusInternalServerError,
					"User role not found in context",
					"Authentication middleware must be applied before role middleware",
				))
				return
			}

			// Check if user has required role
			if userRole != requiredRole && userRole != models.RoleSuperAdmin {
				models.SendJSONResponse(w, http.StatusForbidden, models.CreateErrorResponse(
					http.StatusForbidden,
					"Insufficient permissions",
					"Access denied: insufficient role permissions",
				))
				return
			}

			// Continue to next handler
			next.ServeHTTP(w, r)
		})
	}
}

// AdminOrSuperAdminMiddleware checks if user is admin or super-admin
func AdminOrSuperAdminMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get user role from context
			userRole, ok := r.Context().Value("user_role").(models.UserRole)
			if !ok {
				models.SendJSONResponse(w, http.StatusInternalServerError, models.CreateErrorResponse(
					http.StatusInternalServerError,
					"User role not found in context",
					"Authentication middleware must be applied before role middleware",
				))
				return
			}

			// Check if user is admin or super-admin
			if userRole != models.RoleAdmin && userRole != models.RoleSuperAdmin {
				models.SendJSONResponse(w, http.StatusForbidden, models.CreateErrorResponse(
					http.StatusForbidden,
					"Insufficient permissions",
					"Access denied: admin or super-admin role required",
				))
				return
			}

			// Continue to next handler
			next.ServeHTTP(w, r)
		})
	}
}

