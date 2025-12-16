package services

import (
	"database/sql"

	"github.com/bomboskuy/UAS-Backend/app/repositories"
	"github.com/bomboskuy/UAS-Backend/utils"
)

type AuthService struct {
	userRepo       repositories.UserRepository
	roleRepo       repositories.RoleRepository
	permissionRepo repositories.PermissionRepository
}

func NewAuthService(
	userRepo repositories.UserRepository,
	roleRepo repositories.RoleRepository,
	permissionRepo repositories.PermissionRepository,
) *AuthService {
	return &AuthService{
		userRepo:       userRepo,
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token        string      `json:"token"`
	RefreshToken string      `json:"refresh_token"`
	User         UserProfile `json:"user"`
}

type UserProfile struct {
	ID          string   `json:"id"`
	Username    string   `json:"username"`
	FullName    string   `json:"full_name"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
}

func (s *AuthService) Login(req LoginRequest) (*LoginResponse, error) {
	if req.Username == "" || req.Password == "" {
		return nil, sql.ErrNoRows
	}

	user, err := s.userRepo.FindByUsernameOrEmail(req.Username)
	if err != nil {
		return nil, err
	}

	if !user.IsActive {
		return nil, sql.ErrNoRows
	}

	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return nil, sql.ErrNoRows
	}

	role, err := s.roleRepo.FindByID(user.RoleID)
	if err != nil {
		return nil, err
	}

	permissions, _ := s.permissionRepo.FindByRoleID(user.RoleID)

	token, err := utils.GenerateToken(user.ID, user.RoleID, permissions)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User: UserProfile{
			ID:          user.ID,
			Username:    user.Username,
			FullName:    user.FullName,
			Role:        role.Name,
			Permissions: permissions,
		},
	}, nil
}

func (s *AuthService) GetProfile(userID, roleID string, permissions []string) (*UserProfile, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	role, err := s.roleRepo.FindByID(roleID)
	if err != nil {
		return nil, err
	}

	return &UserProfile{
		ID:          user.ID,
		Username:    user.Username,
		FullName:    user.FullName,
		Role:        role.Name,
		Permissions: permissions,
	}, nil
}