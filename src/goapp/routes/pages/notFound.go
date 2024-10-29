package routes

import (
	"net/http"

	"main/pkg/appinsights_wrapper"
	"main/pkg/template"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	logger.TrackEvent("NotFoundHandler")

	template.UseTemplate(&w, r, "notfound", nil)
}
