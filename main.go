package main

import (
	"fmt"
	"log"

	"github.com/bomboskuy/UAS-Backend/db"
	"github.com/bomboskuy/UAS-Backend/routes"
	"github.com/bomboskuy/UAS-Backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Failed load env")
	}

	db.ConnectPostgres()
	db.ConnectMongoDB()

	utils.InitLogger()
	fmt.Println("🚀 UAS Backend Started")

	app := fiber.New(fiber.Config{
		AppName: "UAS Backend",
	})

	app.Use(cors.New())

	routes.Register(app)

	log.Fatal(app.Listen(":8080"))
}