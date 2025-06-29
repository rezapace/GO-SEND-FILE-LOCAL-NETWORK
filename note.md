Kertas Kerja: Rancangan dan Analisis Aplikasi Berbagi Berkas Jaringan Lokal dengan Go
Abstrak: Dokumen ini menguraikan desain, arsitektur, dan blueprint implementasi untuk aplikasi transfer berkas peer-to-peer (P2P) yang berjalan di jaringan area lokal (LAN). Aplikasi ini dikembangkan menggunakan bahasa pemrograman Go untuk backend karena kinerjanya yang tinggi, dukungan konkurensi yang kuat, dan pustaka standar yang kaya. Antarmuka pengguna (UI) dibuat menggunakan HTML dan JavaScript standar, memastikan kesederhanaan dan portabilitas. Fokus utama adalah pada mekanisme penemuan perangkat (device discovery) yang efisien dan proses transfer berkas yang cepat dan andal.

1. Analisis Kebutuhan dan Tujuan
Tujuan utama proyek ini adalah menciptakan solusi yang sederhana dan cepat untuk mengirim berkas antar perangkat dalam satu jaringan lokal tanpa memerlukan server pusat atau koneksi internet.

Kebutuhan Fungsional:

Penemuan Perangkat: Aplikasi harus dapat secara otomatis mendeteksi instance lain dari aplikasi yang berjalan di perangkat lain dalam jaringan yang sama.

Pemilihan Berkas: Pengguna harus dapat memilih satu atau lebih berkas dari perangkat mereka untuk dikirim.

Pengiriman Berkas: Aplikasi harus dapat mentransfer berkas yang dipilih ke perangkat tujuan yang terdeteksi.

Penerimaan Berkas: Aplikasi di perangkat tujuan harus dapat menerima dan menyimpan berkas yang masuk.

Kebutuhan Non-Fungsional:

Kecepatan: Transfer harus memanfaatkan kecepatan penuh dari jaringan lokal.

Efisiensi: Penggunaan sumber daya (CPU, memori) harus minimal.

Kesederhanaan: Antarmuka harus intuitif dan tidak memerlukan konfigurasi yang rumit (misalnya, memasukkan alamat IP secara manual).

Portabilitas: Aplikasi harus dapat berjalan di berbagai sistem operasi (Windows, macOS, Linux).

2. Arsitektur dan Desain Teknologi
Untuk memenuhi kebutuhan tersebut, kami mengusulkan arsitektur hibrida yang menggabungkan dua protokol jaringan utama untuk tugas yang berbeda: UDP untuk penemuan perangkat dan TCP (via HTTP) untuk transfer berkas.

Komponen Utama:

Backend (Go): Sebuah executable tunggal yang menangani semua logika inti.

UDP Broadcaster/Listener: Bertugas untuk mengirimkan pesan "halo" ke seluruh jaringan untuk mengumumkan keberadaannya dan mendengarkan pesan serupa dari perangkat lain.

HTTP Server: Berjalan di atas TCP, server ini menyediakan endpoint untuk menerima unggahan berkas. Ini adalah pilihan yang sangat baik karena Go memiliki server HTTP yang sangat andal dan beperforma tinggi di pustaka standarnya.

Frontend (HTML/JavaScript): Sebuah halaman web sederhana yang berkomunikasi dengan backend Go lokal.

Menyediakan tombol untuk memulai penemuan perangkat.

Menampilkan daftar perangkat yang ditemukan.

Menyediakan form untuk memilih berkas dan memilih tujuan.

Menggunakan fetch API di JavaScript untuk mengirim data berkas ke perangkat tujuan melalui HTTP POST.

3. Alur Kerja (Blueprint)
Berikut adalah alur kerja langkah demi langkah dari aplikasi:

A. Fase Penemuan Perangkat (Device Discovery)

Inisiasi: Pengguna di Perangkat A menekan tombol "Cari Perangkat" di antarmuka web.

Permintaan ke Backend: JavaScript mengirim permintaan ke backend Go di Perangkat A (misalnya, ke http://localhost:8080/discover).

Siaran UDP: Backend Go di Perangkat A mengirimkan paket UDP broadcast ke seluruh subnet lokal (misalnya, ke alamat 255.255.255.255 pada port tertentu, katakanlah 8888). Paket ini berisi informasi identifikasi, seperti nama perangkat atau alamat IP.

Mendengarkan Respons: Semua instance lain dari aplikasi (misalnya, di Perangkat B) terus-menerus mendengarkan di port UDP 8888.

Respons UDP: Ketika Perangkat B menerima pesan siaran dari Perangkat A, ia merespons dengan mengirimkan paket UDP langsung kembali ke alamat IP Perangkat A. Respons ini berisi informasi tentang Perangkat B.

Pengumpulan Peer: Backend di Perangkat A mengumpulkan semua respons yang masuk dalam jangka waktu singkat (misalnya, 3 detik).

Menampilkan Hasil: Daftar perangkat yang merespons (peers) dikirim kembali ke frontend Perangkat A dan ditampilkan kepada pengguna.

B. Fase Transfer Berkas

Seleksi: Pengguna di Perangkat A memilih berkas yang akan dikirim dan memilih Perangkat B dari daftar perangkat yang ditemukan.

Permintaan HTTP POST: JavaScript di frontend Perangkat A membuat objek FormData yang berisi data berkas. Kemudian, ia mengirimkan permintaan HTTP POST ke endpoint unggah di Perangkat B. Alamat tujuannya adalah alamat IP Perangkat B yang diperoleh selama fase penemuan (misalnya, http://<IP_Perangkat_B>:8080/upload).

Penerimaan Berkas: Server HTTP di Perangkat B menerima permintaan POST. Handler di Go akan mem-parsing multipart/form-data, mengekstrak berkas, dan menyimpannya ke direktori lokal (misalnya, folder Downloads/LocalSend).

Umpan Balik: Setelah berkas berhasil disimpan, server di Perangkat B mengirimkan respons HTTP 200 OK. Frontend di Perangkat A menerima respons ini dan menampilkan pesan sukses kepada pengguna.

4. Analisis Kinerja dan Efisiensi
Bahasa Go: Sebagai bahasa yang dikompilasi, Go menawarkan kinerja yang mendekati C/C++ dengan manajemen memori yang aman dan konkurensi yang jauh lebih mudah. Goroutine memungkinkan backend untuk menangani penemuan UDP dan permintaan server HTTP secara bersamaan tanpa memblokir.

UDP untuk Penemuan: UDP bersifat connectionless dan ringan. Mengirim satu paket siaran jauh lebih efisien untuk penemuan daripada mencoba membuat koneksi TCP ke setiap kemungkinan IP di jaringan.

HTTP untuk Transfer: Meskipun ada sedikit overhead pada header HTTP, protokol ini sangat andal (karena berjalan di atas TCP), didukung secara universal, dan sangat mudah diimplementasikan baik di backend Go maupun frontend JavaScript. Ini menyederhanakan kode secara signifikan.

Sumber Daya Minimal: Dengan tidak adanya dependensi eksternal yang berat (hanya pustaka standar Go) dan antarmuka web yang sederhana, aplikasi ini akan memiliki jejak memori dan penggunaan CPU yang sangat rendah.

5. Kesimpulan
Arsitektur yang diusulkan ini memberikan keseimbangan yang sangat baik antara kesederhanaan implementasi, efisiensi sumber daya, dan kecepatan transfer. Dengan memanfaatkan kekuatan Go untuk tugas-tugas jaringan tingkat rendah dan HTTP untuk transfer data yang andal, kita dapat membangun aplikasi berbagi berkas lokal yang kuat dan mudah digunakan. Prototipe yang menyertai dokumen ini akan menunjukkan kelayakan dan fungsionalitas dari desain ini.