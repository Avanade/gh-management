package routes

import (
	"main/pkg/template"
	"net/http"
)

func AdminIndex(w http.ResponseWriter, r *http.Request) {
	template.UseTemplate(&w, r, "admin/index", nil)
}
