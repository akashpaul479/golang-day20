package project

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type OfficeStaff struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type officeHandler struct {
	db2 *sql.DB
}

func NewOfficeHandler(db2 *sql.DB) *officeHandler {
	return &officeHandler{
		db2: db2,
	}
}

// GET all office staff
func (h *officeHandler) GetOfficeStaff(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db2.Query("SELECT id, name, email, role FROM officestaff")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var staff []OfficeStaff
	for rows.Next() {
		var o OfficeStaff
		if err := rows.Scan(&o.ID, &o.Name, &o.Email, &o.Role); err != nil {
			http.Error(w, "rows scan failed", http.StatusInternalServerError)
			return
		}
		staff = append(staff, o)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(staff)
}

// INSERT new office staff
func (h *officeHandler) InsertOfficeStaff(w http.ResponseWriter, r *http.Request) {
	var newStaff OfficeStaff
	if err := json.NewDecoder(r.Body).Decode(&newStaff); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res, err := h.db2.Exec("INSERT INTO officestaff(name, email, role) VALUES (?, ?, ?)",
		newStaff.Name, newStaff.Email, newStaff.Role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, _ := res.LastInsertId()
	newStaff.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newStaff)
}

// UPDATE office staff by ID
func (h *officeHandler) UpdateOfficeStaff(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var updated OfficeStaff
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err := h.db2.Exec("UPDATE officestaff SET name=?, email=?, role=? WHERE id=?",
		updated.Name, updated.Email, updated.Role, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	updated.ID = id

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

// DELETE office staff by ID
func (h *officeHandler) DeleteOfficeStaff(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	_, err := h.db2.Exec("DELETE FROM officestaff WHERE id=?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
