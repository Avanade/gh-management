package routes

import (
	template "main/pkg/template"
	"net/http"

	"strconv"
	"strings"
	"github.com/gorilla/mux"
)

func ExternalLinksHandler(w http.ResponseWriter, r *http.Request) {
	template.UseTemplate(&w, r, "admin/externallinks", nil)
}
func ExternalLinksForm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	action := vars["action"]
	template.UseTemplate(&w, r, "admin/externallinks/form", struct {
		Id     int
		Action string
	}{
		Id:     id,
		Action: strings.Title(action),
	})
}