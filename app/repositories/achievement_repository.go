package repositories

import "UAS-Backend/app/models"

type AchievementRepository interface {
	Create(achievement *models.Achievement) (string, error)
	FindByID(id string) (*models.Achievement, error)
	Update(id string, achievement *models.Achievement) error
	SoftDelete(id string) error
}
