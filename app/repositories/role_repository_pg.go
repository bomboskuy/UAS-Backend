package repositories

import (
	"database/sql"

	"github.com/bomboskuy/UAS-Backend/app/models"
)

type roleRepositoryPg struct {
	db *sql.DB
}

func NewRoleRepositoryPg(db *sql.DB) RoleRepository {
	return &roleRepositoryPg{db: db}
}

func (r *roleRepositoryPg) FindByID(id string) (*models.Role, error) {
	var role models.Role
	query := `SELECT id, name, description, created_at FROM roles WHERE id=$1`
	err := r.db.QueryRow(query, id).Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepositoryPg) FindByName(name string) (*models.Role, error) {
	var role models.Role
	query := `SELECT id, name, description, created_at FROM roles WHERE name=$1`
	err := r.db.QueryRow(query, name).Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &role, nil
}