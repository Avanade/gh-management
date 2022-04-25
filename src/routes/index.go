package routes

import (
	"net/http"
	session "webserver/pkg/session"
	template "webserver/pkg/template"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {

	session, err := session.Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := session.Values["profile"]

	tmpl := template.UseTemplate("index")
	tmpl.Execute(w, user)
}
