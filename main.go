package main

import (
	"log"
	"os"
	"rosered/passbook-management/src/database"
	"rosered/passbook-management/src/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Setup database
	db := database.Connect()
	database.AutoMigrate(db)

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	app := fiber.New()

	routes.Setup(app, db)

	app.Listen(":" + port)
}
