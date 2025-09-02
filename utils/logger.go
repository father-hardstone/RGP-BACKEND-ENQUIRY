package utils

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// LogEntry represents a single log entry with all request details
type LogEntry struct {
	Timestamp  time.Time
	Method     string
	RemoteAddr string
	Path       string
	Protocol   string
	Duration   time.Duration
	StatusCode int
	UserAgent  string
}

// Logger handles all application logging with formatted output
type Logger struct {
	startTime time.Time
}

// NewLogger creates a new logger instance
func NewLogger() *Logger {
	return &Logger{
		startTime: time.Now(),
	}
}

// LogRequest logs HTTP request details in a formatted table
func (l *Logger) LogRequest(entry LogEntry) {
	// Print table header if it's the first log
	if l.startTime.Equal(entry.Timestamp) {
		l.printTableHeader()
	}

	// Format the log entry as a table row
	l.printTableRow(entry)
}

// printTableHeader prints the table header with column separators
func (l *Logger) printTableHeader() {
	header := fmt.Sprintf("%-25s | %-7s | %-15s | %-30s | %-8s | %-10s | %-3s",
		"Timestamp", "Method", "Remote Address", "Path", "Protocol", "Duration", "Status")

	// Print header
	fmt.Println(strings.Repeat("=", len(header)))
	fmt.Println(header)
	fmt.Println(strings.Repeat("=", len(header)))
}

// printTableRow prints a single log entry as a table row
func (l *Logger) printTableRow(entry LogEntry) {
	// Format timestamp
	timestamp := entry.Timestamp.Format("2006-01-02 15:04:05")

	// Format remote address (remove port if present)
	remoteAddr := entry.RemoteAddr
	if strings.Contains(remoteAddr, ":") {
		remoteAddr = strings.Split(remoteAddr, ":")[0]
	}

	// Truncate path if too long
	path := entry.Path
	if len(path) > 28 {
		path = path[:25] + "..."
	}

	// Format duration
	duration := entry.Duration.String()
	if len(duration) > 8 {
		duration = duration[:8]
	}

	// Format status code with color
	statusColor := getStatusColor(entry.StatusCode)
	statusStr := fmt.Sprintf("%d", entry.StatusCode)

	// Print the formatted row
	row := fmt.Sprintf("%-25s | %-7s | %-15s | %-30s | %-8s | %-10s | %s",
		timestamp,
		entry.Method,
		remoteAddr,
		path,
		entry.Protocol,
		duration,
		statusColor+statusStr+"\033[0m",
	)

	fmt.Println(row)
}

// getStatusColor returns ANSI color codes based on HTTP status code
func getStatusColor(statusCode int) string {
	switch {
	case statusCode >= 200 && statusCode < 300:
		return "\033[32m" // Green for success
	case statusCode >= 300 && statusCode < 400:
		return "\033[36m" // Cyan for redirect
	case statusCode >= 400 && statusCode < 500:
		return "\033[33m" // Yellow for client error
	case statusCode >= 500:
		return "\033[31m" // Red for server error
	default:
		return "\033[0m" // Default color
	}
}

// LogStartup logs application startup information
func (l *Logger) LogStartup(port string) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("🚀 RGP Backend Enquiry API Starting...")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("📍 Server Address: http://localhost:%s\n", port)
	fmt.Printf("⏰ Start Time: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf("🌐 Environment: %s\n", getEnvironment())
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("📊 Request Logs:")
	fmt.Println()
}

// LogShutdown logs application shutdown information
func (l *Logger) LogShutdown() {
	fmt.Println()
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("🛑 Server Shutting Down...")
	fmt.Printf("⏰ Shutdown Time: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf("⏱️  Total Uptime: %v\n", time.Since(l.startTime))
	fmt.Println(strings.Repeat("=", 60))
}

// LogDatabaseConnection logs database connection status
func (l *Logger) LogDatabaseConnection(dbName string, success bool) {
	if success {
		fmt.Printf("✅ Database Connected: %s\n", dbName)
	} else {
		fmt.Printf("❌ Database Connection Failed: %s\n", dbName)
	}
}

// getEnvironment returns the current environment
func getEnvironment() string {
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development"
	}
	return env
}
