// handlers/post_handler.go
package handlers

import (
	"blog-api/database"
	"blog-api/models"
	"blog-api/utils"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// GET /posts
func GetPosts(w http.ResponseWriter, r *http.Request) {
	// Ambil parameter query
	query := r.URL.Query().Get("q")        // kata kunci pencarian
	pageStr := r.URL.Query().Get("page")   // halaman ke berapa
	limitStr := r.URL.Query().Get("limit") // jumlah per halaman

	// Set default
	page := 1
	limit := 5

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	offset := (page - 1) * limit

	var rows *sql.Rows
	var count int
	var err error

	if query != "" {
		// Hitung total yang cocok
		countQuery := "SELECT COUNT(*) FROM posts WHERE title LIKE ? OR content LIKE ?"
		err = database.DB.QueryRow(countQuery, "%"+query+"%", "%"+query+"%").Scan(&count)
		if err != nil {
			utils.SendError(w, "Gagal menghitung data", http.StatusInternalServerError, err)
			return
		}

		// Ambil data dengan pencarian + pagination
		searchQuery := `
			SELECT id, title, content, created_at 
			FROM posts 
			WHERE title LIKE ? OR content LIKE ?
			ORDER BY created_at DESC
			LIMIT ? OFFSET ?
		`
		rows, err = database.DB.Query(searchQuery, "%"+query+"%", "%"+query+"%", limit, offset)
	} else {
		// Hitung total semua
		err = database.DB.QueryRow("SELECT COUNT(*) FROM posts").Scan(&count)
		if err != nil {
			utils.SendError(w, "Gagal menghitung total data", http.StatusInternalServerError, err)
			return
		}

		// Ambil semua data dengan pagination
		allQuery := `
			SELECT id, title, content, created_at 
			FROM posts 
			ORDER BY created_at DESC
			LIMIT ? OFFSET ?
		`
		rows, err = database.DB.Query(allQuery, limit, offset)
	}

	if err != nil {
		utils.SendError(w, "Gagal mengambil data", http.StatusInternalServerError, err)
		return
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		rows.Scan(&p.ID, &p.Title, &p.Content, &p.CreatedAt)
		posts = append(posts, p)
	}

	// Hitung total halaman
	totalPages := (count + limit - 1) / limit // ceil(count / limit)

	// Response dengan metadata
	response := struct {
		Data       []models.Post `json:"data"`
		Page       int           `json:"page"`
		Limit      int           `json:"limit"`
		Total      int           `json:"total"`
		TotalPages int           `json:"total_pages"`
		HasNext    bool          `json:"has_next"`
		HasPrev    bool          `json:"has_prev"`
	}{
		Data:       posts,
		Page:       page,
		Limit:      limit,
		Total:      count,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	// http://192.168.1.5:8080/posts?page=2&limit=3
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
		Scan(&p.ID, &p.Title, &p.Content, &p.CreatedAt)
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
	err = row.Scan(&p.ID, &p.Title, &p.Content, &p.CreatedAt)
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
	err = row.Scan(&p.ID, &p.Title, &p.Content, &p.CreatedAt)
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
