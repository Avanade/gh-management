package routes

import (
	"main/models"
	template "main/pkg/template"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request, data *models.TypPageData) {
	// var test map[string]string
	// test["data"] = "hi"
	// test["data2"] = "hello"
	// data.Content = test
	tmpl := template.UseTemplate("index")
	tmpl.Execute(w, data)
}
