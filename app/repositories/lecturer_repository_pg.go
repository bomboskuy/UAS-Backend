package repositories

import (
	"database/sql"

	"github.com/bomboskuy/UAS-Backend/app/models"
)

type LecturerRepository interface {
	Create(lecturer *models.Lecturer) error
	FindAll() ([]models.Lecturer, error)
	FindByID(id string) (*models.Lecturer, error)
	FindByUserID(userID string) (*models.Lecturer, error)
	FindAdvisees(lecturerID string) ([]models.Student, error)
}

type LecturerRepositoryPg struct {
	db *sql.DB
}

func NewLecturerRepositoryPg(db *sql.DB) LecturerRepository {
	return &LecturerRepositoryPg{db: db}
}

func (r *LecturerRepositoryPg) Create(l *models.Lecturer) error {
	query := `
		INSERT INTO lecturers (
			id, user_id, lecturer_id, department, created_at
		) VALUES ($1,$2,$3,$4,$5)
	`
	_, err := r.db.Exec(
		query,
		l.ID,
		l.UserID,
		l.LecturerID,
		l.Department,
		l.CreatedAt,
	)
	return err
}

func (r *LecturerRepositoryPg) FindAll() ([]models.Lecturer, error) {
	query := `
		SELECT id, user_id, lecturer_id, department, created_at
		FROM lecturers
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lecturers []models.Lecturer
	for rows.Next() {
		var l models.Lecturer
		if err := rows.Scan(
			&l.ID,
			&l.UserID,
			&l.LecturerID,
			&l.Department,
			&l.CreatedAt,
		); err != nil {
			return nil, err
		}
		lecturers = append(lecturers, l)
	}

	return lecturers, nil
}

func (r *LecturerRepositoryPg) FindByID(id string) (*models.Lecturer, error) {
	query := `
		SELECT id, user_id, lecturer_id, department, created_at
		FROM lecturers
		WHERE id = $1
	`

	var l models.Lecturer
	err := r.db.QueryRow(query, id).Scan(
		&l.ID,
		&l.UserID,
		&l.LecturerID,
		&l.Department,
		&l.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &l, nil
}

func (r *LecturerRepositoryPg) FindByUserID(userID string) (*models.Lecturer, error) {
	query := `
		SELECT id, user_id, lecturer_id, department, created_at
		FROM lecturers
		WHERE user_id = $1
	`

	var l models.Lecturer
	err := r.db.QueryRow(query, userID).Scan(
		&l.ID,
		&l.UserID,
		&l.LecturerID,
		&l.Department,
		&l.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &l, nil
}

func (r *LecturerRepositoryPg) FindAdvisees(lecturerID string) ([]models.Student, error) {
	query := `
		SELECT id, user_id, student_id, program_study,
		       academic_year, advisor_id, created_at
		FROM students
		WHERE advisor_id = $1
	`

	rows, err := r.db.Query(query, lecturerID)
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
