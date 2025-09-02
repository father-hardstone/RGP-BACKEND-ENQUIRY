package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"syedibrahimshah067/RGP-BACKEND-ENQUIRY/main/app"
)

// main is the entry point of the application
// Initializes the app, starts the server, and handles graceful shutdown
func main() {
	// Create and configure the application
	appInstance, err := app.NewApp()
	if err != nil {
		log.Fatal("Failed to create application:", err)
	}

	// Setup graceful shutdown
	setupGracefulShutdown(appInstance)

	// Start the server
	if err := appInstance.Start(); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

// setupGracefulShutdown configures signal handling for graceful shutdown
// Listens for SIGINT and SIGTERM signals to shutdown the application cleanly
func setupGracefulShutdown(appInstance *app.App) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Received shutdown signal...")
		appInstance.Shutdown()
		os.Exit(0)
	}()
}
