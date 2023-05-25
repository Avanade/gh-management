package routes

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/coreos/go-oidc"

	auth "main/pkg/authentication"
	ghmgmt "main/pkg/ghmgmtdb"
	"main/pkg/msgraph"
	session "main/pkg/session"
)

func CallbackHandler(w http.ResponseWriter, r *http.Request) {

	// Check session
	session, err := session.Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.URL.Query().Get("state") != session.Values["state"] {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	//Retrieve token
	authenticator, err := auth.NewAuthenticator(r.Host)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := authenticator.Config.Exchange(context.TODO(), r.URL.Query().Get("code"))
	if err != nil {
		log.Printf("no token found: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "No id_token field in oauth2 token.", http.StatusInternalServerError)
		return
	}

	oidcConfig := &oidc.Config{
		ClientID: os.Getenv("CLIENT_ID"),
	}

	idToken, err := authenticator.Provider.Verifier(oidcConfig).Verify(context.TODO(), rawIDToken)

	if err != nil {
		http.Error(w, "Failed to verify ID Token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the userInfo
	var profile map[string]interface{}
	if err := idToken.Claims(&profile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["userHasPhoto"] = hasPhoto
	session.Values["userPhoto"] = userPhoto
	errS := session.Save(r, w)

	// Insert Azure User
	name := fmt.Sprint(profile["name"])
	errIU := ghmgmt.InsertUser(userPrincipalName, name, "", "", "")
	if errIU != nil {
		http.Error(w, errIU.Error(), http.StatusInternalServerError)
		return
	}

	if errS != nil {
		http.Error(w, errS.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to index
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
