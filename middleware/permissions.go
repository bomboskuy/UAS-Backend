package middleware

import (
	"strings"

	"github.com/bomboskuy/UAS-Backend/helper"
	"github.com/bomboskuy/UAS-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return helper.Unauthorized(c, "Token tidak ditemukan")
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := utils.ValidateToken(tokenString)
		if err != nil || claims.UserID == "" {
			return helper.Unauthorized(c, "Token tidak valid atau expired")
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("role_id", claims.RoleID)
		c.Locals("permissions", claims.Permissions)

		return c.Next()
	}
}

func RequirePermission(permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		permissions, ok := c.Locals("permissions").([]string)
		if !ok {
			return helper.Forbidden(c, "Permission data not found")
		}

		hasPermission := false
		for _, p := range permissions {
			if p == permission {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			return helper.Forbidden(c, "Anda tidak memiliki akses untuk resource ini")
		}

		return c.Next()
	}
}