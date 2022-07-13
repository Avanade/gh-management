package template

import (
	"fmt"
	"html/template"
	"main/models"
	session "main/pkg/session"
	"net/http"
	"os"
	"strings"
)

// This parses the master page layout and the required page template.
func UseTemplate(w *http.ResponseWriter, r *http.Request, page string, pageData interface{}) error {

	sessionaz, err := session.Store.Get(r, "auth-session")
	if err != nil {
		http.Error(*w, err.Error(), http.StatusInternalServerError)
		return err
	}

	sessiongh, err := session.GetGitHubUserData(*w, r)
	if err != nil {
		return err
	}

	isAdmin, err := session.IsUserAdmin(*w, r)
	if err != nil {
		return err
	}

	approvalSystemUrl := os.Getenv("APPROVAL_SYSTEM_APP_URL")
	// Data on master page
	var menu []models.TypMenu
	menu = append(menu, models.TypMenu{Name: "Dashboard", Url: "/", IconPath: "/public/icons/dashboard.svg"})
	menu = append(menu, models.TypMenu{Name: "Projects", Url: "/projects", IconPath: "/public/icons/projects.svg"})
	menu = append(menu, models.TypMenu{Name: "Communities", Url: "/communities/list", IconPath: "/public/icons/communities.svg"})
	menu = append(menu, models.TypMenu{Name: "Activities", Url: "/activities", IconPath: "/public/icons/activity.svg"})
	menu = append(menu, models.TypMenu{Name: "Guidance", Url: "/guidance", IconPath: "/public/icons/guidance.svg"})
	menu = append(menu, models.TypMenu{Name: "Approvals", Url: approvalSystemUrl, IconPath: "/public/icons/approvals.svg"})
	//menu = append(menu, models.TypMenu{Name: "Search", Url: "/search", IconPath: "/public/icons/search.svg"})
	if isAdmin {
		menu = append(menu, models.TypMenu{Name: "Admin", Url: "/admin/members", IconPath: "/public/icons/lock.svg"})
	}

	var externalLinks []models.TypMenu
	externalLinks = append(externalLinks, models.TypMenu{Name: "Tech Community Calendar", Url: "https://techcommunitycalendar.com/", IconPath: "/public/icons/calendar.svg"})
	externalLinks = append(externalLinks, models.TypMenu{Name: "Stack Overflow at Avanade", Url: "https://avanade.stackenterprise.co/", IconPath: "/public/icons/questionmark.svg"})
	externalLinks = append(externalLinks, models.TypMenu{Name: "Open Innovation Meetup", Url: "https://opentechmeetup.com/", IconPath: "/public/icons/microphone.svg"})
	masterPageData := models.TypHeaders{Menu: menu, ExternalLinks: externalLinks, Page: getUrlPath(r.URL.Path)}

	data := models.TypPageData{
		Header:    masterPageData,
		Profile:   sessionaz.Values["profile"],
		ProfileGH: sessiongh,
		Content:   pageData}

	tmpl := template.Must(
		template.ParseFiles("templates/master.html",
			fmt.Sprintf("templates/%v.html", page)))
	return tmpl.Execute(*w, data)
}

func getUrlPath(path string) string {
	p := strings.Split(path, "/")
	if p[1] == "" {
		return "Dashboard"
	} else {
		return strings.Title(p[1])
	}
}
