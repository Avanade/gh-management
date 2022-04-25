package template

import (
	"fmt"
	"html/template"
)

func UseTemplate(page string) *template.Template {
	return template.Must(
		template.ParseFiles("main.html",
			fmt.Sprintf("templates/%v.html", page)))
}
