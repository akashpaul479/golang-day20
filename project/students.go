package project

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Student struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type studentHandler struct {
	db3 *sql.DB
}

func NewStudentHandler(db3 *sql.DB) *studentHandler {
	return &studentHandler{
		db3: db3,
	}
}

func (s *studentHandler) GetStudents(w http.ResponseWriter, r *http.Request) {
	rows, err := s.db3.Query("SELECT id , name , email FROM students")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var students []Student
	for rows.Next() {
		var s Student
		if err := rows.Scan(&s.ID, &s.Name, &s.Email); err != nil {
			http.Error(w, "rows scan failed", http.StatusInternalServerError)
			return
		}
		students = append(students, s)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}
func (s *studentHandler) InsertStudents(w http.ResponseWriter, r *http.Request) {
	var Newstudent Student
	if err := json.NewDecoder(r.Body).Decode(&Newstudent); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res, err := s.db3.Exec("INSERT INTO students(name , email)VALUES (? , ?)", Newstudent.Name, Newstudent.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, _ := res.LastInsertId()
	Newstudent.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Newstudent)
}
func (s *studentHandler) UpdateStudents(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var updated Student
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err := s.db3.Exec("UPDATE students  SET name=?, email=? WHERE id=? ", updated.Name, updated.Email, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	updated.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}
func (s *studentHandler) DeleteStudents(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	_, err := s.db3.Exec("DELETE FROM students WHERE id=?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (s *studentHandler) getstudentsbyId(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var student Student
	err = s.db3.QueryRow("Select id , name , email FROM students WHERE id=?", id).
		Scan(&student.ID, &student.Name, &student.Email)

	if err == sql.ErrNoRows {
		http.Error(w, "student not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)

}
