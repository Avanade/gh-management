package routes

import (
	"net/http"

	"main/pkg/template"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	template.UseTemplate(&w, r, "notfound", nil)
}
