package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"localsend/internal/discovery"
)

// HTTPServer handles HTTP requests
type HTTPServer struct {
	port            int
	downloadDir     string
	discoveryService *discovery.Service
	server          *http.Server
}

// NewHTTPServer creates a new HTTP server
func NewHTTPServer(port int, downloadDir string, discoveryService *discovery.Service) *HTTPServer {
	return &HTTPServer{
		port:            port,
		downloadDir:     downloadDir,
		discoveryService: discoveryService,
	}
}

// Start starts the HTTP server
func (s *HTTPServer) Start() error {
	mux := http.NewServeMux()

	// Serve static files (frontend)
	mux.HandleFunc("/", s.handleIndex)
	mux.HandleFunc("/static/", s.handleStatic)

	// API endpoints
	mux.HandleFunc("/api/discover", s.handleDiscover)
	mux.HandleFunc("/api/peers", s.handleGetPeers)
	mux.HandleFunc("/api/upload", s.handleUpload)
	mux.HandleFunc("/api/send", s.handleSendFile)

	// File upload endpoint (for receiving files from other devices)
	mux.HandleFunc("/upload", s.handleReceiveFile)

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: mux,
	}

	fmt.Printf("HTTP server starting on port %d\n", s.port)
	return s.server.ListenAndServe()
}

// Stop stops the HTTP server
func (s *HTTPServer) Stop() {
	if s.server != nil {
		s.server.Close()
		fmt.Println("HTTP server stopped")
	}
}

// handleIndex serves the main HTML page
func (s *HTTPServer) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(indexHTML))
}

// handleStatic serves static files
func (s *HTTPServer) handleStatic(w http.ResponseWriter, r *http.Request) {
	// For now, we'll embed CSS and JS in the HTML
	http.NotFound(w, r)
}

// handleDiscover triggers device discovery
func (s *HTTPServer) handleDiscover(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	devices, err := s.discoveryService.DiscoverDevices()
	if err != nil {
		http.Error(w, fmt.Sprintf("Discovery failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"devices": devices,
	})
}

// handleGetPeers returns current list of peers
func (s *HTTPServer) handleGetPeers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	peers := s.discoveryService.GetPeers()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"peers":   peers,
	})
}

// handleUpload handles file selection from frontend
func (s *HTTPServer) handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// This endpoint is for the frontend to upload files to be sent
	// We'll store them temporarily and return file info
	err := r.ParseMultipartForm(32 << 20) // 32MB max
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		http.Error(w, "No files uploaded", http.StatusBadRequest)
		return
	}

	var uploadedFiles []map[string]interface{}

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			continue
		}
		defer file.Close()

		// Create temp directory for outgoing files
		tempDir := filepath.Join(os.TempDir(), "localsend_temp")
		os.MkdirAll(tempDir, 0755)

		// Save file temporarily
		tempPath := filepath.Join(tempDir, fileHeader.Filename)
		dst, err := os.Create(tempPath)
		if err != nil {
			continue
		}
		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			continue
		}

		fileInfo := map[string]interface{}{
			"name": fileHeader.Filename,
			"size": fileHeader.Size,
			"path": tempPath,
		}
		uploadedFiles = append(uploadedFiles, fileInfo)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"files":   uploadedFiles,
	})
}

// handleSendFile sends files to a target device
func (s *HTTPServer) handleSendFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		TargetIP   string   `json:"targetIP"`
		TargetPort int      `json:"targetPort"`
		FilePaths  []string `json:"filePaths"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Send files to target device
	success := true
	var errors []string

	for _, filePath := range request.FilePaths {
		err := s.sendFileToDevice(request.TargetIP, request.TargetPort, filePath)
		if err != nil {
			success = false
			errors = append(errors, fmt.Sprintf("Failed to send %s: %v", filepath.Base(filePath), err))
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": success,
		"errors":  errors,
	})
}

// handleReceiveFile receives files from other devices
func (s *HTTPServer) handleReceiveFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(32 << 20) // 32MB max
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		http.Error(w, "No files received", http.StatusBadRequest)
		return
	}

	var savedFiles []string

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			continue
		}
		defer file.Close()

		// Save to download directory
		destPath := filepath.Join(s.downloadDir, fileHeader.Filename)
		
		// Handle duplicate filenames
		counter := 1
		originalPath := destPath
		for {
			if _, err := os.Stat(destPath); os.IsNotExist(err) {
				break
			}
			ext := filepath.Ext(originalPath)
			name := originalPath[:len(originalPath)-len(ext)]
			destPath = fmt.Sprintf("%s_%d%s", name, counter, ext)
			counter++
		}

		dst, err := os.Create(destPath)
		if err != nil {
			continue
		}
		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			os.Remove(destPath)
			continue
		}

		savedFiles = append(savedFiles, filepath.Base(destPath))
		fmt.Printf("Received file: %s\n", destPath)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Received %d files", len(savedFiles)),
		"files":   savedFiles,
	})
}

// sendFileToDevice sends a file to a target device
func (s *HTTPServer) sendFileToDevice(targetIP string, targetPort int, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Create multipart form
	pr, pw := io.Pipe()
	mw := NewMultipartWriter(pw)

	go func() {
		defer pw.Close()
		defer mw.Close()

		part, err := mw.CreateFormFile("files", filepath.Base(filePath))
		if err != nil {
			return
		}

		_, err = io.Copy(part, file)
		if err != nil {
			return
		}
	}()

	// Send HTTP POST request
	url := fmt.Sprintf("http://%s:%d/upload", targetIP, targetPort)
	req, err := http.NewRequest("POST", url, pr)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", mw.FormDataContentType())

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status: %s", resp.Status)
	}

	fmt.Printf("Successfully sent file %s to %s:%d\n", filepath.Base(filePath), targetIP, targetPort)
	return nil
}