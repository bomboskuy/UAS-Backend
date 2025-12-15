package repositories

import "UAS-Backend/app/models"

type LecturerRepository interface {
	FindByUserID(userID string) (*models.Lecturer, error)
	FindAdvisees(lecturerID string) ([]models.Student, error)
}
