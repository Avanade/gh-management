package routes

import (
	"net/http"
	 template "main/pkg/template"
)

func CustomizeExternalLinks(w http.ResponseWriter, r *http.Request) {
	template.UseTemplate(&w, r, "admin/externallinks", nil)
}



// func AddCustomizeExternalLinks(w http.ResponseWriter, r *http.Request) {

// 	sessionaz, _ := session.Store.Get(r, "auth-session")
// 	iprofile := sessionaz.Values["profile"]
// 	profile := iprofile.(map[string]interface{})
// 	username := profile["preferred_username"]

// 	template.UseTemplate(&w, r, "externallinks/new", data)

// }