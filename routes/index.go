package routes

import (
	"database/sql"

	"github.com/bomboskuy/UAS-Backend/app/repositories"
	"github.com/bomboskuy/UAS-Backend/app/services"
	"github.com/bomboskuy/UAS-Backend/db"
	"github.com/bomboskuy/UAS-Backend/helper"
	"github.com/bomboskuy/UAS-Backend/middleware"
	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App) {
	// Initialize repositories
	userRepo := repositories.NewUserRepositoryPg(db.DB)
	roleRepo := repositories.NewRoleRepositoryPg(db.DB)
	permissionRepo := repositories.NewPermissionRepositoryPg(db.DB)

	// Initialize services
	authService := services.NewAuthService(userRepo, roleRepo, permissionRepo)

	api := app.Group("/api/v1")

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/login", func(c *fiber.Ctx) error {
		var req services.LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return helper.BadRequest(c, "Invalid request body", nil)
		}

		result, err := authService.Login(req)
		if err == sql.ErrNoRows {
			return helper.Unauthorized(c, "Username atau password salah")
		}
		if err != nil {
			return helper.InternalServerError(c, "Login failed")
		}

		return helper.Success(c, "Login berhasil", result)
	})

	auth.Get("/profile", middleware.AuthRequired(), func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(string)
		roleID := c.Locals("role_id").(string)
		permissions := c.Locals("permissions").([]string)

		profile, err := authService.GetProfile(userID, roleID, permissions)
		if err != nil {
			return helper.NotFound(c, "User not found")
		}

		return helper.Success(c, "Profile retrieved", profile)
	})

	auth.Post("/logout", middleware.AuthRequired(), func(c *fiber.Ctx) error {
		return helper.Success(c, "Logout berhasil", nil)
	})
}
```

### Testing Postman FR-001
```
POST http://localhost:8080/api/v1/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "password123"
}

Response:
{
  "status": "success",
  "message": "Login berhasil",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": "uuid",
      "username": "admin",
      "full_name": "Admin User",
      "role": "Admin",
      "permissions": ["achievement:create", ...]
    }
  }
}