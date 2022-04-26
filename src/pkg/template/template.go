package template

import (
	"fmt"
	"html/template"
	"main/models"
	"net/http"
)

// This parses the master page layout and the required page template.
func UseTemplate(w *http.ResponseWriter, data *models.TypPageData, page string) error {
	tmpl := template.Must(
		template.ParseFiles("templates/master.html",
			fmt.Sprintf("templates/%v.html", page)))
	return tmpl.Execute(*w, *data)
}
