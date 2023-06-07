package routes

import (
	"net/http"

	"main/pkg/session"
	"main/pkg/template"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	gitHubUser, err := session.GetGitHubUserData(w, r)
	if err != nil {
		return
	}

	template.UseTemplate(&w, r, "index", gitHubUser)
}
