package routes

import (
	"main/models"
	session "main/pkg/session"
	template "main/pkg/template"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request, pageHeaders *models.TypHeaders) {
	var data models.TypPageData
	data.Header = pageHeaders

	gitHubUser, err := session.GetGitHubUserData(w, r)
	if err != nil {
		return
	}
	data.Content = gitHubUser

	template.UseTemplate(&w, r, &data, "index")
}
