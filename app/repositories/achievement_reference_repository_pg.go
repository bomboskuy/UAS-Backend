package repositories

import (
	"database/sql"
	"time"

	"github.com/bomboskuy/UAS-Backend/app/models"
)

type AchievementReferenceRepository interface {
	Create(ref *models.AchievementReference) error
	FindByID(id string) (*models.AchievementReference, error)
	FindByStudentID(studentID string) ([]models.AchievementReference, error)
	FindByAdvisorID(lecturerID string) ([]models.AchievementReference, error)
	FindVerifiedByStudentID(studentID string) ([]models.AchievementReference, error)
	UpdateStatus(id string, status string, verifierID *string, note *string) error
	FindAll() ([]models.AchievementReference, error)

	// 🔥 TAMBAHAN
	CountByStatus() (map[string]int, error)
}

type achievementReferenceRepositoryPg struct {
	db *sql.DB
}

func NewAchievementReferenceRepositoryPg(db *sql.DB) AchievementReferenceRepository {
	return &achievementReferenceRepositoryPg{db: db}
}

func (r *achievementReferenceRepositoryPg) Create(ref *models.AchievementReference) error {
	query := `
		INSERT INTO achievement_references 
		(id, student_id, mongo_achievement_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(
		query,
		ref.ID,
		ref.StudentID,
		ref.MongoAchievementID,
		ref.Status,
		ref.CreatedAt,
		ref.UpdatedAt,
	)
	return err
}

func (r *achievementReferenceRepositoryPg) FindByID(id string) (*models.AchievementReference, error) {
	var ref models.AchievementReference
	query := `
		SELECT id, student_id, mongo_achievement_id, status,
		       submitted_at, verified_at, verified_by, rejection_note,
		       created_at, updated_at
		FROM achievement_references WHERE id=$1
	`
	err := r.db.QueryRow(query, id).Scan(
		&ref.ID,
		&ref.StudentID,
		&ref.MongoAchievementID,
		&ref.Status,
		&ref.SubmittedAt,
		&ref.VerifiedAt,
		&ref.VerifiedBy,
		&ref.RejectionNote,
		&ref.CreatedAt,
		&ref.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &ref, nil
}

func (r *achievementReferenceRepositoryPg) FindByStudentID(studentID string) ([]models.AchievementReference, error) {
	query := `
		SELECT id, student_id, mongo_achievement_id, status,
		       submitted_at, verified_at, verified_by, rejection_note,
		       created_at, updated_at
		FROM achievement_references
		WHERE student_id=$1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var refs []models.AchievementReference
	for rows.Next() {
		var ref models.AchievementReference
		rows.Scan(
			&ref.ID,
			&ref.StudentID,
			&ref.MongoAchievementID,
			&ref.Status,
			&ref.SubmittedAt,
			&ref.VerifiedAt,
			&ref.VerifiedBy,
			&ref.RejectionNote,
			&ref.CreatedAt,
			&ref.UpdatedAt,
		)
		refs = append(refs, ref)
	}

	return refs, nil
}

func (r *achievementReferenceRepositoryPg) FindByAdvisorID(lecturerID string) ([]models.AchievementReference, error) {
	query := `
		SELECT ar.id, ar.student_id, ar.mongo_achievement_id, ar.status,
		       ar.submitted_at, ar.verified_at, ar.verified_by, ar.rejection_note,
		       ar.created_at, ar.updated_at
		FROM achievement_references ar
		INNER JOIN students s ON ar.student_id = s.id
		WHERE s.advisor_id=$1
		ORDER BY ar.created_at DESC
	`

	rows, err := r.db.Query(query, lecturerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var refs []models.AchievementReference
	for rows.Next() {
		var ref models.AchievementReference
		rows.Scan(
			&ref.ID,
			&ref.StudentID,
			&ref.MongoAchievementID,
			&ref.Status,
			&ref.SubmittedAt,
			&ref.VerifiedAt,
			&ref.VerifiedBy,
			&ref.RejectionNote,
			&ref.CreatedAt,
			&ref.UpdatedAt,
		)
		refs = append(refs, ref)
	}

	return refs, nil
}

func (r *achievementReferenceRepositoryPg) FindAll() ([]models.AchievementReference, error) {
	query := `
		SELECT id, student_id, mongo_achievement_id, status,
		       submitted_at, verified_at, verified_by, rejection_note,
		       created_at, updated_at
		FROM achievement_references
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var refs []models.AchievementReference
	for rows.Next() {
		var ref models.AchievementReference
		rows.Scan(
			&ref.ID,
			&ref.StudentID,
			&ref.MongoAchievementID,
			&ref.Status,
			&ref.SubmittedAt,
			&ref.VerifiedAt,
			&ref.VerifiedBy,
			&ref.RejectionNote,
			&ref.CreatedAt,
			&ref.UpdatedAt,
		)
		refs = append(refs, ref)
	}

	return refs, nil
}

func (r *achievementReferenceRepositoryPg) FindVerifiedByStudentID(studentID string) ([]models.AchievementReference, error) {
	query := `
		SELECT id, student_id, mongo_achievement_id, status,
		       submitted_at, verified_at, verified_by,
		       rejection_note, created_at, updated_at
		FROM achievement_references
		WHERE student_id=$1 AND status='verified'
		ORDER BY verified_at DESC
	`

	rows, err := r.db.Query(query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var refs []models.AchievementReference
	for rows.Next() {
		var ref models.AchievementReference
		rows.Scan(
			&ref.ID,
			&ref.StudentID,
			&ref.MongoAchievementID,
			&ref.Status,
			&ref.SubmittedAt,
			&ref.VerifiedAt,
			&ref.VerifiedBy,
			&ref.RejectionNote,
			&ref.CreatedAt,
			&ref.UpdatedAt,
		)
		refs = append(refs, ref)
	}

	return refs, nil
}

func (r *achievementReferenceRepositoryPg) UpdateStatus(
	id string,
	status string,
	verifierID *string,
	note *string,
) error {
	now := time.Now()

	switch status {
	case "submitted":
		_, err := r.db.Exec(
			`UPDATE achievement_references 
			 SET status=$1, submitted_at=$2, updated_at=$3 
			 WHERE id=$4`,
			status, now, now, id,
		)
		return err

	case "verified":
		_, err := r.db.Exec(
			`UPDATE achievement_references 
			 SET status=$1, verified_at=$2, verified_by=$3, updated_at=$4 
			 WHERE id=$5`,
			status, now, verifierID, now, id,
		)
		return err

	case "rejected":
		_, err := r.db.Exec(
			`UPDATE achievement_references 
			 SET status=$1, rejection_note=$2, updated_at=$3 
			 WHERE id=$4`,
			status, note, now, id,
		)
		return err
	}

	return nil
}

func (r *achievementReferenceRepositoryPg) CountByStatus() (map[string]int, error) {
	query := `
		SELECT status, COUNT(*) 
		FROM achievement_references
		GROUP BY status
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := map[string]int{
		"draft":     0,
		"submitted": 0,
		"verified":  0,
		"rejected":  0,
	}

	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		result[status] = count
	}

	return result, nil
}
