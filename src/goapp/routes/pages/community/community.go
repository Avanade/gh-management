package routes

import (
	"net/http"

	"main/pkg/session"
	"main/pkg/template"

	"github.com/gorilla/mux"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	template.UseTemplate(&w, r, "community/index", nil)
}

func FormHandler(w http.ResponseWriter, r *http.Request) {
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

	template.UseTemplate(&w, r, "/community/form", data)
}

func OnBoardingHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id := params["id"]

	data := map[string]string{
		"Id": id,
	}
	template.UseTemplate(&w, r, "community/onboarding", data)
}

// NOT YET FINAL
func CommunityApproversHandler(w http.ResponseWriter, r *http.Request) {
	template.UseTemplate(&w, r, "/community/communityapprovers", nil)
}
