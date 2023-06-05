package routes

import (
	"log"
	"main/pkg/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func ListApprovalTypes(w http.ResponseWriter, r *http.Request) {
	template.UseTemplate(&w, r, "admin/approvaltypes/index", nil)
}

func ApprovalTypeForm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	action := vars["action"]
	template.UseTemplate(&w, r, "admin/approvaltypes/form", struct {
		Id     int
		Action string
	}{
		Id:     id,
		Action: strings.Title(action),
	})
}
