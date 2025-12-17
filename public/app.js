// Konfigurasi
// const API_BASE = "http://localhost:8080";
// const API_BASE = "http://192.168.1.5:8080";
// const API_BASE = window.location.origin;
const API_BASE = "";
const API_KEY = "rahasia123"; // sesuaikan dengan .env

// State
let currentPage = 1;
let currentQuery = "";

// Elemen DOM
const postList = document.getElementById("postList");
const addForm = document.getElementById("addForm");
const searchInput = document.getElementById("search");
const searchBtn = document.getElementById("searchBtn");
const prevBtn = document.getElementById("prev");
const nextBtn = document.getElementById("next");
const pageInfo = document.getElementById("pageInfo");
const loadingEl = document.getElementById("loading");
const errorEl = document.getElementById("error");

// Helper: tampilkan error
function showError(msg) {
  errorEl.textContent = msg;
  errorEl.style.display = "block";
  setTimeout(() => (errorEl.style.display = "none"), 5000);
}

// Helper: ambil data posts
async function fetchPosts(page = 1, query = "") {
  loadingEl.style.display = "block";
  try {
    let url = `${API_BASE}/posts?page=${page}&limit=5`;
    if (query) url += `&q=${encodeURIComponent(query)}`;

    const res = await fetch(url);
    if (!res.ok) throw new Error("Gagal ambil data");
    const data = await res.json();

    // Pastikan data.data selalu array
    if (!Array.isArray(data.data)) {
      throw new Error("Format respons tidak valid");
    }

    // Render
    postList.innerHTML = "";
    if (data.data.length === 0) {
      postList.innerHTML = "<li>Tidak ada post.</li>";
    } else {
      data.data.forEach((post) => {
        const li = document.createElement("li");
        li.dataset.id = post.id;
        li.innerHTML = `
          <div class="post-view" id="view-${post.id}">
            <strong>${post.title}</strong><br>
            <small>${new Date(post.created_at).toLocaleString()}</small><br>
            ${post.content}
            <br>
            <button class="edit-btn" data-id="${post.id}">✏️ Edit</button>
            <button class="delete-btn" data-id="${post.id}">🗑️ Hapus</button>
          </div>
          <div class="post-edit" id="edit-${post.id}" style="display:none;">
            <input type="text" class="edit-title" value="${
              post.title
            }" required />
            <textarea class="edit-content" required>${post.content}</textarea>
            <button class="save-btn" data-id="${post.id}">💾 Simpan</button>
            <button class="cancel-btn" data-id="${post.id}">❌ Batal</button>
          </div>
        `;
        postList.appendChild(li);
      });

      // Event hapus
      document.querySelectorAll(".delete-btn").forEach((btn) => {
        btn.addEventListener("click", (e) => {
          const id = e.target.dataset.id;
          deletePost(id);
        });
      });

      // Event edit
      document.querySelectorAll(".edit-btn").forEach((btn) => {
        btn.addEventListener("click", (e) => {
          const id = e.target.dataset.id;
          document.getElementById(`view-${id}`).style.display = "none";
          document.getElementById(`edit-${id}`).style.display = "block";
        });
      });

      // Event cancel
      document.querySelectorAll(".cancel-btn").forEach((btn) => {
        btn.addEventListener("click", (e) => {
          const id = e.target.dataset.id;
          document.getElementById(`view-${id}`).style.display = "block";
          document.getElementById(`edit-${id}`).style.display = "none";
        });
      });

      // Event save
      document.querySelectorAll(".save-btn").forEach((btn) => {
        btn.addEventListener("click", (e) => {
          const id = e.target.dataset.id;
          const title = document
            .querySelector(`#edit-${id} .edit-title`)
            .value.trim();
          const content = document
            .querySelector(`#edit-${id} .edit-content`)
            .value.trim();

          if (!title || !content) {
            showError("Judul dan isi tidak boleh kosong");
            return;
          }

          updatePost(id, title, content);
        });
      });
    }

    // Update pagination
    currentPage = page;
    currentQuery = query;
    pageInfo.textContent = `Halaman ${page}`;
    prevBtn.disabled = page <= 1;
    nextBtn.disabled = !data.has_next;
  } catch (err) {
    showError(err.message);
  } finally {
    loadingEl.style.display = "none";
  }
}

// Tambah post
addForm.addEventListener("submit", async (e) => {
  e.preventDefault();
  const title = document.getElementById("title").value;
  const content = document.getElementById("content").value;
  const addBtn = document.getElementById("addBtn");
  addBtn.disabled = true;
  addBtn.textContent = "Menyimpan...";

  try {
    const res = await fetch(`${API_BASE}/posts`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-API-Key": API_KEY,
      },
      body: JSON.stringify({ title, content }),
    });
    if (!res.ok) throw new Error("Gagal simpan post");
    addForm.reset();
    showNotif("Post berhasil ditambahkan!");
    fetchPosts(1, currentQuery); // refresh
  } catch (err) {
    showError(err.message);
  } finally {
    addBtn.disabled = false;
    addBtn.textContent = "➕ Tambah Post";
  }
});

// Search
searchBtn.addEventListener("click", () => {
  const q = searchInput.value.trim();
  fetchPosts(1, q);
});

searchInput.addEventListener("keyup", (e) => {
  if (e.key === "Enter") searchBtn.click();
});

// Pagination
prevBtn.addEventListener("click", () => {
  if (currentPage > 1) fetchPosts(currentPage - 1, currentQuery);
});
nextBtn.addEventListener("click", () => {
  fetchPosts(currentPage + 1, currentQuery);
});

// Update post
async function updatePost(id, title, content) {
  const saveBtn = document.querySelector(`#edit-${id} .save-btn`);
  saveBtn.disabled = true;
  saveBtn.textContent = "Menyimpan...";

  try {
    const res = await fetch(`${API_BASE}/posts/${id}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        "X-API-Key": API_KEY,
      },
      body: JSON.stringify({ title, content }),
    });
    if (!res.ok) throw new Error("Gagal update post");

    // Kembali ke tampilan
    document.getElementById(`view-${id}`).style.display = "block";
    document.getElementById(`edit-${id}`).style.display = "none";

    showNotif("Post berhasil diupdate!");
    // Refresh daftar (opsional: update hanya post ini)
    fetchPosts(currentPage, currentQuery);
  } catch (err) {
    showError(err.message);
  } finally {
    saveBtn.disabled = false;
    saveBtn.textContent = "💾 Simpan";
  }
}

// Hapus post
async function deletePost(id) {
  if (!confirm("Hapus post ini?")) return;
  const deleteBtn = document.querySelector(`.delete-btn[data-id="${id}"]`);
  deleteBtn.disabled = true;
  deleteBtn.textContent = "Menghapus...";

  try {
    const res = await fetch(`${API_BASE}/posts/${id}`, {
      method: "DELETE",
      headers: { "X-API-Key": API_KEY },
    });
    if (!res.ok) throw new Error("Gagal hapus post");
    showNotif("Post berhasil dihapus!", "success");
    fetchPosts(currentPage, currentQuery); // refresh
  } catch (err) {
    showError(err.message);
  } finally {
    deleteBtn.disabled = false;
    deleteBtn.textContent = "🗑️ Hapus";
  }
}

// show notif
function showNotif(message, type = "success") {
  const notif = document.getElementById("notif");
  notif.textContent = message;
  notif.style.display = "block";
  notif.style.backgroundColor = type === "success" ? "#d4edda" : "#f8d7da";
  notif.style.color = type === "success" ? "#155724" : "#721c24";
  setTimeout(() => (notif.style.display = "none"), 3000);
}

// Muat awal
fetchPosts();
