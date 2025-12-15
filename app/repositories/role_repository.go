package repositories

import "UAS-Backend/app/models"

type RoleRepository interface {
	FindByID(id string) (*models.Role, error)
	FindByName(name string) (*models.Role, error)
}
