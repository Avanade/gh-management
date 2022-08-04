package routes

import (
	template "main/pkg/template"
	"net/http"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	template.UseTemplate(&w, r, "NotFound", nil)
}
