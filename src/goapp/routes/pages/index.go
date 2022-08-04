package routes

import (
	session "main/pkg/session"
	template "main/pkg/template"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	gitHubUser, err := session.GetGitHubUserData(w, r)
	if err != nil {
		return
	}

	template.UseTemplate(&w, r, "index", gitHubUser)
}
