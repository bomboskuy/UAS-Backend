package services

import (
	"github.com/gofiber/fiber/v2"

	"github.com/bomboskuy/UAS-Backend/app/models"
	"github.com/bomboskuy/UAS-Backend/app/repositories"
	"github.com/bomboskuy/UAS-Backend/helper"
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

//
// ==========================
// LOGIN
// ==========================
//
func (s *AuthService) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest(c, "Invalid request body", nil)
	}

	if req.Username == "" || req.Password == "" {
		return helper.Unauthorized(c, "Username atau password salah")
	}

	user, err := s.userRepo.FindByUsernameOrEmail(req.Username)
	if err != nil || !user.IsActive {
		return helper.Unauthorized(c, "Username atau password salah")
	}

	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return helper.Unauthorized(c, "Username atau password salah")
	}

	role, err := s.roleRepo.FindByID(user.RoleID)
	if err != nil {
		return helper.InternalServerError(c, "Failed to get role")
	}

	permissions, _ := s.permissionRepo.FindByRoleID(user.RoleID)

	token, err := utils.GenerateToken(user.ID, user.RoleID, permissions)
	if err != nil {
		return helper.InternalServerError(c, "Failed to generate token")
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return helper.InternalServerError(c, "Failed to generate refresh token")
	}

	return helper.Success(c, "Login berhasil", models.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User: models.UserProfile{
			ID:          user.ID,
			Username:    user.Username,
			FullName:    user.FullName,
			Role:        role.Name,
			Permissions: permissions,
		},
	})
}

//
// ==========================
// PROFILE
// ==========================
//
func (s *AuthService) Profile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	roleID := c.Locals("role_id").(string)
	permissions := c.Locals("permissions").([]string)

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return helper.NotFound(c, "User not found")
	}

	role, err := s.roleRepo.FindByID(roleID)
	if err != nil {
		return helper.InternalServerError(c, "Role not found")
	}

	return helper.Success(c, "Profile retrieved", models.UserProfile{
		ID:          user.ID,
		Username:    user.Username,
		FullName:    user.FullName,
		Role:        role.Name,
		Permissions: permissions,
	})
}

//
// ==========================
// LOGOUT
// ==========================
//
func (s *AuthService) Logout(c *fiber.Ctx) error {
	// Stateless JWT → logout cukup di client
	return helper.Success(c, "Logout berhasil", nil)
}
