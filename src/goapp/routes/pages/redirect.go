package routes

import (
	"net/http"
	"html/template"
)

func LoginRedirectHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	redirect := "/"
	if len(q["redirect"]) > 0 {
		redirect = q["redirect"][0]
	}
	data := map[string]interface{}{
		"redirect": redirect,
	}
	tmpl := template.Must(template.ParseFiles("templates/loginredirect.html"))
	tmpl.Execute(w, data)
}