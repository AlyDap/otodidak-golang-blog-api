module blog-api

go 1.25.4

require (
	github.com/glebarez/sqlite v1.11.0
	github.com/go-chi/chi/v5 v5.2.3
	github.com/joho/godotenv v1.5.1
)

require (
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/glebarez/go-sqlite v1.21.2 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	golang.org/x/sys v0.7.0 // indirect
	gorm.io/gorm v1.25.7 // indirect
	modernc.org/libc v1.22.5 // indirect
	modernc.org/mathutil v1.5.0 // indirect
	modernc.org/memory v1.5.0 // indirect
	modernc.org/sqlite v1.23.1 // indirect
)

// install dependensi
// go get github.com/go-chi/chi/v5
// go get github.com/mattn/go-sqlite3 (NO)
// go get github.com/glebarez/sqlite

//
// set CGO_ENABLED=1
// go env CGO_ENABLED
// go run main.go

// install godotenv untuk baca .env
// go get github.com/joho/godotenv

// go mod tidy 
