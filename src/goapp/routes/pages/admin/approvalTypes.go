package routes

import (
	"net/http"
	"strconv"
	"strings"

	"main/pkg/template"

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
