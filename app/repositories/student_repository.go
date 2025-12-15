package repositories

import "UAS-Backend/app/models"

type StudentRepository interface {
	Create(student *models.Student) error
	FindByUserID(userID string) (*models.Student, error)
	FindByAdvisorID(lecturerID string) ([]models.Student, error)
}
