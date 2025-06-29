package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"localsend/internal/discovery"
	"localsend/internal/server"
	"localsend/internal/config"
)

func main() {
	// Load configuration
	cfg := config.Load()

	fmt.Printf("Starting LocalSend application...\n")
	fmt.Printf("HTTP Server: http://localhost:%d\n", cfg.HTTPPort)
	fmt.Printf("UDP Discovery Port: %d\n", cfg.UDPPort)

	// Create channels for graceful shutdown
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	// Start UDP discovery service
	discoveryService := discovery.NewService(cfg.UDPPort, cfg.DeviceName)
	go func() {
		if err := discoveryService.Start(); err != nil {
			log.Printf("Discovery service error: %v", err)
		}
	}()

	// Start HTTP server
	httpServer := server.NewHTTPServer(cfg.HTTPPort, cfg.DownloadDir, discoveryService)
	go func() {
		if err := httpServer.Start(); err != nil {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	// Wait for a few seconds to ensure services are running
	time.Sleep(2 * time.Second)
	fmt.Println("\nApplication is ready!")
	fmt.Println("Open your browser and go to: http://localhost:8080")
	fmt.Println("Press Ctrl+C to stop the application")

	// Wait for shutdown signal
	<-stopChan
	fmt.Println("\nShutting down...")

	// Graceful shutdown
	discoveryService.Stop()
	httpServer.Stop()

	fmt.Println("Application stopped.")
}