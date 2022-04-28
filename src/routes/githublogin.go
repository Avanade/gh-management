package routes

import (
	"fmt"
	session "main/pkg/session"
	"net/http"
	"os"
)

func GithubLoginHandler(w http.ResponseWriter, r *http.Request) {
	githubClientID := os.Getenv("ghclientid")
	githubCallbackUrl := os.Getenv("ghcallbackurl")
	state, err := session.GetState(w, r)
	if err != nil {
		return
	}

	redirectURL := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&state=%s", githubClientID, githubCallbackUrl, state)

	http.Redirect(w, r, redirectURL, 301)
}
