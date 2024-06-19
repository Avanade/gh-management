package routes

import (
	"net/http"
	"os"
	"strconv"

	"main/pkg/template"

	"github.com/gorilla/mux"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	template.UseTemplate(&w, r, "activities/index", struct {
		OrganizationName string
	}{
		OrganizationName: os.Getenv("ORGANIZATION_NAME"),
	})
}

func FormHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	action := vars["action"]

	caser := cases.Title(language.Und, cases.NoLower)

	template.UseTemplate(&w, r, "activities/form", struct {
		Id               int
		Action           string
		OrganizationName string
	}{
		Id:               id,
		Action:           caser.String(action),
		OrganizationName: os.Getenv("ORGANIZATION_NAME"),
	})
}
