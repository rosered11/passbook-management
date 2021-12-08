package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"rosered/passbook-management/src/authentication"
	"rosered/passbook-management/src/dto"
	"rosered/passbook-management/src/utilities"

	"golang.org/x/oauth2"
)

type PassbookController interface {
	CreateGetPassbookHandler(w http.ResponseWriter, r *http.Request)
}

type DefaultPassbookController struct {
	context context.Context
	authen  authentication.Authentication
}

func NewPassbookController(context context.Context, authen authentication.Authentication) PassbookController {
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
