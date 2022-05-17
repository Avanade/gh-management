package routes

import (
	template "main/pkg/template"
	"net/http"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	template.UseTemplate(&w, r, "NotFound", nil)
}
