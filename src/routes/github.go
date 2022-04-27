package routes

import (
	"main/models"
	template "main/pkg/template"
	"net/http"
)

func GithubHandler(w http.ResponseWriter, r *http.Request, pageHeaders *models.TypHeaders) {
	var data models.TypPageData
	data.Header = pageHeaders
	template.UseTemplate(&w, r, &data, "github")
}
