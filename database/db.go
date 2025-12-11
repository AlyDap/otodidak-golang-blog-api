// database/db.go
package database

import (
	"database/sql"
	"log"

	_ "github.com/glebarez/sqlite"
)

var DB *sql.DB

func InitDB() {
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
