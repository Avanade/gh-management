package routes

import (
	"net/http"
	// "github.com/gorilla/mux"
	 template "main/pkg/template"
)

func CustomizeExternalLinks(w http.ResponseWriter, r *http.Request) {
	template.UseTemplate(&w, r, "admin/externallinks", nil)
}



func GetCustomizeExternalLinks(w http.ResponseWriter, r *http.Request) {
	// req := mux.Vars(r)

	// sessionaz, _ := session.Store.Get(r, "auth-session")
	// iprofile := sessionaz.Values["profile"]
	// profile := iprofile.(map[string]interface{})
	// username := profile["preferred_username"]

	// template.UseTemplate(&w, r, "externallinks/new", data)

}