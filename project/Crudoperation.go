package project

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func CrudOperation() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("err loading .env file: %v", err)
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	fmt.Println("dsn:", dsn)
	var err error
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("db ping error:%v", err)
	}
	fmt.Println("connection succesfull!")

	stdhandler := NewStudentHandler(db)
	r := mux.NewRouter()
	r.HandleFunc("/students", stdhandler.GetStudents).Methods("GET")
	r.HandleFunc("/students", stdhandler.InsertStudents).Methods("POST")
	r.HandleFunc("/students/{id}", stdhandler.UpdateStudents).Methods("PUT")
	r.HandleFunc("/students/{id}", stdhandler.DeleteStudents).Methods("DELETE")
	r.HandleFunc("/students/{id}", stdhandler.getstudentsbyId).Methods("GET")

	lechandler := NewLecturerHandler(db)
	r.HandleFunc("/lecturer", lechandler.GetLecturers).Methods("GET")
	r.HandleFunc("/lecturer", lechandler.InsertLecturer).Methods("POST")
	r.HandleFunc("/lecturer/{id}", lechandler.UpdateLecturer).Methods("PUT")
	r.HandleFunc("/lecturer/{id}", lechandler.DeleteLecturer).Methods("DELETE")

	officehandler := NewOfficeHandler(db)
	r.HandleFunc("/officestaff", officehandler.GetOfficeStaff).Methods("GET")
	r.HandleFunc("/officestaff", officehandler.InsertOfficeStaff).Methods("POST")
	r.HandleFunc("/officestaff/{id}", officehandler.UpdateOfficeStaff).Methods("PUT")
	r.HandleFunc("/officestaff/{id}", officehandler.DeleteOfficeStaff).Methods("DELETE")

	examHandler := NewExamHandler(db)
	r.HandleFunc("/exams", examHandler.GetExams).Methods("GET")
	r.HandleFunc("/exams", examHandler.InsertExam).Methods("POST")
	r.HandleFunc("/exams/{id}", examHandler.UpdateExam).Methods("PUT")
	r.HandleFunc("/exams/{id}", examHandler.DeleteExam).Methods("DELETE")

	examresultHandler := NewExamResultHandler(db)
	r.HandleFunc("/exams/{id}/results", examresultHandler.GetExamResults).Methods("GET")
	r.HandleFunc("/exams/{id}/results", examresultHandler.InsertExamResult).Methods("POST")

	log.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
