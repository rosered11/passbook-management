package routes

import (
	"context"
	"rosered/passbook-management/src/authentication"
	"rosered/passbook-management/src/controllers"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
)


func Setup(app *fiber.App) {
	ctx := context.Background()

	// authentication
	authen := authentication.NewAuthentication(ctx)

	// controller
	authenController := controllers.NewAuthenController(ctx, authen)
	passbookController := controllers.NewPassbookController(ctx, authen)

	// setup path endpoint
	// prefix endpoing
	api := app.Group("api")

	openid := api.Group("openid")
	// controller path
	auth := openid.Group("auth")
	auth.Get("/callback", adaptor.HTTPHandlerFunc(authenController.CallbackHandler))

	passbook := openid.Group("passbook")
	passbook.Get("/", adaptor.HTTPHandlerFunc(passbookController.CreateGetPassbookHandler))
}
