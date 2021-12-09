package controllers

import (
	"context"
	"net/http"
	"rosered/passbook-management/src/authentication"
)

type AuthenController interface{}

type DefaultAuthenController struct {
	context context.Context
	authen  authentication.Authentication
}

func NewAuthenController(context context.Context, authen authentication.Authentication) DefaultAuthenController {
	return DefaultAuthenController{context: context, authen: authen}
}

// http
func (authenController DefaultAuthenController) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	_, config, _ := authenController.authen.GetProvider()

	originPath, err := r.Cookie("oring-path")
	if err != nil {
		http.Error(w, "state not found", http.StatusBadRequest)
		return
	}

	state, err := r.Cookie("state")
	if err != nil {
		http.Error(w, "state not found", http.StatusBadRequest)
		return
	}

	if r.URL.Query().Get("state") != state.Value {
		http.Error(w, "state did not match", http.StatusBadRequest)
		return
	}

	oauth2Token, err := config.Exchange(authenController.context, r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, originPath.Value+"?access_token="+oauth2Token.AccessToken, http.StatusFound)
}
