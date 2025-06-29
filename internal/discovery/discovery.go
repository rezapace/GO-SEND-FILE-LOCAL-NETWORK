package discovery

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"
)

// Device represents a discovered device
type Device struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

// Message represents UDP discovery message
type Message struct {
	Type       string `json:"type"`       // "discover" or "response"
	DeviceName string `json:"deviceName"`
	IP         string `json:"ip"`
	Port       int    `json:"port"`
}

// Service handles device discovery
type Service struct {
	udpPort    int
	deviceName string
	conn       *net.UDPConn
	peers      map[string]*Device
	mutex      sync.RWMutex
	stopChan   chan bool
	running    bool
}

// NewService creates a new discovery service
func NewService(udpPort int, deviceName string) *Service {
	return &Service{
		udpPort:    udpPort,
		deviceName: deviceName,
		peers:      make(map[string]*Device),
		stopChan:   make(chan bool),
	}
}

// Start begins the discovery service
func (s *Service) Start() error {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", s.udpPort))
	if err != nil {
		return fmt.Errorf("failed to resolve UDP address: %v", err)
	}

	s.conn, err = net.ListenUDP("udp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on UDP: %v", err)
	}

	s.running = true
	fmt.Printf("Discovery service started on UDP port %d\n", s.udpPort)

	// Start listening for messages
	go s.listen()

	// Start periodic cleanup of old peers
	go s.cleanupPeers()

	return nil
}

// Stop stops the discovery service
func (s *Service) Stop() {
	if !s.running {
		return
	}

	s.running = false
	close(s.stopChan)

	if s.conn != nil {
		s.conn.Close()
	}

	fmt.Println("Discovery service stopped")
}

// listen handles incoming UDP messages
func (s *Service) listen() {
	buffer := make([]byte, 1024)

	for s.running {
		s.conn.SetReadDeadline(time.Now().Add(1 * time.Second))
		n, addr, err := s.conn.ReadFromUDP(buffer)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				continue
			}
			if s.running {
				fmt.Printf("Error reading UDP message: %v\n", err)
			}
			continue
		}

		var msg Message
		if err := json.Unmarshal(buffer[:n], &msg); err != nil {
			fmt.Printf("Error unmarshaling message: %v\n", err)
			continue
		}

		s.handleMessage(&msg, addr)
	}
}

// handleMessage processes incoming discovery messages
func (s *Service) handleMessage(msg *Message, addr *net.UDPAddr) {
	switch msg.Type {
	case "discover":
		// Someone is looking for devices, respond with our info
		s.sendResponse(addr)
	case "response":
		// Someone responded to our discovery, add them to peers
		s.addPeer(msg, addr)
	}
}

// sendResponse sends a response to a discovery request
func (s *Service) sendResponse(addr *net.UDPAddr) {
	localIP := s.getLocalIP()
	response := Message{
		Type:       "response",
		DeviceName: s.deviceName,
		IP:         localIP,
		Port:       8080, // HTTP server port
	}

	data, err := json.Marshal(response)
	if err != nil {
		fmt.Printf("Error marshaling response: %v\n", err)
		return
	}

	_, err = s.conn.WriteToUDP(data, addr)
	if err != nil {
		fmt.Printf("Error sending response: %v\n", err)
	}
}

// addPeer adds a discovered peer to the list
func (s *Service) addPeer(msg *Message, addr *net.UDPAddr) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	device := &Device{
		Name: msg.DeviceName,
		IP:   msg.IP,
		Port: msg.Port,
	}

	s.peers[addr.IP.String()] = device
	fmt.Printf("Discovered device: %s (%s:%d)\n", device.Name, device.IP, device.Port)
}

// DiscoverDevices broadcasts a discovery message
func (s *Service) DiscoverDevices() ([]*Device, error) {
	if !s.running {
		return nil, fmt.Errorf("discovery service not running")
	}

	// Clear existing peers
	s.mutex.Lock()
	s.peers = make(map[string]*Device)
	s.mutex.Unlock()

	// Broadcast discovery message
	localIP := s.getLocalIP()
	msg := Message{
		Type:       "discover",
		DeviceName: s.deviceName,
		IP:         localIP,
		Port:       8080,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("error marshaling discovery message: %v", err)
	}

	// Broadcast to subnet
	broadcastAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("255.255.255.255:%d", s.udpPort))
	if err != nil {
		return nil, fmt.Errorf("error resolving broadcast address: %v", err)
	}

	_, err = s.conn.WriteToUDP(data, broadcastAddr)
	if err != nil {
		return nil, fmt.Errorf("error broadcasting discovery message: %v", err)
	}

	// Wait for responses
	time.Sleep(3 * time.Second)

	// Return discovered devices
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	devices := make([]*Device, 0, len(s.peers))
	for _, device := range s.peers {
		devices = append(devices, device)
	}

	return devices, nil
}

// GetPeers returns the current list of discovered peers
func (s *Service) GetPeers() []*Device {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	devices := make([]*Device, 0, len(s.peers))
	for _, device := range s.peers {
		devices = append(devices, device)
	}

	return devices
}

// cleanupPeers removes old peers periodically
func (s *Service) cleanupPeers() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// For now, we keep all peers. In a real implementation,
			// you might want to ping peers and remove unresponsive ones
		case <-s.stopChan:
			return
		}
	}
}

// getLocalIP returns the local IP address
func (s *Service) getLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "127.0.0.1"
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}