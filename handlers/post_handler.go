// handlers/post_handler.go
package handlers

import (
	"blog-api/database"
	"blog-api/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// GET /posts
func GetPosts(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, title, content, created_at FROM posts")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		rows.Scan(&p.ID, &p.Title, &p.Content, &p.Created_at)
		posts = append(posts, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

// GET /posts/{id}
func GetPost(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	var p models.Post
	err = database.DB.QueryRow("SELECT id, title, content, created_at FROM posts WHERE id = ?", id).
		Scan(&p.ID, &p.Title, &p.Content, &p.Created_at)
	if err != nil {
		http.Error(w, "Post tidak ditemukan", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

// POST /posts
func CreatePost(w http.ResponseWriter, r *http.Request) {
	var p models.Post
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validasi input
	if p.Title == "" {
		http.Error(w, "Title tidak boleh kosong", http.StatusBadRequest)
		return
	}
	if p.Content == "" {
		http.Error(w, "Content tidak boleh kosong", http.StatusBadRequest)
		return
	}

	res, err := database.DB.Exec("INSERT INTO posts (title, content) VALUES (?, ?)", p.Title, p.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := res.LastInsertId()
	p.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

// PUT /posts/{id}
func UpdatePost(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	var p models.Post
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validasi input
	if p.Title == "" {
		http.Error(w, "Title tidak boleh kosong", http.StatusBadRequest)
		return
	}
	if p.Content == "" {
		http.Error(w, "Content tidak boleh kosong", http.StatusBadRequest)
		return
	}

	_, err = database.DB.Exec("UPDATE posts SET title = ?, content = ? WHERE id = ?", p.Title, p.Content, id)
	if err != nil {
		http.Error(w, "Gagal update post", http.StatusInternalServerError)
		return
	}

	p.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

// DELETE /posts/{id}
func DeletePost(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	_, err = database.DB.Exec("DELETE FROM posts WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Gagal hapus post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
