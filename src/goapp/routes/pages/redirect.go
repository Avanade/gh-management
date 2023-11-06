package routes

import (
	"html/template"
	"net/http"
)

func LoginRedirectHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	redirect := "/"
	if len(q["redirect"]) > 0 {
		redirect = q["redirect"][0]
	}
	if len(q["search"]) > 0 {
		redirect = redirect + "&search=" + q["search"][0]
	}
	data := map[string]interface{}{
		"redirect": redirect,
	}
	tmpl := template.Must(template.ParseFiles("templates/loginredirect.html"))
	tmpl.Execute(w, data)
}

func GHAuthenticatedHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/ghauthenticated.html"))
	tmpl.Execute(w, nil)
}

func GHAuthenticatingHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/ghauthenticating.html"))
	tmpl.Execute(w, nil)
}

func GitRedirectHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	redirect := "/"
	if len(q["redirect"]) > 0 {
		redirect = q["redirect"][0]
	}
	if len(q["search"]) > 0 {
		redirect = redirect + "&search=" + q["search"][0]
	}
	data := map[string]interface{}{
		"redirect": redirect,
	}
	tmpl := template.Must(template.ParseFiles("templates/gitredirect.html"))
	tmpl.Execute(w, data)
}
