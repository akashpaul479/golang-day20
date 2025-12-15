package project

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Lecturer struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Department string `json:"department"`
}

type lecturerHandler struct {
	db1 *sql.DB
}

func NewLecturerHandler(db1 *sql.DB) *lecturerHandler {
	return &lecturerHandler{
		db1: db1,
	}
}

// GET all lecturers
func (h *lecturerHandler) GetLecturers(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db1.Query("SELECT id, name, email, department FROM lecturer")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var lecturers []Lecturer
	for rows.Next() {
		var l Lecturer
		if err := rows.Scan(&l.ID, &l.Name, &l.Email, &l.Department); err != nil {
			http.Error(w, "rows scan failed", http.StatusInternalServerError)
			return
		}
		lecturers = append(lecturers, l)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lecturers)
}

// INSERT new lecturer
func (h *lecturerHandler) InsertLecturer(w http.ResponseWriter, r *http.Request) {
	var newLecturer Lecturer
	if err := json.NewDecoder(r.Body).Decode(&newLecturer); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res, err := h.db1.Exec("INSERT INTO lecturer(name, email, department) VALUES (?, ?, ?)",
		newLecturer.Name, newLecturer.Email, newLecturer.Department)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, _ := res.LastInsertId()
	newLecturer.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newLecturer)
}

// UPDATE lecturer by ID
func (h *lecturerHandler) UpdateLecturer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var updated Lecturer
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err := h.db1.Exec("UPDATE lecturer SET name=?, email=?, department=? WHERE id=?",
		updated.Name, updated.Email, updated.Department, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	updated.ID = id

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

// DELETE lecturer by ID
func (h *lecturerHandler) DeleteLecturer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	_, err := h.db1.Exec("DELETE FROM lecturer WHERE id=?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
