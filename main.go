// main.go
package main

import (
	"blog-api/database"
	"blog-api/handlers"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	database.InitDB()
	defer database.DB.Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/posts", func(r chi.Router) {
		r.Get("/", handlers.GetPosts)
		r.Post("/", handlers.CreatePost)
		r.Get("/{id}", handlers.GetPost)
		r.Put("/{id}", handlers.UpdatePost)
		r.Delete("/{id}", handlers.DeletePost)
	})

	// log.Println("Server jalan di :8080")
	// http.ListenAndServe(":8080", r)

	// deploy ke local server
	log.Println("Server jalan di http://0.0.0.0:8080 (akses via IP lokal di jaringan)")
	http.ListenAndServe("0.0.0.0:8080", r)
	// 0.0.0.0 berarti: dengarkan semua antarmuka jaringan (bukan hanya localhost).
}
