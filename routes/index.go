package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/bomboskuy/UAS-Backend/app/services"
	"github.com/bomboskuy/UAS-Backend/middleware"
)

func Register(
	app *fiber.App,
	authService *services.AuthService,
	userService *services.UserService,
	achievementService *services.AchievementService,
) {
	api := app.Group("/api/v1")

	// =====================
	// AUTH
	// =====================
	auth := api.Group("/auth")
	auth.Post("/login", authService.Login)
	auth.Get("/profile", middleware.AuthRequired(), authService.Profile)
	auth.Post("/logout", middleware.AuthRequired(), authService.Logout)

	// =====================
	// USERS (ADMIN)
	// =====================
	users := api.Group("/users",
		middleware.AuthRequired(),
		middleware.RequirePermission("user:manage"),
	)

	users.Get("/", userService.GetAll)
	users.Get("/:id", userService.GetByID)
	users.Post("/", userService.Create)
	users.Put("/:id", userService.Update)
	users.Delete("/:id", userService.Delete)

	// =====================
	// ACHIEVEMENTS
	// =====================
	achievements := api.Group("/achievements",
		middleware.AuthRequired(),
	)

	achievements.Post("/",
		middleware.RequirePermission("achievement:create"),
		achievementService.Create,
	)

	achievements.Get("/",
		middleware.RequirePermission("achievement:read"),
		achievementService.GetAll,
	)

	achievements.Post("/:id/submit",
		middleware.RequirePermission("achievement:update"),
		achievementService.Submit,
	)

	achievements.Post("/:id/verify",
		middleware.RequirePermission("achievement:verify"),
		achievementService.Verify,
	)

	achievements.Post("/:id/reject",
		middleware.RequirePermission("achievement:verify"),
		achievementService.Reject,
	)

	// =====================
	// STUDENTS
	// =====================
	students := api.Group("/students",
		middleware.AuthRequired(),
		middleware.RequirePermission("user:manage"),
	)

	students.Get("/", userService.GetStudents)
	students.Get("/:id", userService.GetStudentByID)
	students.Put("/:id/advisor", userService.AssignAdvisor)

	// =====================
	// LECTURERS
	// =====================
	lecturers := api.Group("/lecturers",
		middleware.AuthRequired(),
	)

	lecturers.Get("/",
		middleware.RequirePermission("user:manage"),
		userService.GetLecturers,
	)

	lecturers.Get("/:id/advisees",
		middleware.RequirePermission("achievement:verify"),
		userService.GetAdvisees,
	)

	// =====================
	// REPORTS
	// =====================
	reports := api.Group("/reports",
		middleware.AuthRequired(),
	)

	reports.Get("/statistics",
		middleware.RequirePermission("achievement:read"),
		achievementService.Statistics,
	)

	reports.Get("/student/:id",
		middleware.RequirePermission("achievement:read"),
		achievementService.StudentReport,
	)
}
