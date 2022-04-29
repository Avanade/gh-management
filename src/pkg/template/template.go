package template

import (
	"fmt"
	"html/template"
	"main/models"
	session "main/pkg/session"
	"net/http"
)

// This parses the master page layout and the required page template.
func UseTemplate(w *http.ResponseWriter, r *http.Request, page string, pageData interface{}) error {

	session, err := session.Store.Get(r, "auth-session")
	if err != nil {
		http.Error(*w, err.Error(), http.StatusInternalServerError)
		return err
	}

	// Data on master page
	var menu []models.TypMenu
	menu = append(menu, models.TypMenu{Name: "Home", Url: "/"})
	menu = append(menu, models.TypMenu{Name: "Github", Url: "/github"})
	pageHeaders := models.TypHeaders{Menu: menu}

	data := models.TypPageData{
		Header:  pageHeaders,
		Profile: session.Values["profile"],
		Content: pageData}

	tmpl := template.Must(
		template.ParseFiles("templates/master.html",
			fmt.Sprintf("templates/%v.html", page)))
	return tmpl.Execute(*w, data)
}
