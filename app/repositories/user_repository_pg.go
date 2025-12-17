package repositories

import (
	"github.com/bomboskuy/UAS-Backend/app/models"
	"database/sql"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByID(id string) (*models.User, error)
	FindByUsernameOrEmail(value string) (*models.User, error)
	FindAll() ([]models.User, error)
	Update(user *models.User) error
	Delete(id string) error
}

type userRepositoryPg struct {
	db *sql.DB
}

func NewUserRepositoryPg(db *sql.DB) UserRepository {
	return &userRepositoryPg{db: db}
}

func (r *userRepositoryPg) Create(user *models.User) error {
	query := `
		INSERT INTO users (id, username, email, password_hash, full_name, role_id, is_active)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
	`
	_, err := r.db.Exec(
		query,
		user.ID,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.FullName,
		user.RoleID,
		user.IsActive,
	)
	return err
}

func (r *userRepositoryPg) FindByID(id string) (*models.User, error) {
	var user models.User

	query := `
		SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at
		FROM users WHERE id=$1
	`
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.RoleID,
		&user.IsActive,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryPg) FindByUsernameOrEmail(value string) (*models.User, error) {
	var user models.User

	query := `
		SELECT id, username, email, password_hash, full_name, role_id, is_active
		FROM users
		WHERE username=$1 OR email=$1
	`
	err := r.db.QueryRow(query, value).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.RoleID,
		&user.IsActive,
	)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryPg) FindAll() ([]models.User, error) {
	rows, err := r.db.Query(`SELECT id, username, email, full_name, role_id FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		rows.Scan(&u.ID, &u.Username, &u.Email, &u.FullName, &u.RoleID)
		users = append(users, u)
	}
	return users, nil
}

func (r *userRepositoryPg) Update(user *models.User) error {
	query := `
		UPDATE users SET username=$2, email=$3, full_name=$4, role_id=$5, is_active=$6
		WHERE id=$1
	`
	_, err := r.db.Exec(
		query,
		user.ID,
		user.Username,
		user.Email,
		user.FullName,
		user.RoleID,
		user.IsActive,
	)
	return err
}

func (r *userRepositoryPg) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM users WHERE id=$1`, id)
	return err
}

