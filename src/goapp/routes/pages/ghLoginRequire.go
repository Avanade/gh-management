package routes

import (
	"net/http"

	"main/pkg/template"
)

func GHLoginRequire(w http.ResponseWriter, r *http.Request) {
	template.UseTemplate(&w, r, "ghloginrequire", nil)
}
