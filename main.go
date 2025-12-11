// main.go
package main

import (
	"blog-api/database"
	"blog-api/handlers"
	mymiddleware "blog-api/middleware" // alias
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	// muat env jika ada
	godotenv.Load()

	database.InitDB()
	defer database.DB.Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	// r.Use(mymiddleware.Logger)

	// Endpoint publik
	r.Get("/posts", handlers.GetPosts)
	r.Get("/posts/{id}", handlers.GetPost)

	// Endpoint terlindungi
	// r.Route("/posts", func(r chi.Router) {
	// 	r.Use(mymiddleware.AuthMiddleware)
	// 	r.Post("/", handlers.CreatePost)
	// 	r.Put("/{id}", handlers.UpdatePost)
	// 	r.Delete("/{id}", handlers.DeletePost)
	// })

	// Endpoint terlindungi — hanya untuk write operations
	r.With(mymiddleware.AuthMiddleware).Post("/posts", handlers.CreatePost)
	r.With(mymiddleware.AuthMiddleware).Put("/posts/{id}", handlers.UpdatePost)
	r.With(mymiddleware.AuthMiddleware).Delete("/posts/{id}", handlers.DeletePost)

	// Layani file statis Halaman utama & aset
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./public/"))))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/index.html")
	})

	// log.Println("Server jalan di :8080")
	// http.ListenAndServe(":8080", r)

	// deploy ke local server
	log.Println("Server jalan di http://0.0.0.0:8080 (akses via IP lokal di jaringan)")
	http.ListenAndServe("0.0.0.0:8080", r)
	// 0.0.0.0 berarti: dengarkan semua antarmuka jaringan (bukan hanya localhost).
}
