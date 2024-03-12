package routes

import (
	"main/pkg/template"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	template.UseTemplate(&w, r, "otherRequests/index", nil)
}
