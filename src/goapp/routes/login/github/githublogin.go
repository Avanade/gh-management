package routes

import (
	"log"
	"net/http"

	auth "main/pkg/authentication"
	"main/pkg/session"
)

func GithubLoginHandler(w http.ResponseWriter, r *http.Request) {
	state, err := session.GetState(w, r)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ghauth := auth.GetGitHubOauthConfig(r.Host)

	http.Redirect(w, r, ghauth.AuthCodeURL(state), http.StatusTemporaryRedirect)
}
