package repositories

import (
	"database/sql"

	"github.com/bomboskuy/UAS-Backend/app/models"
)

type lecturerRepositoryPg struct {
	db *sql.DB
}

func NewLecturerRepositoryPg(db *sql.DB) LecturerRepository {
	return &lecturerRepositoryPg{db: db}
}

func (r *lecturerRepositoryPg) FindByUserID(userID string) (*models.Lecturer, error) {
	var lecturer models.Lecturer
	query := `SELECT id, user_id, lecturer_id, department, created_at FROM lecturers WHERE user_id=$1`
	err := r.db.QueryRow(query, userID).Scan(&lecturer.ID, &lecturer.UserID, &lecturer.LecturerID, &lecturer.Department, &lecturer.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &lecturer, nil
}

func (r *lecturerRepositoryPg) FindAdvisees(lecturerID string) ([]models.Student, error) {
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