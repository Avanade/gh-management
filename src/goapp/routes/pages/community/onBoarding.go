package routes

import (
	"net/http"

	"main/pkg/template"

	"github.com/gorilla/mux"
)

func CommunityOnBoarding(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id := params["id"]

	data := map[string]string{
		"Id": id,
	}
	template.UseTemplate(&w, r, "community/onboarding", data)
}
