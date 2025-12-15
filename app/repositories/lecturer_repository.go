package repositories

import "github.com/bomboskuy/UAS-Backend/app/models"

type LecturerRepository interface {
	FindByUserID(userID string) (*models.Lecturer, error)
	FindAdvisees(lecturerID string) ([]models.Student, error)
}
