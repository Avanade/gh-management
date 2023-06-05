package routes

import (
	"log"
	session "main/pkg/session"
	"net/http"
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
