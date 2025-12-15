package main

import (
	"fmt"
	"log"

	"github.com/bomboskuy/UAS-Backend/db"
	"github.com/bomboskuy/UAS-Backend/utils"
    "github.com/bomboskuy/UAS-Backend/routes"
	"github.com/joho/godotenv"
    "github.com/gofiber/fiber/v2"

)
func main() {
    if err := godotenv.Load(); err != nil {
		log.Fatal("Failed load env")
	}

	db.ConnectPostgres()

    utils.InitLogger()
    fmt.Println("Go Clean Architecture Started 🚀")

    app := fiber.New(fiber.Config{
		AppName: "UAS Backend",
	})

	// Register routes
	routes.Register(app)

	// Run server
	log.Fatal(app.Listen(":8080"))

}