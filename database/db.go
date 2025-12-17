// database/db.go
package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	// _ "github.com/glebarez/sqlite"
)

var DB *sql.DB

func InitDB() {
	// Baca koneksi dari env (lebih aman)
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Fallback untuk development
		dsn = "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"
	}

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Gagal buka koneksi PostgreSQL:", err)
	}

	// Cek koneksi
	if err = DB.Ping(); err != nil {
		log.Fatal("Gagal ping database:", err)
	}

	// Buat tabel jika belum ada
	createTable := `
	CREATE TABLE IF NOT EXISTS posts (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);`

	_, err = DB.Exec(createTable)
	if err != nil {
		log.Fatal("Gagal buat tabel posts:", err)
	}

	log.Println("Koneksi ke PostgreSQL berhasil")
}

func InitDB_sqlite() {
	var err error
	DB, err = sql.Open("sqlite", "./blog.db")
	if err != nil {
		log.Fatal("Gagal buka database:", err)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = DB.Exec(createTable)
	if err != nil {
		log.Fatal("Gagal buat tabel:", err)
	}
}
