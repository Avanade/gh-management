package routes

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"main/pkg/session"
	"main/pkg/template"

	"github.com/gorilla/mux"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	sessionaz, err := session.Store.Get(r, "auth-session")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	template.UseTemplate(&w, r, "activities/index", sessionaz)
}

func FormHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	action := vars["action"]
	template.UseTemplate(&w, r, "activities/form", struct {
		Id     int
		Action string
	}{
		Id:     id,
		Action: strings.ToTitle(action),
	})
}
