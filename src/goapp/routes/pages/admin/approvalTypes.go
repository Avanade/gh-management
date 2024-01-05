package routes

import (
	"net/http"
	"strconv"

	"main/pkg/template"

	"github.com/gorilla/mux"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func ApprovalTypesHandler(w http.ResponseWriter, r *http.Request) {
	template.UseTemplate(&w, r, "admin/approvaltypes/index", nil)
}

func ApprovalTypeFormHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	action := vars["action"]
	caser := cases.Title(language.Und, cases.NoLower)
	template.UseTemplate(&w, r, "admin/approvaltypes/form", struct {
		Id     int
		Action string
	}{
		Id:     id,
		Action: caser.String(action),
	})
}
