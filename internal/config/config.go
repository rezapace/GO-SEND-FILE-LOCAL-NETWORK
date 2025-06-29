package config

import (
	"os"
	"path/filepath"
	"runtime"
)

// Config holds application configuration
type Config struct {
	HTTPPort    int
	UDPPort     int
	DeviceName  string
	DownloadDir string
}

// Load returns the default configuration
func Load() *Config {
	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}

	// Create download directory path
	downloadDir := filepath.Join(homeDir, "Downloads", "LocalSend")

	// Create download directory if it doesn't exist
	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		downloadDir = "./downloads" // fallback to current directory
		os.MkdirAll(downloadDir, 0755)
	}

	// Get device name (hostname or default)
	deviceName, err := os.Hostname()
	if err != nil {
		deviceName = "Unknown-" + runtime.GOOS
	}

	return &Config{
		HTTPPort:    8080,
		UDPPort:     8888,
		DeviceName:  deviceName,
		DownloadDir: downloadDir,
	}
}

// GetLocalIP returns the local IP address (placeholder for now)
func GetLocalIP() string {
	// This will be implemented in the discovery package
	return "localhost"
}