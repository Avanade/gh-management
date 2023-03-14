package routes

import (
	"main/pkg/template"
	"net/http"
)

func CustomizeExternalLinks(w http.ResponseWriter, r *http.Request) {
	template.UseTemplate(&w, r, "admin/externallinks", nil)
}
