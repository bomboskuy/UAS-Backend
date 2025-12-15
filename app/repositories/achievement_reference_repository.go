package repositories

import "github.com/bomboskuy/UAS-Backend/app/models"

type AchievementReferenceRepository interface {
	Create(ref *models.AchievementReference) error
	FindByID(id string) (*models.AchievementReference, error)
	FindByStudentID(studentID string) ([]models.AchievementReference, error)
	UpdateStatus(id string, status string, verifierID *string, note *string) error
}
