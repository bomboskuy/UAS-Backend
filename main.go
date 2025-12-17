package main

import (
	"fmt"
	"log"

	"github.com/bomboskuy/UAS-Backend/app/repositories"
	"github.com/bomboskuy/UAS-Backend/app/services"
	"github.com/bomboskuy/UAS-Backend/db"
	"github.com/bomboskuy/UAS-Backend/routes"
	"github.com/bomboskuy/UAS-Backend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	// Load ENV
	if err := godotenv.Load(); err != nil {
		log.Fatal("Failed load env")
	}

	// Connect DB
	db.ConnectPostgres()
	db.ConnectMongoDB()

	utils.InitLogger()

	fmt.Println("🚀 UAS Backend Started")

	// Fiber App
	app := fiber.New(fiber.Config{
		AppName: "UAS Backend",
	})

	app.Use(cors.New())

	// =====================
	// REPOSITORIES
	// =====================
	userRepo := repositories.NewUserRepositoryPg(db.DB)
	roleRepo := repositories.NewRoleRepositoryPg(db.DB)
	permissionRepo := repositories.NewPermissionRepositoryPg(db.DB)

	studentRepo := repositories.NewStudentRepositoryPg(db.DB)
	lecturerRepo := repositories.NewLecturerRepositoryPg(db.DB)

	achievementRepo := repositories.NewAchievementRepositoryMongo(db.MongoDB)
	achievementRefRepo := repositories.NewAchievementReferenceRepositoryPg(db.DB)

	// =====================
	// SERVICES
	// =====================
	authService := services.NewAuthService(
		userRepo,
		roleRepo,
		permissionRepo,
	)

	userService := services.NewUserService(
		userRepo,
		roleRepo,
		studentRepo,
		lecturerRepo,
	)

	achievementService := services.NewAchievementService(
		achievementRepo,
		achievementRefRepo,
		studentRepo,
		lecturerRepo,
	)

	// =====================
	// ROUTES (NO LOGIC)
	// =====================
	routes.Register(
		app,
		authService,
		userService,
		achievementService,
	)

	log.Fatal(app.Listen(":3000"))
}
