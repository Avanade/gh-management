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

	hasPhoto, userPhoto, err := session.GetUserPhoto(*w, r)
	if err != nil {
		return err
	}

	approvalSystemUrl := os.Getenv("APPROVAL_SYSTEM_APP_URL")
	// Data on master page
	var menu []models.TypMenu
	menu = append(menu, models.TypMenu{Name: "Dashboard", Url: "/", IconPath: "/public/icons/dashboard.svg", External: false})
	menu = append(menu, models.TypMenu{Name: "Repositories", Url: "/repositories", IconPath: "/public/icons/projects.svg", External: false})
	menu = append(menu, models.TypMenu{Name: "Communities", Url: "/communities/list", IconPath: "/public/icons/communities.svg", External: false})
	menu = append(menu, models.TypMenu{Name: "Activities", Url: "/activities", IconPath: "/public/icons/activity.svg", External: false})
	menu = append(menu, models.TypMenu{Name: "Guidance", Url: "/guidance", IconPath: "/public/icons/guidance.svg", External: false})
	menu = append(menu, models.TypMenu{Name: "Approvals", Url: approvalSystemUrl, IconPath: "/public/icons/approvals.svg", External: true})
	if isAdmin {
		menu = append(menu, models.TypMenu{Name: "Admin", Url: "/admin", IconPath: "/public/icons/lock.svg", External: false})
	}

	var externalLinks []models.TypMenu
	externalLinks = append(externalLinks, models.TypMenu{Name: "Tech Community Calendar", Url: "https://techcommunitycalendar.com/", IconPath: "/public/icons/calendar.svg", External: true})
	externalLinks = append(externalLinks, models.TypMenu{Name: "Stack Overflow at Avanade", Url: "https://avanade.stackenterprise.co/", IconPath: "/public/icons/questionmark.svg", External: true})
	externalLinks = append(externalLinks, models.TypMenu{Name: "Open Innovation Meetup", Url: "https://opentechinnovation.com/", IconPath: "/public/icons/microphone.svg", External: true})
	masterPageData := models.TypHeaders{Menu: menu, ExternalLinks: externalLinks, Page: getUrlPath(r.URL.Path)}

	data := models.TypPageData{
		Header:    masterPageData,
		Profile:   sessionaz.Values["profile"],
		ProfileGH: sessiongh,
		Content:   pageData,
		HasPhoto:  hasPhoto,
		UserPhoto: userPhoto}

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
