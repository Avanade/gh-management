package routes

import (
	"log"
	"net/http"
	"os"
	"strconv"

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
		"profileGH":                             sessiongh,
		"isAdmin":                               isAdmin,
		"innersource":                           os.Getenv("GH_ORG_INNERSOURCE"),
		"opensource":                            os.Getenv("GH_ORG_OPENSOURCE"),
		"innersourceGeneralLegalGuidelinesLink": os.Getenv("LINK_INNERSOURCE_GENERAL_LEGAL_GUIDELINES"),
	}
	template.UseTemplate(&w, r, "projects/index", data)
}

func MakePublicHandler(w http.ResponseWriter, r *http.Request) {
	token := os.Getenv("GH_TOKEN")
	innerSourceOrgName := os.Getenv("GH_ORG_INNERSOURCE")
	openSourceOrgName := os.Getenv("GH_ORG_OPENSOURCE")
	sessiongh, _ := session.GetGitHubUserData(w, r)
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

	data := map[string]interface{}{
		"isInvalidToken":      isInvalidToken,
		"isInnerSourceMember": isInnerSourceMember,
		"isOpenSourceMember":  isOpenSourceMember,
		"innersourceOrg":      innerSourceOrgName,
		"opensourceOrg":       openSourceOrgName,
		"OrganizationName":    os.Getenv("ORGANIZATION_NAME"),
	}

	template.UseTemplate(&w, r, "projects/makepublic", data)
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
		"Id":                           id,
		"users":                        users,
		"email":                        username,
		"isInnersourceMember":          isInnerSourceMember,
		"isOpensourceMember":           isOpenSourceMember,
		"innersourceOrg":               innerSourceOrgName,
		"opensourceOrg":                openSourceOrgName,
		"isInvalidToken":               isInvalidToken,
		"innersourceGeneralGuidelines": os.Getenv("LINK_INNERSOURCE_GENERAL_GUIDELINES"),
		"OrganizationName":             os.Getenv("ORGANIZATION_NAME"),
	}

	template.UseTemplate(&w, r, "projects/form", data)
}

func ViewByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	isAdmin, _ := session.IsUserAdmin(w, r)
	sessiongh, _ := session.GetGitHubUserData(w, r)

	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := profile["preferred_username"]

	id, err := strconv.ParseInt(vars["githubId"], 10, 64)
	if err != nil {
		http.Redirect(w, r, "/repositories", http.StatusNotFound)
		return
	}

	projects := db.GetProjectByGithubId(id)
	if projects == nil {
		http.Redirect(w, r, "/repositories", http.StatusNotFound)
		return
	}

	projectId := projects[0]["Id"].(int64)
	orgName := projects[0]["Organization"].(string)

	repoOwners, err := db.RepoOwnersByUserAndProjectId(projectId, username.(string))
	if err != nil {
		log.Println(err)
		return
	}

	isOwner := false

	if len(repoOwners) > 0 {
		isOwner = true
	}

	token := os.Getenv("GH_TOKEN")
	isMember, err := ghAPI.IsOrganizationMember(token, orgName, sessiongh.Username)
	if err != nil {
		log.Println(err)
		return
	}

	data := map[string]interface{}{
		"id":        projectId,
		"profileGH": sessiongh,
		"isAdmin":   isAdmin,
		"isOwner":   isOwner,
		"isMember":  isMember,
		"orgName":   orgName,
	}

	template.UseTemplate(&w, r, "projects/view", data)
}

func ViewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	isAdmin, _ := session.IsUserAdmin(w, r)
	sessiongh, _ := session.GetGitHubUserData(w, r)

	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := profile["preferred_username"]

	orgName := vars["org"]

	id, err := db.GetProjectIdByOrgName(orgName, vars["repo"])
	if err != nil {
		http.Redirect(w, r, "/repositories", http.StatusNotFound)
		return
	}

	repoOwners, err := db.RepoOwnersByUserAndProjectId(id, username.(string))
	if err != nil {
		log.Println(err)
		return
	}

	isOwner := false

	if len(repoOwners) > 0 {
		isOwner = true
	}

	token := os.Getenv("GH_TOKEN")
	isMember, err := ghAPI.IsOrganizationMember(token, orgName, sessiongh.Username)
	if err != nil {
		log.Println(err)
		return
	}

	data := map[string]interface{}{
		"id":        id,
		"profileGH": sessiongh,
		"isAdmin":   isAdmin,
		"isOwner":   isOwner,
		"isMember":  isMember,
		"orgName":   orgName,
	}

	template.UseTemplate(&w, r, "projects/view", data)
}
