# 📚 Book API - Golang Project

Sebuah RESTful API sederhana untuk mengelola data buku menggunakan Golang. Cocok untuk pembelajaran konsep handler, routing, testing, dan penerapan clean code.

## ✨ Fitur

- CRUD Buku (Create, Read, Update, Delete)
- Penyimpanan data di memori (map)
- Penanganan ID otomatis (`auto-increment`)
- Unit testing dengan `net/http/httptest`
- Routing menggunakan `go-chi/chi/v5`
- Singleton pattern untuk in-memory storage

## 📁 Struktur Folder

```
book-api/
├── handler/         # Handler HTTP (CreateBook, GetBook, dll)
├── middleware/      # Middleware (opsional)
├── model/           # Struct model Book dan BookStore
├── router/          # Inisialisasi semua route dan middleware
├── utils/           # Utils untuk support kebutuhan lain-lain (opsional)
├── main.go          # Entry point
└── go.mod           # Modul Go
```

## 🚀 Cara Menjalankan

### 1. Clone repository

```bash
git clone https://github.com/bearmastery/test-sismedika.git
cd test-sismedika
```

### 2. Jalankan aplikasi

```bash
go run main.go
```

### 3. Tes endpoint dengan `curl` atau `Postman` atau `test.http`

Contoh:

```bash
curl -X POST -H "Content-Type: application/json" \
-d '{"title":"Buku Pertama", "author":"Riki", "published_year":2024}' \
http://localhost:8080/books
```

## 🧪 Menjalankan Unit Test

```bash
go test ./... -v
```

Check coverage:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 📦 Dependency

- [`go-chi/chi/v5`](https://github.com/go-chi/chi) – HTTP router
