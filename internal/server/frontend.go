package server

const indexHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>LocalSend - File Sharing</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 20px;
        }

        .container {
            max-width: 800px;
            margin: 0 auto;
            background: white;
            border-radius: 15px;
            box-shadow: 0 20px 40px rgba(0,0,0,0.1);
            overflow: hidden;
        }

        .header {
            background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
            color: white;
            padding: 30px;
            text-align: center;
        }

        .header h1 {
            font-size: 2.5em;
            margin-bottom: 10px;
        }

        .header p {
            opacity: 0.9;
            font-size: 1.1em;
        }

        .content {
            padding: 30px;
        }

        .section {
            margin-bottom: 30px;
            padding: 20px;
            border: 2px dashed #e0e0e0;
            border-radius: 10px;
            transition: all 0.3s ease;
        }

        .section:hover {
            border-color: #4facfe;
            background: #f8f9ff;
        }

        .section h2 {
            color: #333;
            margin-bottom: 15px;
            font-size: 1.5em;
        }

        .btn {
            background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
            color: white;
            border: none;
            padding: 12px 25px;
            border-radius: 25px;
            cursor: pointer;
            font-size: 16px;
            font-weight: 600;
            transition: all 0.3s ease;
            margin: 5px;
        }

        .btn:hover {
            transform: translateY(-2px);
            box-shadow: 0 5px 15px rgba(79, 172, 254, 0.4);
        }

        .btn:disabled {
            opacity: 0.6;
            cursor: not-allowed;
            transform: none;
        }

        .file-input {
            margin: 10px 0;
        }

        .file-input input[type="file"] {
            display: none;
        }

        .file-input label {
            display: inline-block;
            background: #f0f0f0;
            border: 2px dashed #ccc;
            padding: 20px;
            border-radius: 10px;
            cursor: pointer;
            text-align: center;
            width: 100%;
            transition: all 0.3s ease;
        }

        .file-input label:hover {
            border-color: #4facfe;
            background: #f8f9ff;
        }

        .devices-list {
            margin-top: 15px;
        }

        .device {
            background: #f8f9fa;
            border: 1px solid #e9ecef;
            border-radius: 8px;
            padding: 15px;
            margin: 10px 0;
            cursor: pointer;
            transition: all 0.3s ease;
        }

        .device:hover {
            background: #e3f2fd;
            border-color: #4facfe;
        }

        .device.selected {
            background: #e3f2fd;
            border-color: #4facfe;
            box-shadow: 0 2px 8px rgba(79, 172, 254, 0.3);
        }

        .device-name {
            font-weight: 600;
            color: #333;
        }

        .device-ip {
            color: #666;
            font-size: 0.9em;
        }

        .status {
            margin-top: 15px;
            padding: 10px;
            border-radius: 5px;
            display: none;
        }

        .status.success {
            background: #d4edda;
            color: #155724;
            border: 1px solid #c3e6cb;
        }

        .status.error {
            background: #f8d7da;
            color: #721c24;
            border: 1px solid #f5c6cb;
        }

        .status.info {
            background: #d1ecf1;
            color: #0c5460;
            border: 1px solid #bee5eb;
        }

        .file-list {
            margin-top: 10px;
        }

        .file-item {
            background: #f8f9fa;
            padding: 8px 12px;
            margin: 5px 0;
            border-radius: 5px;
            border-left: 4px solid #4facfe;
        }

        .loading {
            display: inline-block;
            width: 20px;
            height: 20px;
            border: 3px solid #f3f3f3;
            border-top: 3px solid #4facfe;
            border-radius: 50%;
            animation: spin 1s linear infinite;
            margin-right: 10px;
        }

        @keyframes spin {
            0% { transform: rotate(0deg); }
            100% { transform: rotate(360deg); }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🚀 LocalSend</h1>
            <p>Kirim file dengan mudah di jaringan lokal</p>
        </div>
        
        <div class="content">
            <!-- Device Discovery Section -->
            <div class="section">
                <h2>🔍 Cari Perangkat</h2>
                <p>Temukan perangkat lain di jaringan yang sama</p>
                <button class="btn" onclick="discoverDevices()" id="discoverBtn">
                    Cari Perangkat
                </button>
                <div class="devices-list" id="devicesList"></div>
            </div>

            <!-- File Selection Section -->
            <div class="section">
                <h2>📁 Pilih File</h2>
                <p>Pilih file yang ingin dikirim</p>
                <div class="file-input">
                    <label for="fileInput">
                        📎 Klik untuk memilih file atau drag & drop di sini
                    </label>
                    <input type="file" id="fileInput" multiple>
                </div>
                <div class="file-list" id="fileList"></div>
            </div>

            <!-- Send Section -->
            <div class="section">
                <h2>📤 Kirim File</h2>
                <p>Kirim file ke perangkat yang dipilih</p>
                <button class="btn" onclick="sendFiles()" id="sendBtn" disabled>
                    Kirim File
                </button>
            </div>

            <!-- Status Section -->
            <div class="status" id="status"></div>
        </div>
    </div>

    <script>
        let selectedDevice = null;
        let selectedFiles = [];
        let discoveredDevices = [];

        // File input handler
        document.getElementById('fileInput').addEventListener('change', function(e) {
            handleFiles(e.target.files);
        });

        // Drag and drop handlers
        const fileInput = document.querySelector('.file-input label');
        
        fileInput.addEventListener('dragover', function(e) {
            e.preventDefault();
            this.style.borderColor = '#4facfe';
            this.style.background = '#f8f9ff';
        });

        fileInput.addEventListener('dragleave', function(e) {
            e.preventDefault();
            this.style.borderColor = '#ccc';
            this.style.background = '#f0f0f0';
        });

        fileInput.addEventListener('drop', function(e) {
            e.preventDefault();
            this.style.borderColor = '#ccc';
            this.style.background = '#f0f0f0';
            handleFiles(e.dataTransfer.files);
        });

        function handleFiles(files) {
            selectedFiles = Array.from(files);
            displaySelectedFiles();
            updateSendButton();
        }

        function displaySelectedFiles() {
            const fileList = document.getElementById('fileList');
            fileList.innerHTML = '';
            
            selectedFiles.forEach(file => {
                const fileItem = document.createElement('div');
                fileItem.className = 'file-item';
                fileItem.innerHTML = '<strong>' + file.name + '</strong> ' +
                    '<span style="color: #666;">(' + formatFileSize(file.size) + ')</span>';
                fileList.appendChild(fileItem);
            });
        }

        function formatFileSize(bytes) {
            if (bytes === 0) return '0 Bytes';
            const k = 1024;
            const sizes = ['Bytes', 'KB', 'MB', 'GB'];
            const i = Math.floor(Math.log(bytes) / Math.log(k));
            return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
        }

        async function discoverDevices() {
            const btn = document.getElementById('discoverBtn');
            const originalText = btn.innerHTML;
            
            btn.innerHTML = '<span class="loading"></span>Mencari...';
            btn.disabled = true;
            
            showStatus('Mencari perangkat di jaringan...', 'info');
            
            try {
                const response = await fetch('/api/discover', {
                    method: 'POST'
                });
                
                const data = await response.json();
                
                if (data.success) {
                    discoveredDevices = data.devices || [];
                    displayDevices();
                    showStatus('Ditemukan ' + discoveredDevices.length + ' perangkat', 'success');
                } else {
                    showStatus('Gagal mencari perangkat', 'error');
                }
            } catch (error) {
                showStatus('Error: ' + error.message, 'error');
            } finally {
                btn.innerHTML = originalText;
                btn.disabled = false;
            }
        }

        function displayDevices() {
            const devicesList = document.getElementById('devicesList');
            devicesList.innerHTML = '';
            
            if (discoveredDevices.length === 0) {
                devicesList.innerHTML = '<p style="color: #666; text-align: center; padding: 20px;">Tidak ada perangkat ditemukan</p>';
                return;
            }
            
            discoveredDevices.forEach((device, index) => {
                const deviceElement = document.createElement('div');
                deviceElement.className = 'device';
                deviceElement.onclick = () => selectDevice(index);
                deviceElement.innerHTML = '<div class="device-name">' + device.name + '</div>' +
                    '<div class="device-ip">' + device.ip + ':' + device.port + '</div>';
                devicesList.appendChild(deviceElement);
            });
        }

        function selectDevice(index) {
            // Remove previous selection
            document.querySelectorAll('.device').forEach(d => d.classList.remove('selected'));
            
            // Select new device
            document.querySelectorAll('.device')[index].classList.add('selected');
            selectedDevice = discoveredDevices[index];
            
            updateSendButton();
            showStatus('Perangkat dipilih: ' + selectedDevice.name, 'info');
        }

        function updateSendButton() {
            const sendBtn = document.getElementById('sendBtn');
            sendBtn.disabled = !selectedDevice || selectedFiles.length === 0;
        }

        async function sendFiles() {
            if (!selectedDevice || selectedFiles.length === 0) {
                showStatus('Pilih perangkat dan file terlebih dahulu', 'error');
                return;
            }
            
            const btn = document.getElementById('sendBtn');
            const originalText = btn.innerHTML;
            
            btn.innerHTML = '<span class="loading"></span>Mengirim...';
            btn.disabled = true;
            
            showStatus('Mengirim file...', 'info');
            
            try {
                // First upload files to our server
                const formData = new FormData();
                selectedFiles.forEach(file => {
                    formData.append('files', file);
                });
                
                const uploadResponse = await fetch('/api/upload', {
                    method: 'POST',
                    body: formData
                });
                
                const uploadData = await uploadResponse.json();
                
                if (!uploadData.success) {
                    throw new Error('Gagal mengupload file');
                }
                
                // Then send files to target device
                const sendResponse = await fetch('/api/send', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        targetIP: selectedDevice.ip,
                        targetPort: selectedDevice.port,
                        filePaths: uploadData.files.map(f => f.path)
                    })
                });
                
                const sendData = await sendResponse.json();
                
                if (sendData.success) {
                    showStatus('File berhasil dikirim ke ' + selectedDevice.name + '!', 'success');
                    // Clear selections
                    selectedFiles = [];
                    document.getElementById('fileInput').value = '';
                    displaySelectedFiles();
                } else {
                    showStatus('Gagal mengirim file: ' + (sendData.errors || []).join(', '), 'error');
                }
            } catch (error) {
                showStatus('Error: ' + error.message, 'error');
            } finally {
                btn.innerHTML = originalText;
                updateSendButton();
            }
        }

        function showStatus(message, type) {
            const status = document.getElementById('status');
            status.className = 'status ' + type;
            status.innerHTML = message;
            status.style.display = 'block';
            
            // Auto hide after 5 seconds for success messages
            if (type === 'success') {
                setTimeout(() => {
                    status.style.display = 'none';
                }, 5000);
            }
        }

        // Auto-discover devices on page load
        window.addEventListener('load', function() {
            setTimeout(discoverDevices, 1000);
        });
    </script>
</body>
</html>`