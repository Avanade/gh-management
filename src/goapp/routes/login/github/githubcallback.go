package routes

import (
	"encoding/json"
	"fmt"
	auth "main/pkg/authentication"
	githubAPI "main/pkg/github"
	session "main/pkg/session"
	"net/http"
	"os"
	"strconv"

	"golang.org/x/oauth2"

	db "main/pkg/ghmgmtdb"
	ghmgmt "main/pkg/ghmgmtdb"
	"main/pkg/msgraph"
)

func GithubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	sessionaz, err := session.Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Check session and state
	state, err := session.GetState(w, r)

	session, err := session.Store.Get(r, "gh-auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.URL.Query().Get("state") != state {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	ghauth := auth.GetGitHubOauthConfig()

	// Exchange temporary code for access token
	code := r.URL.Query().Get("code")

	ghAccessToken, err := ghauth.Exchange(oauth2.NoContext, code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ghProfile := auth.GetGitHubUserProfile(ghAccessToken.AccessToken)

	// Save GitHub auth data on session cookies
	session.Values["ghAccessToken"] = ghAccessToken.AccessToken
	session.Values["ghProfile"] = ghProfile

	// Convert string to interface{} array
	var p map[string]interface{}
	err = json.Unmarshal([]byte(ghProfile), &p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Save and Validate github account
	azProfile := sessionaz.Values["profile"].(map[string]interface{})
	userPrincipalName := fmt.Sprintf("%s", azProfile["preferred_username"])
	ghId := strconv.FormatFloat(p["id"].(float64), 'f', 0, 64)
	ghUser := fmt.Sprintf("%s", p["login"])
	id, err := strconv.ParseInt(ghId, 10, 64)
	resultUUG, errUUG := ghmgmt.UpdateUserGithub(userPrincipalName, ghId, ghUser, 0)
	if errUUG != nil {
		http.Error(w, errUUG.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["ghIsValid"] = resultUUG["IsValid"].(bool)

	isDirect, _ := msgraph.IsDirectMember(fmt.Sprintf("%s", azProfile["oid"]))
	isEnterpriseMember, _ := msgraph.IsGithubEnterpriseMember(fmt.Sprintf("%s", azProfile["oid"]))

	session.Values["ghIsDirect"] = isDirect
	session.Values["ghIsEnterpriseMember"] = isEnterpriseMember

	CheckMembership(ghUser, &id)

	err = session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func GithubForceSaveHandler(w http.ResponseWriter, r *http.Request) {
	sessionaz, err := session.Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Check session and state
	session, err := session.Store.Get(r, "gh-auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ghProfile := session.Values["ghProfile"].(string)

	var p map[string]interface{}
	err = json.Unmarshal([]byte(ghProfile), &p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Save and Validate github account
	azProfile := sessionaz.Values["profile"].(map[string]interface{})
	userPrincipalName := fmt.Sprintf("%s", azProfile["preferred_username"])
	ghId := strconv.FormatFloat(p["id"].(float64), 'f', 0, 64)
	ghUser := fmt.Sprintf("%s", p["login"])

	resultUUG, errUUG := ghmgmt.UpdateUserGithub(userPrincipalName, ghId, ghUser, 1)
	if errUUG != nil {
		http.Error(w, errUUG.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["ghIsValid"] = resultUUG["IsValid"].(bool)

	err = session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func CheckMembership(ghusername string, id *int64) {
	token := os.Getenv("GH_TOKEN")
	inner, outer, _ := githubAPI.OrganizationsIsMember(token, ghusername)
	if !inner {
		githubAPI.OrganizationInvitation(token, ghusername, os.Getenv("GH_ORG_INNERSOURCE"))

	}
	if !outer {
		githubAPI.OrganizationInvitation(token, ghusername, os.Getenv("GH_ORG_OPENSOURCE"))

	}
}

func CheckAvaInnerSource(w http.ResponseWriter, r *http.Request) {

	org := os.Getenv("GH_ORG_INNERSOURCE")
	token := os.Getenv("GH_TOKEN")

	collabs := githubAPI.ListOutsideCollaborators(token, org)
	for _, collab := range collabs {
		githubAPI.RemoveOutsideCollaborator(token, org, *collab.Login)
	}
}

func CheckAvaOpenSource(w http.ResponseWriter, r *http.Request) {
	org := os.Getenv("GH_ORG_OPENSOURCE")
	var OutsidecollabsList []string
	token := os.Getenv("GH_TOKEN")
	repos, _ := githubAPI.GetRepositoriesFromOrganization(org)
	Outsidecollabs := githubAPI.ListOutsideCollaborators(token, org)
	for _, list := range Outsidecollabs {
		OutsidecollabsList = append(OutsidecollabsList, *list.Login)
	}
	var OutsideRepocollabsList []string
	for _, collab := range repos {
		var RepocollabsList []string

		var Adminmember []string
		OutsideRepocollabsList = nil

		Repocollabs := githubAPI.RepositoriesListCollaborators(token, org, collab.Name)
		for _, list := range Repocollabs {

			RepocollabsList = append(RepocollabsList, *list.Login)
			if *list.RoleName == "admin" {
				Adminmember = append(Adminmember, *list.Login)

			}
		}

		for _, list := range RepocollabsList {
			for _, Outsidelist := range OutsidecollabsList {
				if list == Outsidelist {
					OutsideRepocollabsList = append(OutsideRepocollabsList, Outsidelist)
				}
			}
		}
		if len(OutsideRepocollabsList) > 0 {

			for _, admin := range Adminmember {
				email, err := db.UsersGetEmail(admin)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				githubAPI.EmailAdmin(admin, email, collab.Name, OutsideRepocollabsList)
			}

		}

	}

}

func ClearOrgMembers(w http.ResponseWriter, r *http.Request) {
	token := os.Getenv("GH_TOKEN")

	// Remove GitHub users from innersource who are not employees
	organization := os.Getenv("GH_ORG_INNERSOURCE")
	EmailSupport := os.Getenv("EMAIL_SUPPORT")
	var ConvertedOutsidecollabsList []string
	users := githubAPI.OrgListMembers(token, organization)
	for _, list := range users {
		email, _ := db.UsersGetEmail(*list.Login)
		if len(email) > 0 {
			activeuser, _ := msgraph.ActiveUsers(email)
			if activeuser == nil {
				githubAPI.RemoveOrganizationsMember(token, organization, *list.Login)

			}
		} else {
			githubAPI.RemoveOrganizationsMember(token, organization, *list.Login)

		}

	}

	// Convert users who are not employees to an outside collaborator
	organizationsOpen := os.Getenv("GH_ORG_OPENSOURCE")

	usersOpenorg := githubAPI.OrgListMembers(token, organizationsOpen)
	for _, list := range usersOpenorg {
		email, _ := db.UsersGetEmail(*list.Login)
		if len(email) > 0 {
			activeuser, _ := msgraph.ActiveUsers(email)
			if activeuser == nil {
				githubAPI.ConvertMemberToOutsideCollaborator(token, organizationsOpen, *list.Login)
				ConvertedOutsidecollabsList = append(ConvertedOutsidecollabsList, *list.Login)
			}
		} else {
			githubAPI.ConvertMemberToOutsideCollaborator(token, organizationsOpen, *list.Login)
			ConvertedOutsidecollabsList = append(ConvertedOutsidecollabsList, *list.Login)
		}

	}

	if len(ConvertedOutsidecollabsList) > 0 {
		// Email list of new outside collaborators to ospo
		githubAPI.EmailAdminConvertToColaborator(EmailSupport, ConvertedOutsidecollabsList)

		// Email repo admins with converted users
		repos, _ := githubAPI.GetRepositoriesFromOrganization(organizationsOpen)
		for _, repo := range repos {

			RepoAdmins := githubAPI.GetRepoAdmin(organizationsOpen, repo.Name)
			Repocollabs := githubAPI.RepositoriesListCollaborators(token, organizationsOpen, repo.Name)
			var ConvertedInRepo []string
			for _, collab1 := range ConvertedOutsidecollabsList {
				for _, collab2 := range Repocollabs {
					if collab1 == *collab2.Login {
						ConvertedInRepo = append(ConvertedInRepo, collab1)
					}
				}
			}

			for _, collab := range RepoAdmins {
				collabemail, _ := db.UsersGetEmail(collab)

				if len(ConvertedInRepo) > 0 {
					githubAPI.EmailRepoAdminConvertToColaborator(collabemail, repo.Name, ConvertedInRepo)
				}
			}

		}
	}

}
