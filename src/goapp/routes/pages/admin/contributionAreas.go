package routes

import (
	"main/pkg/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func ListContributionAreas(w http.ResponseWriter, r *http.Request) {
	template.UseTemplate(&w, r, "admin/contributionareas/index", nil)
}

func ContributionAreasForm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	action := vars["action"]
	template.UseTemplate(&w, r, "admin/contributionareas/form", struct {
		Id     int
		Action string
	}{
		Id:     id,
		Action: strings.Title(action),
	})
}
