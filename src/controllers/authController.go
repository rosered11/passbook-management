package controllers

import (
	"github.com/gofiber/fiber/v2"
)

func Login(c *fiber.Ctx) error {

	return c.JSON(fiber.Map{
		"message": "hello",
	})
}

// func handleRedirect(w http.ResponseWriter, r *http.Request) {
// 	http.Redirect(w, r, oauth2Config.AuthCodeURL(state), http.StatusFound)
// }
