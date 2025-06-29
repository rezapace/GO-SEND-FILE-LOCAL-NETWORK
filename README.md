# LocalSend - Aplikasi Berbagi File Jaringan Lokal

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20macOS%20%7C%20Linux-lightgrey.svg)]()

LocalSend adalah aplikasi berbagi file peer-to-peer yang dirancang khusus untuk jaringan area lokal (LAN). Aplikasi ini memungkinkan transfer file yang cepat dan aman antar perangkat dalam satu jaringan tanpa memerlukan koneksi internet atau server pusat.

## 📋 Daftar Isi

- [Fitur Utama](#-fitur-utama)
- [Persyaratan Sistem](#-persyaratan-sistem)
- [Instalasi dan Penggunaan](#-instalasi-dan-penggunaan)
- [Arsitektur Aplikasi](#-arsitektur-aplikasi)
- [Dokumentasi API](#-dokumentasi-api)
- [Konfigurasi](#-konfigurasi)
- [Keamanan](#-keamanan)
- [Troubleshooting](#-troubleshooting)
- [Pengembangan](#-pengembangan)
- [Kontribusi](#-kontribusi)
- [Lisensi](#-lisensi)

## ✨ Fitur Utama

### 🔍 **Penemuan Perangkat Otomatis**
- Deteksi otomatis perangkat lain di jaringan menggunakan UDP broadcast
- Tidak memerlukan konfigurasi manual alamat IP
- Real-time discovery dengan refresh otomatis

### 📁 **Transfer File Multi-Platform**
- Mendukung transfer multiple file sekaligus
- Drag & drop interface yang intuitif
- Penanganan duplikasi nama file otomatis
- Progress tracking untuk setiap transfer

### 🚀 **Performa Tinggi**
- Transfer langsung antar perangkat tanpa server perantara
- Memanfaatkan kecepatan penuh jaringan lokal
- Penggunaan memori dan CPU yang minimal
- Concurrent handling untuk multiple connections

### 🌐 **Cross-Platform Compatibility**
- Windows (x64)
- macOS (Intel & Apple Silicon)
- Linux (x64)
- Web interface yang universal

### 💻 **Antarmuka Web Modern**
- Responsive design untuk berbagai ukuran layar
- Real-time status updates
- Material design inspired UI
- Embedded frontend tanpa dependensi eksternal

## 🔧 Persyaratan Sistem

### Minimum Requirements
- **Go**: Version 1.21 atau lebih baru
- **RAM**: 64 MB
- **Storage**: 10 MB ruang kosong
- **Network**: Koneksi ke jaringan lokal (WiFi/Ethernet)

### Supported Operating Systems
- Windows 10/11 (x64)
- macOS 10.15+ (Catalina atau lebih baru)
- Linux distributions dengan kernel 3.2+

### Network Requirements
- Semua perangkat harus berada dalam subnet yang sama
- Port 8080 (HTTP) dan 8888 (UDP) harus tersedia
- Firewall harus mengizinkan komunikasi pada port tersebut

## 🚀 Instalasi dan Penggunaan

### Metode 1: Menjalankan dari Source Code

```bash
# 1. Clone repository
git clone <repository-url>
cd "sendfile local network"

# 2. Verifikasi Go installation
go version

# 3. Download dependencies
go mod tidy

# 4. Jalankan aplikasi
go run main.go
```

### Metode 2: Build Executable

```bash
# Build untuk sistem operasi saat ini
go build -ldflags "-s -w" -o localsend main.go

# Jalankan executable
./localsend
```

### Metode 3: Cross-Platform Build

```bash
# Windows (64-bit)
GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o localsend.exe main.go

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o localsend-mac-intel main.go

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o localsend-mac-arm main.go

# Linux (64-bit)
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o localsend-linux main.go
```

### Penggunaan Aplikasi

#### 1. **Memulai Aplikasi**
```bash
./localsend
```
Output yang diharapkan:
```
Starting LocalSend application...
HTTP Server: http://localhost:8080
UDP Discovery Port: 8888
Discovery service started on UDP port 8888
HTTP server starting on port 8080

Application is ready!
Open your browser and go to: http://localhost:8080
Press Ctrl+C to stop the application
```

#### 2. **Mengakses Web Interface**
- Buka browser web
- Navigasi ke `http://localhost:8080`
- Interface akan menampilkan dashboard utama

#### 3. **Mengirim File**
1. **Penemuan Perangkat**:
   - Klik tombol "Cari Perangkat"
   - Tunggu hingga daftar perangkat muncul (3-5 detik)
   - Perangkat yang ditemukan akan ditampilkan dengan nama dan IP

2. **Pemilihan File**:
   - Klik area "Pilih File" atau gunakan drag & drop
   - Pilih satu atau multiple file
   - File yang dipilih akan ditampilkan dalam daftar

3. **Transfer File**:
   - Pilih perangkat tujuan dari daftar
   - Klik tombol "Kirim File"
   - Monitor progress transfer

#### 4. **Menerima File**
File yang diterima akan otomatis disimpan di:
- **Windows**: `%USERPROFILE%\Downloads\LocalSend\`
- **macOS**: `~/Downloads/LocalSend/`
- **Linux**: `~/Downloads/LocalSend/`

## 🏗️ Arsitektur Aplikasi

### Gambaran Umum
LocalSend menggunakan arsitektur hibrida yang menggabungkan protokol UDP untuk device discovery dan HTTP untuk transfer file.

```
┌─────────────────┐    UDP Broadcast    ┌─────────────────┐
│   Device A      │◄──────────────────►│   Device B      │
│                 │     (Port 8888)     │                 │
│  ┌───────────┐  │                     │  ┌───────────┐  │
│  │Discovery  │  │                     │  │Discovery  │  │
│  │Service    │  │                     │  │Service    │  │
│  └───────────┘  │                     │  └───────────┘  │
│  ┌───────────┐  │    HTTP Transfer    │  ┌───────────┐  │
│  │HTTP       │  │◄──────────────────►│  │HTTP       │  │
│  │Server     │  │     (Port 8080)     │  │Server     │  │
│  └───────────┘  │                     │  └───────────┘  │
│  ┌───────────┐  │                     │  ┌───────────┐  │
│  │Web        │  │                     │  │Web        │  │
│  │Frontend   │  │                     │  │Frontend   │  │
│  └───────────┘  │                     │  └───────────┘  │
└─────────────────┘                     └─────────────────┘
```

### Komponen Utama

#### 1. **Discovery Service** (`internal/discovery/`)
- **Fungsi**: Menangani penemuan perangkat menggunakan UDP broadcast
- **Port**: 8888 (UDP)
- **Protokol**: JSON over UDP
- **Features**:
  - Broadcast discovery messages
  - Listen untuk discovery requests
  - Maintain peer list dengan automatic cleanup
  - Local IP detection

#### 2. **HTTP Server** (`internal/server/`)
- **Fungsi**: Menangani transfer file dan web interface
- **Port**: 8080 (HTTP)
- **Endpoints**:
  - `GET /` - Web interface
  - `POST /api/discover` - Trigger device discovery
  - `GET /api/peers` - Get discovered devices
  - `POST /api/upload` - Upload files from frontend
  - `POST /api/send` - Send files to target device
  - `POST /upload` - Receive files from other devices

#### 3. **Configuration Management** (`internal/config/`)
- **Fungsi**: Mengelola konfigurasi aplikasi
- **Features**:
  - Auto-detection sistem operasi
  - Dynamic download directory creation
  - Hostname detection untuk device naming

#### 4. **Web Frontend** (`internal/server/frontend.go`)
- **Teknologi**: HTML5, CSS3, JavaScript (ES6+)
- **Features**:
  - Responsive design
  - Drag & drop file upload
  - Real-time device discovery
  - Progress tracking
  - Error handling

### Struktur Direktori

```
localsend/
├── main.go                     # Entry point aplikasi
├── go.mod                      # Go module definition
├── go.sum                      # Dependency checksums
├── README.md                   # Dokumentasi lengkap
├── LICENSE                     # MIT License
├── note.md                     # Design document
└── internal/                   # Internal packages
    ├── config/
    │   └── config.go          # Configuration management
    ├── discovery/
    │   └── discovery.go       # UDP device discovery
    └── server/
        ├── server.go          # HTTP server implementation
        ├── multipart.go       # Multipart form utilities
        └── frontend.go        # Embedded web interface
```

### Flow Diagram

#### Device Discovery Flow
```
[User clicks "Discover"] 
        ↓
[Frontend sends POST to /api/discover]
        ↓
[Backend broadcasts UDP message]
        ↓
[Other devices receive broadcast]
        ↓
[Other devices send UDP response]
        ↓
[Backend collects responses]
        ↓
[Device list returned to frontend]
```

#### File Transfer Flow
```
[User selects files and target device]
        ↓
[Files uploaded to local temp storage]
        ↓
[Frontend sends transfer request]
        ↓
[Backend creates HTTP POST to target]
        ↓
[Target device receives files]
        ↓
[Files saved to download directory]
        ↓
[Success response sent back]
```

## 📡 Dokumentasi API

### REST Endpoints

#### `GET /`
**Deskripsi**: Menampilkan web interface utama

**Response**: HTML page

#### `POST /api/discover`
**Deskripsi**: Memulai proses device discovery

**Request**: Empty body

**Response**:
```json
{
  "success": true,
  "devices": [
    {
      "name": "MacBook-Pro",
      "ip": "192.168.1.100",
      "port": 8080
    }
  ]
}
```

#### `GET /api/peers`
**Deskripsi**: Mendapatkan daftar perangkat yang sudah ditemukan

**Response**:
```json
{
  "success": true,
  "peers": [
    {
      "name": "Windows-PC",
      "ip": "192.168.1.101",
      "port": 8080
    }
  ]
}
```

#### `POST /api/upload`
**Deskripsi**: Upload file dari frontend untuk persiapan pengiriman

**Request**: `multipart/form-data` dengan field `files`

**Response**:
```json
{
  "success": true,
  "files": [
    {
      "name": "document.pdf",
      "size": 1024000,
      "path": "/tmp/localsend_temp/document.pdf"
    }
  ]
}
```

#### `POST /api/send`
**Deskripsi**: Mengirim file ke perangkat target

**Request**:
```json
{
  "targetIP": "192.168.1.101",
  "targetPort": 8080,
  "filePaths": ["/tmp/localsend_temp/document.pdf"]
}
```

**Response**:
```json
{
  "success": true,
  "errors": []
}
```

#### `POST /upload`
**Deskripsi**: Menerima file dari perangkat lain

**Request**: `multipart/form-data` dengan field `files`

**Response**:
```json
{
  "success": true,
  "message": "Received 1 files",
  "files": ["document.pdf"]
}
```

### UDP Protocol

#### Discovery Message Format
```json
{
  "type": "discover",
  "deviceName": "MacBook-Pro",
  "ip": "192.168.1.100",
  "port": 8080
}
```

#### Response Message Format
```json
{
  "type": "response",
  "deviceName": "Windows-PC",
  "ip": "192.168.1.101",
  "port": 8080
}
```

## ⚙️ Konfigurasi

### Konfigurasi Default
Aplikasi menggunakan konfigurasi default yang dapat ditemukan di `internal/config/config.go`:

```go
type Config struct {
    HTTPPort    int    // 8080
    UDPPort     int    // 8888
    DeviceName  string // Hostname sistem
    DownloadDir string // ~/Downloads/LocalSend/
}
```

### Kustomisasi Konfigurasi

#### 1. **Mengubah Port**
```go
// internal/config/config.go
return &Config{
    HTTPPort:    9090,  // Ubah dari 8080
    UDPPort:     9999,  // Ubah dari 8888
    DeviceName:  deviceName,
    DownloadDir: downloadDir,
}
```

#### 2. **Mengubah Download Directory**
```go
// Contoh: Gunakan direktori custom
downloadDir := "/custom/download/path"
```

#### 3. **Mengubah Device Name**
```go
// Contoh: Set nama device custom
deviceName := "My-Custom-Device"
```

### Environment Variables (Future Enhancement)
Untuk versi mendatang, konfigurasi dapat dilakukan melalui environment variables:

```bash
export LOCALSEND_HTTP_PORT=9090
export LOCALSEND_UDP_PORT=9999
export LOCALSEND_DEVICE_NAME="Custom-Device"
export LOCALSEND_DOWNLOAD_DIR="/custom/path"
```

## 🔒 Keamanan

### ⚠️ Peringatan Keamanan

**PENTING**: LocalSend dirancang khusus untuk jaringan lokal yang terpercaya. Aplikasi ini **TIDAK** menyediakan:

- Enkripsi data dalam transit
- Autentikasi pengguna
- Otorisasi akses file
- Proteksi terhadap man-in-the-middle attacks

### Rekomendasi Penggunaan

#### ✅ **Aman untuk digunakan**:
- Jaringan rumah pribadi
- Jaringan kantor internal yang terpercaya
- Jaringan lab atau development environment
- Hotspot pribadi

#### ❌ **TIDAK aman untuk digunakan**:
- WiFi publik (café, hotel, bandara)
- Jaringan yang tidak terpercaya
- Jaringan dengan pengguna yang tidak dikenal
- Environment production yang memerlukan keamanan tinggi

### Best Practices

1. **Network Isolation**:
   - Gunakan hanya di jaringan yang terisolasi
   - Pastikan firewall aktif di level router

2. **File Validation**:
   - Selalu scan file yang diterima dengan antivirus
   - Jangan menjalankan executable yang diterima tanpa verifikasi

3. **Access Control**:
   - Tutup aplikasi setelah selesai digunakan
   - Monitor log untuk aktivitas yang mencurigakan

4. **Regular Updates**:
   - Update aplikasi secara berkala
   - Monitor security advisories

### Future Security Enhancements

Pengembangan selanjutnya akan mencakup:
- TLS encryption untuk HTTP transfers
- Device authentication menggunakan certificates
- File integrity verification
- Access control lists
- Audit logging

## 🔧 Troubleshooting

### Masalah Umum dan Solusi

#### 1. **Perangkat Tidak Ditemukan**

**Gejala**: Daftar perangkat kosong setelah discovery

**Penyebab Umum**:
- Perangkat tidak dalam subnet yang sama
- Firewall memblokir UDP port 8888
- Aplikasi tidak berjalan di perangkat target

**Solusi**:
```bash
# 1. Verifikasi koneksi jaringan
ping <target-ip>

# 2. Check firewall (macOS)
sudo pfctl -sr | grep 8888

# 3. Check firewall (Linux)
sudo iptables -L | grep 8888

# 4. Test UDP connectivity
nc -u <target-ip> 8888

# 5. Restart aplikasi di semua perangkat
```

#### 2. **Transfer File Gagal**

**Gejala**: Error saat mengirim file

**Penyebab Umum**:
- Koneksi jaringan terputus
- Ruang disk tidak cukup
- File terlalu besar
- Target device tidak responsif

**Solusi**:
```bash
# 1. Check disk space
df -h

# 2. Test HTTP connectivity
curl -I http://<target-ip>:8080

# 3. Check file permissions
ls -la <file-path>

# 4. Monitor network traffic
netstat -an | grep 8080
```

#### 3. **Port Sudah Digunakan**

**Gejala**: Error "address already in use"

**Solusi**:
```bash
# 1. Find process using port
lsof -i :8080

# 2. Kill process if necessary
kill -9 <PID>

# 3. Atau ubah port di konfigurasi
# Edit internal/config/config.go
```

#### 4. **Permission Denied**

**Gejala**: Error saat menyimpan file

**Solusi**:
```bash
# 1. Check directory permissions
ls -la ~/Downloads/

# 2. Create directory if not exists
mkdir -p ~/Downloads/LocalSend

# 3. Fix permissions
chmod 755 ~/Downloads/LocalSend
```

### Debugging Mode

Untuk debugging yang lebih detail, tambahkan logging:

```go
// Tambahkan di main.go
import "log"

func main() {
    log.SetFlags(log.LstdFlags | log.Lshortfile)
    // ... rest of code
}
```

### Log Analysis

Monitor log output untuk mengidentifikasi masalah:

```bash
# Jalankan dengan output ke file
./localsend 2>&1 | tee localsend.log

# Analyze log
grep -i error localsend.log
grep -i "failed" localsend.log
```

## 🛠️ Pengembangan

### Setup Development Environment

```bash
# 1. Install Go 1.21+
go version

# 2. Clone repository
git clone <repository-url>
cd "sendfile local network"

# 3. Install dependencies
go mod tidy

# 4. Run tests
go test ./...

# 5. Run with hot reload (install air)
go install github.com/cosmtrek/air@latest
air
```

### Code Structure Guidelines

#### Package Organization
- `main.go`: Application entry point
- `internal/config`: Configuration management
- `internal/discovery`: UDP device discovery
- `internal/server`: HTTP server and file handling

#### Coding Standards
- Follow Go conventions (gofmt, golint)
- Use meaningful variable names
- Add comments for exported functions
- Handle errors appropriately
- Use context for cancellation

### Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./internal/discovery

# Benchmark tests
go test -bench=. ./...
```

### Performance Optimization

#### Memory Usage
```bash
# Profile memory usage
go tool pprof http://localhost:8080/debug/pprof/heap
```

#### CPU Usage
```bash
# Profile CPU usage
go tool pprof http://localhost:8080/debug/pprof/profile
```

### Adding New Features

#### 1. **New API Endpoint**
```go
// internal/server/server.go
func (s *HTTPServer) handleNewFeature(w http.ResponseWriter, r *http.Request) {
    // Implementation
}

// Register in Start() method
mux.HandleFunc("/api/new-feature", s.handleNewFeature)
```

#### 2. **New Configuration Option**
```go
// internal/config/config.go
type Config struct {
    // ... existing fields
    NewOption string
}
```

#### 3. **Frontend Enhancement**
```javascript
// internal/server/frontend.go
// Add new JavaScript functionality
```

## 🤝 Kontribusi

### Cara Berkontribusi

1. **Fork Repository**
   ```bash
   git fork <repository-url>
   ```

2. **Create Feature Branch**
   ```bash
   git checkout -b feature/new-feature
   ```

3. **Make Changes**
   - Follow coding standards
   - Add tests for new functionality
   - Update documentation

4. **Commit Changes**
   ```bash
   git commit -m "feat: add new feature"
   ```

5. **Push to Branch**
   ```bash
   git push origin feature/new-feature
   ```

6. **Create Pull Request**
   - Provide clear description
   - Include test results
   - Reference related issues

### Contribution Guidelines

#### Code Quality
- All code must pass `go fmt`
- All code must pass `go vet`
- Add unit tests for new features
- Maintain test coverage above 80%

#### Documentation
- Update README.md for new features
- Add inline code comments
- Update API documentation
- Include usage examples

#### Issue Reporting
Saat melaporkan bug, sertakan:
- Versi Go yang digunakan
- Sistem operasi dan versi
- Langkah-langkah reproduksi
- Log error yang relevan
- Expected vs actual behavior

### Roadmap

#### Version 2.0 (Planned)
- [ ] TLS encryption
- [ ] User authentication
- [ ] File compression
- [ ] Resume interrupted transfers
- [ ] Mobile app support

#### Version 2.1 (Future)
- [ ] Directory synchronization
- [ ] Real-time chat
- [ ] File versioning
- [ ] Bandwidth throttling
- [ ] Plugin system

## 📄 Lisensi

Project ini menggunakan [MIT License](LICENSE).

```
MIT License

Copyright (c) 2024 LocalSend Contributors

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

---

## 📞 Dukungan

### Community Support
- **GitHub Issues**: [Report bugs dan feature requests](https://github.com/username/localsend/issues)
- **Discussions**: [Community discussions](https://github.com/username/localsend/discussions)
- **Wiki**: [Additional documentation](https://github.com/username/localsend/wiki)

### Professional Support
Untuk dukungan enterprise atau kustomisasi khusus, hubungi tim development.

---

**Dibuat dengan ❤️ menggunakan Go dan teknologi web modern**

*LocalSend - Sharing made simple, secure, and fast.*