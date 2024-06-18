package routes

import (
	"net/http"
	"os"

	"main/pkg/template"
)

func AdminIndexHandler(w http.ResponseWriter, r *http.Request) {
	template.UseTemplate(&w, r, "admin/index", struct {
		OrganizationName string
	}{
		OrganizationName: os.Getenv("ORGANIZATION_NAME"),
	})
}
