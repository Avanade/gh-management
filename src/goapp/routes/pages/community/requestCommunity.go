package routes

import (
	"net/http"

	"main/pkg/session"
	"main/pkg/template"

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
