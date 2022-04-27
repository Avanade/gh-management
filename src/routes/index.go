package routes

import (
	"main/models"
	template "main/pkg/template"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request, pageHeaders *models.TypHeaders) {
	var data models.TypPageData
	data.Header = pageHeaders

	template.UseTemplate(&w, r, &data, "index")
}
