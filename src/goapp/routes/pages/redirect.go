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

func AuthenticationSuccessfulHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/authenticationsuccessful.html"))
	tmpl.Execute(w, nil)
}

func AuthenticationInProgressHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/authenticationinprogress.html"))
	tmpl.Execute(w, nil)
}

func AuthenticationFailedHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/authenticationfailed.html"))
	tmpl.Execute(w, nil)
}

func GHAuthenticationInProgressHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/ghauthenticationinprogress.html"))
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
