# 📝 Blog API Sederhana (Go + SQLite)

API untuk mengelola catatan/blog dengan operasi CRUD.

## 🚀 Jalankan

```bash
go run main.go

```

## 🔑 Autentikasi

Untuk aksi POST, PUT, DELETE, tambahkan header (postman) atau di .env:

```bash
X-API-Key: rahasia123
```

(Ganti rahasia123 dengan nilai di file .env Anda)

### 🔍 Opsional:

Tambahkan di bagian GET /posts:

?q=kata → cari di title/content
?page=2 → halaman ke-2 (default: 1)
?limit=10 → data per halaman (max: 100, default: 5)

### Keterangan

Contoh: 192.168.1.5 = localhost
http://192.168.1.5:8080/posts?q=abc
http://192.168.1.5:8080/posts?page=2&limit=3

Router: chi
Database: SQLite
Fitur: CRUD lengkap (POST, GET all, GET by ID, PUT, DELETE) + Filter & Paginate
Struktur proyek bersih dan mudah dikembangkan
