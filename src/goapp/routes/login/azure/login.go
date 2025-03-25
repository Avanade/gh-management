package routes

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	auth "main/pkg/authentication"
	"main/pkg/session"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Generate random state
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	state := base64.StdEncoding.EncodeToString(b)
	session, err := session.Store.Get(r, "auth-session")
	if err != nil {
		http.Redirect(w, r, "/clearcookies", http.StatusTemporaryRedirect)
		return
	}
	session.Values["state"] = state
	session.Options.MaxAge = 2592000
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	authenticator, err := auth.NewAuthenticator(r.Host)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
