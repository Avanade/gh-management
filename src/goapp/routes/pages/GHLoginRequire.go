package routes

import (
	template "main/pkg/template"
	"net/http"
)

func GHLoginRequire(w http.ResponseWriter, r *http.Request) {
	template.UseTemplate(&w, r, "GHLoginRequire", nil)
}
