package project

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Exam struct {
	ID        int    `json:"id"`
	Course    string `json:"course"`
	ExamDate  string `json:"exam_date"`
	Mode      string `json:"mode"`
	Published bool   `json:"published"`
}

type examHandler struct {
	db4 *sql.DB
}

func NewExamHandler(db4 *sql.DB) *examHandler {
	return &examHandler{db4: db4}
}

// GET all exams
func (h *examHandler) GetExams(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db4.Query("SELECT id, course, exam_date, mode, published FROM exams")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var exams []Exam
	for rows.Next() {
		var e Exam
		if err := rows.Scan(&e.ID, &e.Course, &e.ExamDate, &e.Mode, &e.Published); err != nil {
			http.Error(w, "rows scan failed", http.StatusInternalServerError)
			return
		}
		exams = append(exams, e)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exams)
}

// INSERT new exam
func (h *examHandler) InsertExam(w http.ResponseWriter, r *http.Request) {
	var newExam Exam
	if err := json.NewDecoder(r.Body).Decode(&newExam); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res, err := h.db4.Exec("INSERT INTO exams(course, exam_date, mode, published) VALUES (?, ?, ?, ?)",
		newExam.Course, newExam.ExamDate, newExam.Mode, newExam.Published)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, _ := res.LastInsertId()
	newExam.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newExam)
}

// UPDATE exam by ID
func (h *examHandler) UpdateExam(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var updated Exam
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err := h.db4.Exec("UPDATE exams SET course=?, exam_date=?, mode=?, published=? WHERE id=?",
		updated.Course, updated.ExamDate, updated.Mode, updated.Published, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	updated.ID = id

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

// DELETE exam by ID
func (h *examHandler) DeleteExam(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	_, err := h.db4.Exec("DELETE FROM exams WHERE id=?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
