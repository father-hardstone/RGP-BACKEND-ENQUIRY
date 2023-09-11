# Enquiry Submission API

This Go application provides an API for submitting and storing enquiries in a MongoDB database.

## Prerequisites

Before running the application, ensure you have the following installed:

- Go (Golang)
- MongoDB
- Required Go packages (Gin-Gonic, MongoDB driver)

## Installation

1. Clone this repository:
   `git clone https://github.com/father-hardstone/RGP-BACKEND-ENQUIRY.git`
2. Install Go dependencies:
   `go mod tidy`

## Configuration
   Configure the MongoDB URI in the "main.go" file ("mongoURI" constant).
   i.e. Make sure to replace <"mongodb+srv://XYZ"> on line 30 with your actual mongodb connection URI

## Usage
1. Start the application:
`go run main.go`
3. Make a POST request to `http://localhost:8080/enquiry` with JSON data (see below for JSON format).

4. The API will save the enquiry to the MongoDB database.

## JSON Request Format
   Sample JSON for submitting an enquiry:
```
   {
      "first_name": "John",
      "last_name": "Doe",
      "email": "johndoe@example.com",
      "phone_number": "123-456-7890",
      "company_name": "ABC Inc.",
      "enquiry_type": "General Inquiry",
      "message": "This is a sample message with a 2000 character limit."
   }
```
## Copyright and license:

 `  © 2023 Binary-Phantom Pk `

## Contributing
   Feel free to contribute to this project by opening issues or pull requests.
   
