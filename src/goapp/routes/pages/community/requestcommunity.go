package routes

import (
	session "main/pkg/session"
	template "main/pkg/template"
	"net/http"

	"github.com/gorilla/mux"
)

func CommunityHandler(w http.ResponseWriter, r *http.Request) {
	req := mux.Vars(r)
	id := req["id"]

	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := profile["preferred_username"]

	data := map[string]interface{}{
		"Id":    id,
		"Email": username,
	}

	template.UseTemplate(&w, r, "/community/requestcommunity", data)
}
