// handlers/post_handler.go
package handlers

import (
	"blog-api/database"
	"blog-api/models"
	"blog-api/utils"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// GET /posts
func GetPosts(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, title, content, created_at FROM posts")
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		utils.SendError(w, err.Error(), http.StatusInternalServerError, err)
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
		// http.Error(w, "ID tidak valid", http.StatusBadRequest)
		utils.SendError(w, "ID tidak valid", http.StatusBadRequest, err)
		return
	}

	var p models.Post
	err = database.DB.QueryRow("SELECT id, title, content, created_at FROM posts WHERE id = ?", id).
		Scan(&p.ID, &p.Title, &p.Content, &p.Created_at)
	if err != nil {
		// http.Error(w, "Post tidak ditemukan", http.StatusNotFound)
		utils.SendError(w, "Post tidak ditemukan", http.StatusNotFound, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

// POST /posts
func CreatePost(w http.ResponseWriter, r *http.Request) {
	var p models.Post
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		// http.Error(w, err.Error(), http.StatusBadRequest)
		utils.SendError(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	// validasi input
	if p.Title == "" {
		// http.Error(w, "Title tidak boleh kosong", http.StatusBadRequest)
		utils.SendError(w, "Title tidak boleh kosong", http.StatusBadRequest, nil)
		return
	}
	if p.Content == "" {
		// http.Error(w, "Content tidak boleh kosong", http.StatusBadRequest)
		utils.SendError(w, "Content tidak boleh kosong", http.StatusBadRequest, nil)
		return
	}

	res, err := database.DB.Exec("INSERT INTO posts (title, content) VALUES (?, ?)", p.Title, p.Content)
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		utils.SendError(w, err.Error(), http.StatusInternalServerError, err)
		return
	}

	id, _ := res.LastInsertId()
	// p.ID = int(id)
	// 1. Ambil data dari database menggunakan ID yang baru
	row := database.DB.QueryRow("SELECT id, title, content, created_at FROM posts WHERE id = ?", id)
	err = row.Scan(&p.ID, &p.Title, &p.Content, &p.Created_at)
	if err != nil {
		utils.SendError(w, "Gagal mengambil data post yang baru dibuat", http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

// PUT /posts/{id}
func UpdatePost(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// http.Error(w, "ID tidak valid", http.StatusBadRequest)
		utils.SendError(w, "ID tidak valid", http.StatusBadRequest, err)
		return
	}

	var p models.Post
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		// http.Error(w, err.Error(), http.StatusBadRequest)
		utils.SendError(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	// validasi input
	if p.Title == "" {
		// http.Error(w, "Title tidak boleh kosong", http.StatusBadRequest)
		utils.SendError(w, "Title tidak boleh kosong", http.StatusBadRequest, nil)
		return
	}
	if p.Content == "" {
		// http.Error(w, "Content tidak boleh kosong", http.StatusBadRequest)
		utils.SendError(w, "Content tidak boleh kosong", http.StatusBadRequest, nil)
		return
	}

	_, err = database.DB.Exec("UPDATE posts SET title = ?, content = ? WHERE id = ?", p.Title, p.Content, id)
	if err != nil {
		// http.Error(w, "Gagal update post", http.StatusInternalServerError)
		utils.SendError(w, "Gagal update post", http.StatusInternalServerError, err)
		return
	}

	// p.ID = int(id)
	// 1. Ambil data dari database menggunakan ID yang baru
	row := database.DB.QueryRow("SELECT id, title, content, created_at FROM posts WHERE id = ?", id)
	err = row.Scan(&p.ID, &p.Title, &p.Content, &p.Created_at)
	if err != nil {
		utils.SendError(w, "Gagal mengambil data post yang baru dibuat", http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

// DELETE /posts/{id}
func DeletePost(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// http.Error(w, "ID tidak valid", http.StatusBadRequest)
		utils.SendError(w, "ID tidak valid", http.StatusBadRequest, err)
		return
	}

	_, err = database.DB.Exec("DELETE FROM posts WHERE id = ?", id)
	if err != nil {
		// http.Error(w, "Gagal hapus post", http.StatusInternalServerError)
		utils.SendError(w, "Gagal hapus post", http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
