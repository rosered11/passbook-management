package main

import (
	"log"
	"os"
	"rosered/passbook-management/src/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	app := fiber.New()

	routes.Setup(app)

	app.Listen(":" + port)
}
