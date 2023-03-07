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
