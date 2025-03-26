# Gateway Services API

API ini adalah aplikasi manajemen pengguna menggunakan Golang dan PostgreSQL.

## Persiapan Sebelum Memulai

1. **Pastikan Database Sudah Tersedia**
   - Pastikan Anda sudah membuat database PostgreSQL dengan nama **`db_users`**.

2. **Konfigurasi Aplikasi**
   - Buka file `config.yaml` pada aplikasi ini.
   - Sesuaikan pengaturan koneksi database dan konfigurasi lainnya sesuai dengan lingkungan Anda.

## Langkah-Langkah Menjalankan Aplikasi

### 1. Menjalankan Migrations
   - Migration digunakan untuk mengatur struktur database yang diperlukan oleh aplikasi ini.
   - Masuk ke direktori `migrations/` dengan perintah berikut:
     ```bash
     cd migrations/
     ```
   - Jalankan perintah berikut untuk menjalankan migration:
     ```bash
     go run migration.go ./sql "host=localhost port=5432 user=root dbname=db_users password=password sslmode=disable" up
     ```
   - Pastikan detail koneksi (seperti host, port, user, dan password) sesuai dengan konfigurasi database PostgreSQL Anda.

### 2. Menjalankan Aplikasi
   - Setelah konfigurasi selesai, Anda dapat menjalankan aplikasi dengan perintah:
     ```bash
     go run .
     ```
   - Aplikasi sekarang akan berjalan dan terhubung ke database `db_users`.

## Struktur Direktori

- **config/**: Menyimpan konfigurasi yang berkaitan dengan layanan pihak ketiga.
- **handlers/**: Lapisan yang menangani permintaan dari pengguna, baik dari aplikasi mobile maupun web.
- **migrations/**: Menyimpan file migration SQL dan kode untuk mengatur dan memperbarui struktur database.
- **models/**: Berisi struktur data (constructs di Golang) yang memudahkan dalam membuat kontrak untuk request dan response.
- **respository/**: Lapisan yang berfungsi khusus untuk berinteraksi dengan database, termasuk operasi pencatatan dan pengambilan data.
- **routes/**: Menyimpan definisi endpoint utama yang mengarahkan permintaan ke handler terkait.
- **usecases/**: Lapisan yang menangani logika bisnis, termasuk pengolahan data dari input pengguna atau hasil dari database.
- **util/**: Menyimpan middleware dan fungsi utilitas lainnya.
- **config.yaml**: File konfigurasi utama untuk mengatur koneksi database dan parameter aplikasi lainnya.

## Teknologi yang Digunakan

- **Golang**: Backend aplikasi utama.
- **PostgreSQL**: Database untuk menyimpan data pengguna.
