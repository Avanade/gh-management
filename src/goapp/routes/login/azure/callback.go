package routes

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/coreos/go-oidc"
	"github.com/gorilla/sessions"

	auth "main/pkg/authentication"
	db "main/pkg/ghmgmtdb"
	"main/pkg/msgraph"
	"main/pkg/session"
)

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	// Check session
	session, err := session.Store.Get(r, "auth-session")
	if err != nil {
		http.Redirect(w, r, "/authentication/azure/failed", http.StatusSeeOther)
		return
	}

	if r.URL.Query().Get("state") != session.Values["state"] {
		http.Redirect(w, r, "/authentication/azure/failed", http.StatusSeeOther)
		return
	}

	//Retrieve token
	authenticator, err := auth.NewAuthenticator(r.Host)
	if err != nil {
		http.Redirect(w, r, "/authentication/azure/failed", http.StatusSeeOther)
		return
	}

	token, err := authenticator.Config.Exchange(context.TODO(), r.URL.Query().Get("code"))
	if err != nil {
		http.Redirect(w, r, "/authentication/azure/failed", http.StatusSeeOther)
		return
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		http.Redirect(w, r, "/authentication/azure/failed", http.StatusSeeOther)
		return
	}

	oidcConfig := &oidc.Config{
		ClientID: os.Getenv("CLIENT_ID"),
	}

	idToken, err := authenticator.Provider.Verifier(oidcConfig).Verify(context.TODO(), rawIDToken)

	if err != nil {
		http.Redirect(w, r, "/authentication/azure/failed", http.StatusSeeOther)
		return
	}

	// Get the userInfo
	var profile map[string]interface{}
	if err := idToken.Claims(&profile); err != nil {
		http.Redirect(w, r, "/authentication/azure/failed", http.StatusSeeOther)
		return
	}

	userPrincipalName := fmt.Sprint(profile["preferred_username"])

	session.Values["id_token"] = rawIDToken
	session.Values["access_token"] = token.AccessToken
	session.Values["profile"] = profile
	session.Values["refresh_token"] = token.RefreshToken
	session.Values["expiry"] = token.Expiry.UTC().Format("2006-01-02 15:04:05")
	isAdmin, _ := msgraph.IsUserAdmin(fmt.Sprintf("%s", profile["oid"]))
	session.Values["isUserAdmin"] = isAdmin
	hasPhoto, userPhoto, err := msgraph.GetUserPhoto(fmt.Sprintf("%s", profile["oid"]))
	if err != nil {
		http.Redirect(w, r, "/authentication/azure/failed", http.StatusSeeOther)
		return
	}
	ghUser := db.Users_Get_GHUser(userPrincipalName)
	session.Values["isGHAssociated"] = ghUser != ""
	session.Values["userHasPhoto"] = hasPhoto
	session.Values["userPhoto"] = userPhoto
	session.Options = &sessions.Options{
		Path:     "/",
		Domain:   "",
		MaxAge:   2592000,
		Secure:   true,
		HttpOnly: false,
		SameSite: http.SameSiteNoneMode,
	}
	err = session.Save(r, w)
	if err != nil {
		http.Redirect(w, r, "/authentication/azure/failed", http.StatusSeeOther)
		return
	}

	// Insert Azure User
	name := fmt.Sprint(profile["name"])
	err = db.InsertUser(userPrincipalName, name, "", "", "")
	if err != nil {
		http.Redirect(w, r, "/authentication/azure/failed", http.StatusSeeOther)
		return
	}

	// Redirect to index
	http.Redirect(w, r, "/authentication/azure/successful", http.StatusSeeOther)
}
