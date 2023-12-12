package routes

import (
	"net/http"
	"strconv"
	"strings"

	"main/pkg/template"

	"github.com/gorilla/mux"
)

func ContributionAreasHandler(w http.ResponseWriter, r *http.Request) {
	template.UseTemplate(&w, r, "admin/contributionareas/index", nil)
}

func ContributionAreasFormHandler(w http.ResponseWriter, r *http.Request) {
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
