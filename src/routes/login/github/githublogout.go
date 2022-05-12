package routes

import (
	session "main/pkg/session"
	"net/http"
)

func GitHubLogoutHandler(w http.ResponseWriter, r *http.Request) {
	err := session.RemoveGitHubAccount(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
