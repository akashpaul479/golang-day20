package project

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ExamResult struct {
	ID        int     `json:"id"`
	ExamID    int     `json:"exam_id"`
	StudentID int     `json:"student_id"`
	Marks     float64 `json:"marks"`
	Grade     string  `json:"grade"`
}

type examResultHandler struct {
	db5 *sql.DB
}

func NewExamResultHandler(db5 *sql.DB) *examResultHandler {
	return &examResultHandler{db5: db5}
}

// GET results for an exam
func (h *examResultHandler) GetExamResults(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	examID, _ := strconv.Atoi(params["id"])

	rows, err := h.db5.Query("SELECT id, exam_id, student_id, marks, grade FROM exam_results WHERE exam_id=?", examID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []ExamResult
	for rows.Next() {
		var er ExamResult
		if err := rows.Scan(&er.ID, &er.ExamID, &er.StudentID, &er.Marks, &er.Grade); err != nil {
			http.Error(w, "rows scan failed", http.StatusInternalServerError)
			return
		}
		results = append(results, er)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// INSERT results for an exam
func (h *examResultHandler) InsertExamResult(w http.ResponseWriter, r *http.Request) {
	var newResult ExamResult
	if err := json.NewDecoder(r.Body).Decode(&newResult); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res, err := h.db5.Exec("INSERT INTO exam_results(exam_id, student_id, marks, grade) VALUES (?, ?, ?, ?)",
		newResult.ExamID, newResult.StudentID, newResult.Marks, newResult.Grade)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, _ := res.LastInsertId()
	newResult.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newResult)
}
