package template

import (
	"fmt"
	"html/template"
)

// This parses the master page layout and the required page template.
func UseTemplate(page string) *template.Template {
	return template.Must(
		template.ParseFiles("templates/master.html",
			fmt.Sprintf("templates/%v.html", page)))
}
