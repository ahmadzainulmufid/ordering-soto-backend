# 🍲 Soto Lamongan - Backend RESTful API (Golang)

Dokumentasi resmi untuk **Backend RESTful API** sistem manajemen & pemesanan restoran **Soto Lamongan**. Aplikasi backend ini dibangun menggunakan bahasa pemrograman **Go (Golang)** untuk menangani autentikasi JWT, manajemen katalog produk, antrean dapur (Kitchen Display System), serta integrasi gateway pembayaran online Midtrans.

---

## 📑 Daftar Isi

- [Arsitektur & Teknologi](#-arsitektur--teknologi)
- [Struktur Proyek](#-struktur-proyek)
- [Fitur Utama Backend](#-fitur-utama-backend)
- [Persyaratan Sistem](#-persyaratan-sistem)
- [Instalasi & Panduan Lokal](#-instalasi--panduan-lokal)
- [Konfigurasi Environment (.env)](#-konfigurasi-environment-env)
- [Dokumentasi Endpoint API](#-dokumentasi-endpoint-api)
- [Deployment via Docker](#-deployment-via-docker)

---

## 🚀 Arsitektur & Teknologi

- **Language**: Go (Golang 1.22+)
- **HTTP Framework**: Gin Gonic / Fiber
- **Database**: MySQL 8.0 / PostgreSQL
- **ORM / Driver**: GORM / sqlx
- **Authentication**: JWT (JSON Web Token) dengan Authorization Bearer Header
- **Payment Gateway**: Midtrans Snap API (QRIS, E-Wallet, & Virtual Account)
- **Containerization**: Docker & Docker Compose

---

## 📁 Struktur Proyek

```text
.
├── cmd/
│   └── api/
│       └── main.go             # Entrypoint utama server API Go
├── config/                     # Konfigurasi Database & Environment Variables
├── internal/
│   ├── delivery/
│   │   └── http/               # HTTP Handlers, Routing, & Middleware JWT
│   ├── domain/                 # Struct Data Models / Entities
│   ├── repository/             # Query Database Layer (MySQL)
│   └── usecase/                # Logic Aturan Bisnis Pemesanan
├── pkg/
│   ├── jwt/                    # Helper Generator & Parsing Token JWT
│   └── midtrans/               # Client SDK Payment Gateway Midtrans
├── .env.example
├── Dockerfile
├── go.mod
└── README.md
```

---

## ✨ Fitur Utama Backend

### 1. Autentikasi & Authorization (JWT):

- Keamanan berbasis Token JWT (Authorization: Bearer <TOKEN>).
- Role Management: admin, kitchen, owner.

### 2. Katalog Produk & Status Stok:

- CRUD Produk, Kategori Menu, dan Meja Makan Resto.
- Fitur ubah ketersediaan stok menu secara real-time (is_available).

### 3. Order Pipeline & State Transitions:

- Mendukung tipe pemesanan: dine_in, takeaway, dan delivery.
- Perubahan status transisi antrean: pending ➔ cooking ➔ ready ➔ completed / cancelled.

### 4. Integration Payment Webhook:

- Integrasi Midtrans Payment Gateway.
- Handling HTTP Webhook Callback otomatis dari Midtrans untuk meng-update status pembayaran transaksi secara asynchronous.

---

## 🛠️ Persyaratan Sistem

- Go (Golang) >= 1.22
- MySQL >= 8.0 / PostgreSQL >= 14
- Git

## ⚙️ Instalasi & Panduan Lokal

### 1. Clone Repositori

```
git clone [https://github.com/ahmadzainulmufid/ordering-soto-backend.git](https://github.com/ahmadzainulmufid/ordering-soto-backend.git)
cd ordering-soto-backend
```

### 2. Install Dependencies Go

```
go mod download
```

### 3. Setup Konfigurasi Environment

```
cp .env.example .env
```

### 4. Jalankan Server Go

```
# Running mode pengembang
go run cmd/api/main.go

# Atau build file binary
go build -o server cmd/api/main.go
./server
```

Server REST API default akan berjalan di http://localhost:8080

## 🔑 Konfigurasi Environment (.env)

```
APP_NAME="Soto Lamongan Go API"
APP_ENV=development
APP_PORT=8080

# Database Settings
DB_DRIVER=mysql
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=soto_lamongan_db
DB_USER=root
DB_PASSWORD=secret

# JWT Authentication Secret Key
JWT_SECRET=your_super_secret_jwt_key_here

# Midtrans Integration Keys
MIDTRANS_SERVER_KEY=your_midtrans_server_key
MIDTRANS_CLIENT_KEY=your_midtrans_client_key
MIDTRANS_IS_PRODUCTION=false

```

## 📌 Dokumentasi Endpoint API

Semua request dan response menggunakan standar format JSON (`Content-Type: application/json`).

### 1. Auth Endpoint (`/api/auth`)

| Method | Endpoint          | Deskripsi                              | Auth Required |
| :----- | :---------------- | :------------------------------------- | :------------ |
| `POST` | `/api/auth/login` | Login user & dapatkan Bearer Token JWT | Public        |
| `GET`  | `/api/auth/me`    | Ambil data profil akun terautentikasi  | JWT           |

### 2. Catalog Endpoint (`/api/products` & `/api/categories`)

| Method | Endpoint            | Deskripsi                          | Auth Required       |
| :----- | :------------------ | :--------------------------------- | :------------------ |
| `GET`  | `/api/products`     | Ambil seluruh katalog produk aktif | Public              |
| `POST` | `/api/products`     | Tambah produk menu baru            | JWT (Admin)         |
| `PUT`  | `/api/products/:id` | Update detail/stok produk menu     | JWT (Admin/Kitchen) |
| `GET`  | `/api/categories`   | Ambil seluruh kategori menu        | Public              |
| `GET`  | `/api/tables`       | Ambil seluruh daftar nomor meja    | Public              |

### 3. Order Endpoint (`/api/orders`)

| Method  | Endpoint                       | Deskripsi                                       | Auth Required             |
| :------ | :----------------------------- | :---------------------------------------------- | :------------------------ |
| `POST`  | `/api/orders`                  | Buat transaksi pesanan baru dari pelanggan      | Public                    |
| `GET`   | `/api/admin/orders`            | Ambil semua riwayat transaksi pesanan           | JWT (Admin/Kitchen/Owner) |
| `PATCH` | `/api/admin/orders/:id/status` | Update status pesanan (`cooking`, `ready`, dll) | JWT (Admin/Kitchen)       |

### 4. Midtrans Webhook Callback (`/api/payment`)

| Method | Endpoint                    | Deskripsi                                        | Auth Required |
| :----- | :-------------------------- | :----------------------------------------------- | :------------ |
| `POST` | `/api/payment/notification` | Listener webhook callback otomatis dari Midtrans | Webhook Sign  |
