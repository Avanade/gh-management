package routes

import (
	"log"
	"net/http"
	"os"

	db "main/pkg/ghmgmtdb"
	ghAPI "main/pkg/github"
	"main/pkg/session"
	"main/pkg/template"

	"github.com/gorilla/mux"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	isAdmin, _ := session.IsUserAdmin(w, r)
	sessiongh, _ := session.GetGitHubUserData(w, r)

	data := map[string]interface{}{
		"profileGH": sessiongh,
		"isAdmin":   isAdmin,
	}
	template.UseTemplate(&w, r, "projects/index", data)
}

func MakePublicHandler(w http.ResponseWriter, r *http.Request) {
	template.UseTemplate(&w, r, "projects/makepublic", nil)
}

func FormHandler(w http.ResponseWriter, r *http.Request) {
	req := mux.Vars(r)
	id := req["id"]

	sessionaz, _ := session.Store.Get(r, "auth-session")
	sessiongh, _ := session.GetGitHubUserData(w, r)
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := profile["preferred_username"]

	token := os.Getenv("GH_TOKEN")
	innerSourceOrgName := os.Getenv("GH_ORG_INNERSOURCE")
	openSourceOrgName := os.Getenv("GH_ORG_OPENSOURCE")
	isInvalidToken := false

	isInnerSourceMember, errInnerSource := ghAPI.IsOrganizationMember(token, innerSourceOrgName, sessiongh.Username)
	if errInnerSource != nil {
		log.Println(errInnerSource.Error())
		isInvalidToken = true
	}

	isOpenSourceMember, errOpenSource := ghAPI.IsOrganizationMember(token, openSourceOrgName, sessiongh.Username)
	if errOpenSource != nil {
		log.Println(errOpenSource.Error())
		isInvalidToken = true
	}

	users := db.GetUsersWithGithub()
	data := map[string]interface{}{
		"Id":                  id,
		"users":               users,
		"email":               username,
		"isInnersourceMember": isInnerSourceMember,
		"isOpensourceMember":  isOpenSourceMember,
		"innersourceOrg":      innerSourceOrgName,
		"opensourceOrg":       openSourceOrgName,
		"isInvalidToken":      isInvalidToken,
	}

	template.UseTemplate(&w, r, "projects/form", data)
}

func ViewHandler(w http.ResponseWriter, r *http.Request) {
	template.UseTemplate(&w, r, "projects/view", nil)
}
