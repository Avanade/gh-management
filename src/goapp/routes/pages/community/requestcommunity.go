package routes

import (
	template "main/pkg/template"
	"net/http"
)

func CommunityHandler(w http.ResponseWriter, r *http.Request) {

	template.UseTemplate(&w, r, "community/requestcommunity", nil)
}
