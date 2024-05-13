package template

import (
	"fmt"
	"html/template"
	session "main/pkg/session"
	"net/http"
	"os"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type PageData struct {
	Header         interface{}
	Profile        interface{}
	ProfileGH      interface{}
	Content        interface{}
	IsGHAssociated bool
	HasPhoto       bool
	UserPhoto      string
}

type Headers struct {
	Menu     []Menu
	Title    string
	LogoPath string
	Page     string
}

type Menu struct {
	Name     string
	Url      string
	IconPath string
	External bool
}

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
	title := os.Getenv("APP_TITLE")
	logoPath := os.Getenv("APP_LOGO_PATH")
	// Data on master page
	var menu []Menu
	menu = append(menu, Menu{Name: "Dashboard", Url: "/", IconPath: "/public/icons/dashboard.svg", External: false})
	menu = append(menu, Menu{Name: "Repositories", Url: "/repositories", IconPath: "/public/icons/projects.svg", External: false})
	menu = append(menu, Menu{Name: "Communities", Url: "/communities", IconPath: "/public/icons/communities.svg", External: false})
	menu = append(menu, Menu{Name: "Activities", Url: "/activities", IconPath: "/public/icons/activity.svg", External: false})
	menu = append(menu, Menu{Name: "Guidance", Url: "/guidance", IconPath: "/public/icons/guidance.svg", External: false})
	menu = append(menu, Menu{Name: "Approvals", Url: approvalSystemUrl, IconPath: "/public/icons/approvals.svg", External: true})
	menu = append(menu, Menu{Name: "Other Requests", Url: "/other-requests", IconPath: "/public/icons/otherrequests.svg", External: false})
	if isAdmin {
		menu = append(menu, Menu{Name: "Admin", Url: "/admin", IconPath: "/public/icons/lock.svg", External: false})
	}

	masterPageData := Headers{Title: title, LogoPath: logoPath, Menu: menu, Page: GetUrlPath(r.URL.Path)}

	data := PageData{
		Header:    masterPageData,
		Profile:   sessionaz.Values["profile"],
		ProfileGH: sessiongh,
		Content:   pageData,
		HasPhoto:  hasPhoto,
		UserPhoto: userPhoto,
	}

	if sessionaz.Values["isGHAssociated"] != nil {
		data.IsGHAssociated = sessionaz.Values["isGHAssociated"].(bool)
	}

	tmpl := template.Must(
		template.ParseFiles("templates/master.html",
			fmt.Sprintf("templates/%v.html", page)))
	return tmpl.Execute(*w, data)
}

func GetUrlPath(path string) string {
	p := strings.Split(path, "/")
	if p[1] == "" {
		return "Dashboard"
	} else {
		caser := cases.Title(language.Und, cases.NoLower)
		return caser.String(p[1])
	}
}
