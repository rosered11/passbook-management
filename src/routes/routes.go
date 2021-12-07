package routes

import (
	"rosered/passbook-management/src/controllers"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("api")

	auth := api.Group("auth")
	auth.Post("/login", controllers.Login)
}
