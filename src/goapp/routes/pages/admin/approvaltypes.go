package routes

import (
	"main/pkg/template"
	"net/http"
)

func ApprovalTypes(w http.ResponseWriter, r *http.Request) {
	template.UseTemplate(&w, r, "admin/approvaltypes", nil)
}
