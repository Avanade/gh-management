package template

import (
	"fmt"
	"html/template"
	"main/models"
	session "main/pkg/session"
	"net/http"
	"strings"
)

// This parses the master page layout and the required page template.
func UseTemplate(w *http.ResponseWriter, r *http.Request, page string, pageData interface{}) error {

	sessionaz, err := session.Store.Get(r, "auth-session")
	if err != nil {
		http.Error(*w, err.Error(), http.StatusInternalServerError)
		return err
	}

	// Data on master page
	var menu []models.TypMenu
	menu = append(menu, models.TypMenu{Name: "Approvals", Url: "/approvals/my", IconPath: "/public/icons/approvals.svg"})
	masterPageData := models.TypHeaders{Menu: menu, Page: getUrlPath(r.URL.Path)}

	data := models.TypPageData{
		Header:  masterPageData,
		Profile: sessionaz.Values["profile"],
		Content: pageData}

	tmpl := template.Must(
		template.ParseFiles("templates/master.html",
			fmt.Sprintf("templates/%v.html", page)))
	return tmpl.Execute(*w, data)
}

func getUrlPath(path string) string {
	p := strings.Split(path, "/")
	if p[1] == "" {
		return "Approvals"
	} else {
		return strings.Title(p[1])
	}
}
