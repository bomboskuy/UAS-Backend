package repositories

import (
	"database/sql"
)

type PermissionRepository interface {
	FindByRoleID(roleID string) ([]string, error)
}

type permissionRepositoryPg struct {
	db *sql.DB
}

func NewPermissionRepositoryPg(db *sql.DB) PermissionRepository {
	return &permissionRepositoryPg{db: db}
}

func (r *permissionRepositoryPg) FindByRoleID(roleID string) ([]string, error) {
	query := `
		SELECT p.name FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id=$1
	`
	rows, err := r.db.Query(query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	permissions := []string{}
	for rows.Next() {
		var perm string
		rows.Scan(&perm)
		permissions = append(permissions, perm)
	}

	return permissions, nil
}