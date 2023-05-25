package routes

import (
	auth "main/pkg/authentication"
	session "main/pkg/session"
	"net/http"
)

func GithubLoginHandler(w http.ResponseWriter, r *http.Request) {
	state, err := session.GetState(w, r)
	if err != nil {
		return
	}

	ghauth := auth.GetGitHubOauthConfig(r.Host)

	http.Redirect(w, r, ghauth.AuthCodeURL(state), http.StatusTemporaryRedirect)
}
