-- Exams table
USE college_mm_system;
CREATE TABLE IF NOT EXISTS exams (
  id INT AUTO_INCREMENT PRIMARY KEY,
  course VARCHAR(100) NOT NULL,
  exam_date DATE NOT NULL,
  mode ENUM('offline','online') NOT NULL,
  published BOOLEAN DEFAULT FALSE
);

-- Exam Results table
USE college_mm_system;
CREATE TABLE IF NOT EXISTS exam_results (
  id INT AUTO_INCREMENT PRIMARY KEY,
  exam_id INT NOT NULL,
  student_id INT NOT NULL,
  marks DECIMAL(5,2) NOT NULL,
  grade ENUM('A','B','C','D','E','F') NOT NULL,
  FOREIGN KEY (exam_id) REFERENCES exams(id) ON DELETE CASCADE,
  FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE
);