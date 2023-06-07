package routes

import (
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
	id, _ := strconv.Atoi(vars["id"])

	action := vars["action"]
	template.UseTemplate(&w, r, "admin/approvaltypes/form", struct {
		Id     int
		Action string
	}{
		Id:     id,
		Action: strings.Title(action),
	})
}
