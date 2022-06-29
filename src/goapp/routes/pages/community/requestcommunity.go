package routes

import (
	template "main/pkg/template"
	"net/http"

	"github.com/gorilla/mux"
)

func CommunityHandler(w http.ResponseWriter, r *http.Request) {
	req := mux.Vars(r)
	id := req["id"]

	data := map[string]interface{}{
		"Id": id,
	}

	template.UseTemplate(&w, r, "/community/requestcommunity", data)
}
