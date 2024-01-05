package routes

import (
	"net/http"

	"main/pkg/template"
)

func AdminIndexHandler(w http.ResponseWriter, r *http.Request) {
	template.UseTemplate(&w, r, "admin/index", nil)
}
