package routes

import (
	"net/http"
	"strconv"

	"main/pkg/template"

	"github.com/gorilla/mux"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func ExternalLinksHandler(w http.ResponseWriter, r *http.Request) {
	template.UseTemplate(&w, r, "admin/externallinks/index", nil)
}
func ExternalLinksFormHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	action := vars["action"]
	caser := cases.Title(language.Und, cases.NoLower)
	template.UseTemplate(&w, r, "admin/externallinks/form", struct {
		Id     int
		Action string
	}{
		Id:     id,
		Action: caser.String(action),
	})
}
