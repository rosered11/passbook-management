package routes

import (
	"context"
	"rosered/passbook-management/src/authentication"
	"rosered/passbook-management/src/controllers"
	"rosered/passbook-management/src/database"
	"rosered/passbook-management/src/services"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Setup(app *fiber.App, db *gorm.DB) {
	ctx := context.Background()

	// authentication
	authen := authentication.NewAuthentication(ctx)

	// repository
	repository := database.NewRepository()

	// unitofwork + repositories
	unitofwork := database.NewUnitOfWork(db, repository)

	// services
	passbookService := services.NewPassbookService(unitofwork)

	// controller
	authenController := controllers.NewAuthenController(ctx, authen)
	passbookController := controllers.NewPassbookController(ctx, authen, passbookService)

	// setup path endpoint
	// prefix endpoing
	api := app.Group("api")

	openid := api.Group("openid")
	// controller openid path
	authOpenid := openid.Group("auth")
	authOpenid.Get("/callback", adaptor.HTTPHandlerFunc(authenController.CallbackHandler))

	passbookOpenid := openid.Group("passbook")
	passbookOpenid.Get("/", adaptor.HTTPHandlerFunc(passbookController.CreateGetPassbookHandler))

	// controller path
	passbook := api.Group("passbook")
	passbook.Post("/", passbookController.CreatePassbook)
	passbook.Get("/", passbookController.GetPassbooks)
	passbook.Get("/owner/:owner", passbookController.GetPassbooksWithOwner)
	passbook.Get("/transaction", passbookController.GetTransactions)
	passbook.Get("/transaction/owner/:owner", passbookController.GetTransactionsWithOwner)
	passbook.Get("/transaction/date/:date", passbookController.GetTransactionsWithDate)
	passbook.Get("/transaction/owner/:owner/date/:date", passbookController.GetTransactionsWithDateAndOwner)
}
