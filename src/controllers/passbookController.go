package controllers

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"rosered/passbook-management/src/authentication"
	"rosered/passbook-management/src/dto"
	"rosered/passbook-management/src/services"
	"rosered/passbook-management/src/utilities"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/lestrrat-go/jwx/jwk"
	"golang.org/x/oauth2"
)

type PassbookController interface {
	CreateGetPassbookHandler(w http.ResponseWriter, r *http.Request)
	CreatePassbook(c *fiber.Ctx) error
}

type DefaultPassbookController struct {
	context         context.Context
	authen          authentication.Authentication
	passbookService services.PassbookService
}

func NewPassbookController(context context.Context, authen authentication.Authentication, passbookService services.PassbookService) PassbookController {
	return DefaultPassbookController{context: context, authen: authen}
}

func (passbookController DefaultPassbookController) CreateGetPassbookHandler(w http.ResponseWriter, r *http.Request) {
	provider, config, _ := passbookController.authen.GetProvider()

	if r.URL.Query().Get("access_token") == "" {
		state, err := utilities.RandString(16)
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
		utilities.SetCallbackCookie(w, r, "state", state)
		utilities.SetCallbackCookie(w, r, "oring-path", r.URL.Path)

		http.Redirect(w, r, config.AuthCodeURL(state), http.StatusFound)
		return
	}

	var accessToken *oauth2.Token
	accessToken = &oauth2.Token{AccessToken: r.URL.Query().Get("access_token")}

	userInfo, err := provider.UserInfo(passbookController.context, oauth2.StaticTokenSource(accessToken))
	if err != nil {
		http.Error(w, "Failed to get userinfo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var response dto.AuthenResponse
	userInfo.Claims(&response)

	data, err := json.MarshalIndent(response, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// fiber unuse this function
func (passbookController DefaultPassbookController) CreatePassbook(c *fiber.Ctx) error {
	err := validateToken(c)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{"message": err.Error()})
	}

	c.Status(fiber.StatusNoContent)
	return nil
}

func validateToken(c *fiber.Ctx) error {
	headertoken := string(c.Request().Header.Peek("Authorization"))

	// Split bearer from token
	splitToken := strings.Split(headertoken, " ")
	if len(splitToken) == 1 {
		return errors.New("Invalid token format.")
	}

	token, err := jwt.Parse(splitToken[1], func(token *jwt.Token) (interface{}, error) {
		set, err := jwk.Fetch(c.Context(), os.Getenv("AuthenUrl")+os.Getenv("AuthenJwksPathUrl"))
		if err != nil {
			return nil, err
		}

		keyId, ok := token.Header["kid"].(string)

		if !ok {
			return nil, errors.New("expect jwt kid")
		}

		key, _ := set.LookupKeyID(keyId)

		var rawkey interface{} // This is the raw key, like *rsa.PrivateKey or *ecdsa.PrivateKey
		if err := key.Raw(&rawkey); err != nil {
			return nil, err
		}

		// We know this is an RSA Key so...
		rsa, ok := rawkey.(*rsa.PublicKey)
		if !ok {
			return nil, errors.New(fmt.Sprintf("expected ras key, got %T", rawkey))
		}
		return rsa, nil
	})

	if err != nil || !token.Valid {
		return err
	}

	// claims := token.Claims.(jwt.MapClaims)
	// for key, value := range claims {
	// 	fmt.Printf("%s\t%v\n", key, value)
	// }

	return nil
}
