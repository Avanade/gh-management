package routes

import (
	"log"
	"net/http"

	"main/pkg/session"
)

func GitHubLogoutHandler(w http.ResponseWriter, r *http.Request) {
	err := session.RemoveGitHubAccount(w, r)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
