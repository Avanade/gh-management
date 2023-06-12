package routes

import (
	"net/http"
	"os"

	db "main/pkg/ghmgmtdb"
	ghAPI "main/pkg/github"
	"main/pkg/session"
	"main/pkg/template"

	"github.com/gorilla/mux"
)

func Projects(w http.ResponseWriter, r *http.Request) {
	isAdmin, _ := session.IsUserAdmin(w, r)

	data := map[string]interface{}{
		"isAdmin": isAdmin,
	}
	template.UseTemplate(&w, r, "projects/projects", data)
}

func MakePublic(w http.ResponseWriter, r *http.Request) {
	template.UseTemplate(&w, r, "projects/makepublic", nil)
}

func ProjectById(w http.ResponseWriter, r *http.Request) {
	req := mux.Vars(r)
	id := req["id"]

	users := db.GetUsersWithGithub()
	data := map[string]interface{}{
		"Id":    id,
		"users": users,
	}
	template.UseTemplate(&w, r, "projects/new", data)
}

func NewProject(w http.ResponseWriter, r *http.Request) {
	req := mux.Vars(r)
	id := req["id"]

	sessionaz, _ := session.Store.Get(r, "auth-session")
	sessiongh, _ := session.GetGitHubUserData(w, r)
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := profile["preferred_username"]
	isInnersourceMember, isOpensourceMember, _ := ghAPI.OrganizationsIsMember(os.Getenv("GH_TOKEN"), sessiongh.Username)

	users := db.GetUsersWithGithub()
	data := map[string]interface{}{
		"Id":                  id,
		"users":               users,
		"email":               username,
		"isInnersourceMember": isInnersourceMember,
		"isOpensourceMember":  isOpensourceMember,
		"innersourceOrg":      os.Getenv("GH_ORG_INNERSOURCE"),
		"opensourceOrg":       os.Getenv("GH_ORG_OPENSOURCE"),
	}
	template.UseTemplate(&w, r, "projects/new", data)
}
