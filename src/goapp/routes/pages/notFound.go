package routes

import (
	"net/http"

	"main/pkg/template"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	template.UseTemplate(&w, r, "notfound", nil)
}
