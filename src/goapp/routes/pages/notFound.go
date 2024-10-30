package routes

import (
	"fmt"
	"net/http"

	"main/pkg/session"
	"main/pkg/template"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	s, err := session.Store.Get(r, "auth-session")
	state := s.Values["state"]
	if err != nil || state == nil {
		url := fmt.Sprintf("/loginredirect?redirect=%v", r.URL)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return
	}

	template.UseTemplate(&w, r, "notfound", nil)
}
