package routes

import (
	"main/pkg/template"
	"net/http"
)

func OssContributionSponsorsHandler(w http.ResponseWriter, r *http.Request) {
	template.UseTemplate(&w, r, "admin/ossContributionSponsors/index", nil)
}

func OssContributionSponsorsForm(w http.ResponseWriter, r *http.Request) {
	template.UseTemplate(&w, r, "admin/ossContributionSponsors/form", nil)
}
