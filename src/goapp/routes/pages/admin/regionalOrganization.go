package routes

import (
	"net/http"
	"os"

	"main/pkg/template"
)

func RegionalOrganizationHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"OrganizationName": os.Getenv("ORGANIZATION_NAME"),
	}
	template.UseTemplate(&w, r, "admin/manageorganizations/index", data)
}
