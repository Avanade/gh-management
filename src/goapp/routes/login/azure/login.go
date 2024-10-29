package routes

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"

	"main/pkg/appinsights_wrapper"
	auth "main/pkg/authentication"
	"main/pkg/session"

	"github.com/microsoft/ApplicationInsights-Go/appinsights/contracts"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	// Generate random state
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.LogTrace("Generated random state", contracts.Information)
	state := base64.StdEncoding.EncodeToString(b)
	session, err := session.Store.Get(r, "auth-session")
	if err != nil {
		logger.LogException(err)
		http.Redirect(w, r, "/clearcookies", http.StatusTemporaryRedirect)
		return
	}
	logger.LogTrace("Storing state in session", contracts.Information)
	session.Values["state"] = state
	err = session.Save(r, w)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.LogTrace("Session saved", contracts.Information)
	authenticator, err := auth.NewAuthenticator(r.Host)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.LogTrace(fmt.Sprint("Redirecting to Azure login. State: "+state), contracts.Information)
	http.Redirect(w, r, authenticator.Config.AuthCodeURL(state), http.StatusTemporaryRedirect)
}

func ClearCookies(w http.ResponseWriter, r *http.Request) {
	c := http.Cookie{
		Name:   "auth-session",
		MaxAge: -1}
	http.SetCookie(w, &c)
	cgh := http.Cookie{
		Name:   "gh-auth-session",
		MaxAge: -1}
	http.SetCookie(w, &cgh)

	http.Redirect(w, r, "/login/azure", http.StatusTemporaryRedirect)
}
