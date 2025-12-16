package repositories

import (
	"database/sql"

	"github.com/bomboskuy/UAS-Backend/app/models"
)

type studentRepositoryPg struct {
	db *sql.DB
}

func NewStudentRepositoryPg(db *sql.DB) StudentRepository {
	return &studentRepositoryPg{db: db}
}

func (r *studentRepositoryPg) Create(student *models.Student) error {
	query := `
		INSERT INTO students (id, user_id, student_id, program_study, academic_year, advisor_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(query, student.ID, student.UserID, student.StudentID, student.ProgramStudy, student.AcademicYear, student.AdvisorID, student.CreatedAt)
	return err
}

func (r *studentRepositoryPg) FindByUserID(userID string) (*models.Student, error) {
	var student models.Student
	query := `SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at FROM students WHERE user_id=$1`
	err := r.db.QueryRow(query, userID).Scan(&student.ID, &student.UserID, &student.StudentID, &student.ProgramStudy, &student.AcademicYear, &student.AdvisorID, &student.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &student, nil
}

func (r *studentRepositoryPg) FindByAdvisorID(lecturerID string) ([]models.Student, error) {
	query := `SELECT id, user_id, student_id, program_study, academic_year, advisor_id, created_at FROM students WHERE advisor_id=$1`
	rows, err := r.db.Query(query, lecturerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	students := []models.Student{}
	for rows.Next() {
		var s models.Student
		rows.Scan(&s.ID, &s.UserID, &s.StudentID, &s.ProgramStudy, &s.AcademicYear, &s.AdvisorID, &s.CreatedAt)
		students = append(students, s)
	}

	return students, nil
}