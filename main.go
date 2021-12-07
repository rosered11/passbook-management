package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type Response struct {
	Name string `json:"name"`
	Path string
}

func main() {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, "https://identity-server-4213.herokuapp.com")
	if err != nil {
		log.Fatal(err)
	}
	config := oauth2.Config{
		ClientID:     "defaultClientOpenId",
		ClientSecret: "secret",
		Endpoint:     provider.Endpoint(),
		RedirectURL:  "http://localhost:3000/a",
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("access_token") == "" {
			state, err := randString(16)
			if err != nil {
				http.Error(w, "Internal error", http.StatusInternalServerError)
				return
			}
			setCallbackCookie(w, r, "state", state)
			setCallbackCookie(w, r, "oring-path", r.URL.Path)

			http.Redirect(w, r, config.AuthCodeURL(state), http.StatusFound)
		}

		var accessToken *oauth2.Token
		accessToken = &oauth2.Token{AccessToken: r.URL.Query().Get("access_token")}

		userInfo, err := provider.UserInfo(ctx, oauth2.StaticTokenSource(accessToken))
		if err != nil {
			http.Error(w, "Failed to get userinfo: "+err.Error(), http.StatusInternalServerError)
			return
		}

		var response Response
		userInfo.Claims(&response)
		response.Path = r.URL.Path
		data, err := json.MarshalIndent(response, "", "    ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(data)
	})

	http.HandleFunc("/b", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("access_token") == "" {
			state, err := randString(16)
			if err != nil {
				http.Error(w, "Internal error", http.StatusInternalServerError)
				return
			}
			setCallbackCookie(w, r, "state", state)
			setCallbackCookie(w, r, "oring-path", r.URL.Path)

			http.Redirect(w, r, config.AuthCodeURL(state), http.StatusFound)
		}

		var accessToken *oauth2.Token
		accessToken = &oauth2.Token{AccessToken: r.URL.Query().Get("access_token")}

		userInfo, err := provider.UserInfo(ctx, oauth2.StaticTokenSource(accessToken))
		if err != nil {
			http.Error(w, "Failed to get userinfo: "+err.Error(), http.StatusInternalServerError)
			return
		}

		var response Response
		userInfo.Claims(&response)
		response.Path = r.URL.Path

		data, err := json.MarshalIndent(response, "", "    ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(data)
	})

	http.HandleFunc("/a", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("rawPath: " + r.URL.RawPath)
		log.Printf("path: " + r.URL.Path)
		originPath, err := r.Cookie("oring-path")
		if err != nil {
			http.Error(w, "state not found", http.StatusBadRequest)
			return
		}
		log.Printf("origin-path: " + originPath.Value)
		state, err := r.Cookie("state")
		if err != nil {
			http.Error(w, "state not found", http.StatusBadRequest)
			return
		}

		log.Printf("state: " + state.Value)
		log.Printf("code: " + r.URL.Query().Get("code"))

		if r.URL.Query().Get("state") != state.Value {
			http.Error(w, "state did not match", http.StatusBadRequest)
			return
		}

		oauth2Token, err := config.Exchange(ctx, r.URL.Query().Get("code"))
		if err != nil {
			http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, originPath.Value+"?access_token="+oauth2Token.AccessToken, http.StatusFound)

		userInfo, err := provider.UserInfo(ctx, oauth2.StaticTokenSource(oauth2Token))
		if err != nil {
			http.Error(w, "Failed to get userinfo: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// resp := struct {
		// 	OAuth2Token *oauth2.Token
		// 	UserInfo    *oidc.UserInfo
		// }{oauth2Token, userInfo}

		type claimes struct {
			Subject string `json:"sub"`
			Profile string `json:"profile"`
			Name    string `json:"name"`
		}
		var claims claimes
		userInfo.Claims(&claims)
		data, err := json.MarshalIndent(claims, "", "    ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(data)
	})

	log.Printf("listening on http://%s/", "127.0.0.1:3000")
	log.Fatal(http.ListenAndServe("127.0.0.1:3000", nil))
}

func randString(nByte int) (string, error) {
	b := make([]byte, nByte)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func setCallbackCookie(w http.ResponseWriter, r *http.Request, name, value string) {
	c := &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   int(time.Hour.Seconds()),
		Secure:   r.TLS != nil,
		HttpOnly: true,
	}
	http.SetCookie(w, c)
}
