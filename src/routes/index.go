package routes

import (
	template "main/pkg/template"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	template.UseTemplate(&w, r, "index", nil)
}
