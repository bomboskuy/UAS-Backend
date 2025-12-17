package repositories

import (
	"database/sql"

	"github.com/bomboskuy/UAS-Backend/app/models"
)

type StudentRepository interface {
	Create(student *models.Student) error
	FindAll() ([]models.Student, error)
	FindByID(id string) (*models.Student, error)
	FindByUserID(userID string) (*models.Student, error)
	AssignAdvisor(studentID string, advisorID string) error
}

type StudentRepositoryPg struct {
	db *sql.DB
}

func NewStudentRepositoryPg(db *sql.DB) StudentRepository {
	return &StudentRepositoryPg{db: db}
}

func (r *StudentRepositoryPg) Create(student *models.Student) error {
	query := `
		INSERT INTO students (
			id, user_id, student_id, program_study,
			academic_year, advisor_id, created_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7)
	`
	_, err := r.db.Exec(
		query,
		student.ID,
		student.UserID,
		student.StudentID,
		student.ProgramStudy,
		student.AcademicYear,
		student.AdvisorID,
		student.CreatedAt,
	)
	return err
}

func (r *StudentRepositoryPg) FindAll() ([]models.Student, error) {
	query := `
		SELECT id, user_id, student_id, program_study,
		       academic_year, advisor_id, created_at
		FROM students
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []models.Student
	for rows.Next() {
		var s models.Student
		if err := rows.Scan(
			&s.ID,
			&s.UserID,
			&s.StudentID,
			&s.ProgramStudy,
			&s.AcademicYear,
			&s.AdvisorID,
			&s.CreatedAt,
		); err != nil {
			return nil, err
		}
		students = append(students, s)
	}

	return students, nil
}

func (r *StudentRepositoryPg) FindByID(id string) (*models.Student, error) {
	query := `
		SELECT id, user_id, student_id, program_study,
		       academic_year, advisor_id, created_at
		FROM students
		WHERE id = $1
	`

	var s models.Student
	err := r.db.QueryRow(query, id).Scan(
		&s.ID,
		&s.UserID,
		&s.StudentID,
		&s.ProgramStudy,
		&s.AcademicYear,
		&s.AdvisorID,
		&s.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (r *StudentRepositoryPg) FindByUserID(userID string) (*models.Student, error) {
	query := `
		SELECT id, user_id, student_id, program_study,
		       academic_year, advisor_id, created_at
		FROM students
		WHERE user_id = $1
	`

	var s models.Student
	err := r.db.QueryRow(query, userID).Scan(
		&s.ID,
		&s.UserID,
		&s.StudentID,
		&s.ProgramStudy,
		&s.AcademicYear,
		&s.AdvisorID,
		&s.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (r *StudentRepositoryPg) AssignAdvisor(studentID string, advisorID string) error {
	query := `
		UPDATE students
		SET advisor_id = $1
		WHERE id = $2
	`
	_, err := r.db.Exec(query, advisorID, studentID)
	return err
}
