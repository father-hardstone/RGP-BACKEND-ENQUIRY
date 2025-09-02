# RGP Backend Enquiry API

A modern, modular Go backend API for managing enquiries/queries with MongoDB integration.

## 🎯 Business Logic

### Enquiry Management System
The application serves as a comprehensive enquiry management platform where:

- **Clients can submit enquiries** through a simple form interface
- **Enquiries are stored permanently** in MongoDB with full audit trails
- **No modification or deletion** of enquiries is allowed (business requirement)
- **Multiple enquiries** can be submitted by the same client
- **Enquiry types** can be categorized for better organization

### User Management System
The platform includes a robust user management system for administrative purposes:

- **Admin and Super-Admin roles** for access control
- **Secure authentication** with bcrypt password hashing
- **Username auto-generation** from email addresses
- **Role-based access control** for future secure routes
- **User profile management** with company affiliations

### Data Validation & Security
- **Comprehensive input validation** for all user inputs
- **Email format verification** and uniqueness checks
- **Password strength requirements** (minimum 8 characters)
- **CORS protection** for cross-origin requests
- **Security headers** to prevent common web vulnerabilities

## 🏗️ Technical Implementation

### Architecture Pattern
The application follows a **clean, layered architecture** with clear separation of concerns:

- **Models Layer**: Data structures with validation logic
- **Services Layer**: Business logic and database operations
- **Controllers Layer**: HTTP request handling and response formatting
- **Middleware Layer**: Cross-cutting concerns (CORS, logging, security)
- **Configuration Layer**: Environment-based settings management

### Database Design
- **MongoDB collections**: Separate collections for enquiries and users
- **BSON document storage** with proper indexing on frequently queried fields
- **Timestamp tracking** for creation and modification dates
- **ObjectID primary keys** for unique identification

### API Response Standardization
- **Consistent JSON structure** across all endpoints
- **Standardized error handling** with appropriate HTTP status codes
- **Pagination metadata** for list responses
- **Filtering capabilities** for data exploration

### Performance Features
- **Database connection pooling** for optimal performance
- **Context timeouts** to prevent hanging operations
- **Efficient pagination** with skip/limit optimization
- **Sorted results** by creation date (newest first)

## 🚀 Setup & Installation

### Prerequisites
- **Go 1.18 or higher**
- **MongoDB instance** (local or cloud)
- **Git** for repository cloning

### Environment Configuration

1. **Create `config.env` file** in the root directory:
   ```env
   # MongoDB Connection
   NEW_MONGO_URI=mongodb+srv://username:password@cluster.mongodb.net/
   NEW_DB_NAME=your_database_name
   NEW_COLLECTION_NAME=Enquiries
   USERS_COLLECTION_NAME=Users
   ```

2. **Update the connection string** with your actual MongoDB credentials

### Installation Steps

1. **Clone and navigate** to the project:
   ```bash
   git clone <repository-url>
   cd RGP-BACKEND-ENQUIRY
   ```

2. **Install Go dependencies**:
   ```bash
   go mod tidy
   ```

3. **Run the application**:
   ```bash
   go run main.go
   ```

4. **Access the API** at `http://localhost:8080`

## 🔧 Configuration Details

### MongoDB Setup
- **Connection URI**: Supports both local and cloud MongoDB instances
- **Database**: Separate database for the application
- **Collections**: 
  - `Enquiries`: Stores all client enquiries
  - `Users`: Stores admin and super-admin accounts

### Server Configuration
- **Port**: Default 8080 (configurable in app.go)
- **Host**: Binds to all interfaces (0.0.0.0)
- **Timeout**: 10-second database operation timeout
- **Graceful shutdown**: Handles OS signals for clean termination

### Security Settings
- **CORS**: Enabled for all origins (*)
- **Headers**: Security headers for XSS protection
- **Validation**: Input sanitization and format checking
- **Authentication**: bcrypt password hashing with default cost

## 📊 API Usage Examples

### Create Enquiry
```bash
POST http://localhost:8080/enquiry
Content-Type: application/json

{
  "first_name": "John",
  "last_name": "Doe",
  "email": "john@example.com",
  "phone_number": "+1234567890",
  "company_name": "Tech Corp",
  "enquiry_type": "general",
  "message": "Inquiry about services"
}
```

### Get Enquiries with Pagination
```bash
GET http://localhost:8080/enquiries?page=1&limit=10&enquiry_type=general
```

### Create Admin User
```bash
POST http://localhost:8080/create-user
Content-Type: application/json

{
  "first_name": "Admin",
  "last_name": "User",
  "email": "admin@company.com",
  "password": "securepassword123",
  "role": "admin"
}
```

## 🔍 Troubleshooting

### Common Issues
- **Database connection failed**: Check MongoDB URI and credentials
- **Port already in use**: Change port in app.go or kill existing process
- **Missing dependencies**: Run `go mod tidy` to install packages
- **Environment file not found**: Ensure `config.env` exists in root directory

### Debug Mode
- **Console logging**: All requests are logged with detailed information
- **Database operations**: Service layer includes debug logging for user creation
- **Error responses**: Detailed error messages for debugging

---

**Copyright © 2024 Father-Hardstone. All rights reserved.**
   
