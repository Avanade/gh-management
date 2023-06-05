package routes

import (
	"encoding/json"
	"fmt"
	"log"
	auth "main/pkg/authentication"
	"main/pkg/email"
	githubAPI "main/pkg/github"
	session "main/pkg/session"
	"net/http"
	"os"
	"strconv"
	"time"

	"golang.org/x/oauth2"

	ghmgmt "main/pkg/ghmgmtdb"
	"main/pkg/msgraph"
)

func GithubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	sessionaz, err := session.Store.Get(r, "auth-session")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Check session and state
	state, err := session.GetState(w, r)

	session, err := session.Store.Get(r, "gh-auth-session")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.URL.Query().Get("state") != state {
		log.Println("Invalid state paramerer")
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	ghauth := auth.GetGitHubOauthConfig(r.Host)

	// Exchange temporary code for access token
	code := r.URL.Query().Get("code")

	ghAccessToken, err := ghauth.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Println(err.Error())
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
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Save and Validate github account
	azProfile := sessionaz.Values["profile"].(map[string]interface{})
	userPrincipalName := fmt.Sprintf("%s", azProfile["preferred_username"])
	ghId := strconv.FormatFloat(p["id"].(float64), 'f', 0, 64)
	ghUser := fmt.Sprintf("%s", p["login"])
	id, err := strconv.ParseInt(ghId, 10, 64)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := ghmgmt.UpdateUserGithub(userPrincipalName, ghId, ghUser, 0)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["ghIsValid"] = result["IsValid"].(bool)

	isDirect, _ := msgraph.IsDirectMember(fmt.Sprintf("%s", azProfile["oid"]))
	isEnterpriseMember, _ := msgraph.IsGithubEnterpriseMember(fmt.Sprintf("%s", azProfile["oid"]))

	session.Values["ghIsDirect"] = isDirect
	session.Values["ghIsEnterpriseMember"] = isEnterpriseMember

	CheckMembership(userPrincipalName, ghUser, &id)

	err = session.Save(r, w)
	if err != nil {
		log.Panicln(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func GithubForceSaveHandler(w http.ResponseWriter, r *http.Request) {
	sessionaz, err := session.Store.Get(r, "auth-session")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Check session and state
	session, err := session.Store.Get(r, "gh-auth-session")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ghProfile := session.Values["ghProfile"].(string)

	var p map[string]interface{}
	err = json.Unmarshal([]byte(ghProfile), &p)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Save and Validate github account
	azProfile := sessionaz.Values["profile"].(map[string]interface{})
	userPrincipalName := fmt.Sprintf("%s", azProfile["preferred_username"])
	ghId := strconv.FormatFloat(p["id"].(float64), 'f', 0, 64)
	ghUser := fmt.Sprintf("%s", p["login"])

	result, err := ghmgmt.UpdateUserGithub(userPrincipalName, ghId, ghUser, 1)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["ghIsValid"] = result["IsValid"].(bool)

	err = session.Save(r, w)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func CheckMembership(userPrincipalName, ghusername string, id *int64) {
	token := os.Getenv("GH_TOKEN")
	inner, outer, _ := githubAPI.OrganizationsIsMember(token, ghusername)
	if !inner {
		githubAPI.OrganizationInvitation(token, ghusername, os.Getenv("GH_ORG_INNERSOURCE"))
	}
	if !outer {
		githubAPI.OrganizationInvitation(token, ghusername, os.Getenv("GH_ORG_OPENSOURCE"))

	}
	EmailAcceptOrgInvitation(userPrincipalName, ghusername, inner, outer)
}

func EmailAcceptOrgInvitation(userEmail, ghUsername string, isInnersourceOrgMember, isOpensourceOrgMember bool) {
	opensourceLink := orgInvitationLink(os.Getenv("GH_ORG_OPENSOURCE"))
	innersourceLink := orgInvitationLink(os.Getenv("GH_ORG_INNERSOURCE"))

	body := fmt.Sprintf("Hello %s, <br>", ghUsername)

	if !isInnersourceOrgMember && !isOpensourceOrgMember {
		body = body + fmt.Sprintf("<p>An invitation to join %s and %s was created. You need to join the %s so you get added to the repository you'll request.</p>", innersourceLink, opensourceLink, innersourceLink)
	} else if !isInnersourceOrgMember {
		body = body + fmt.Sprintf("<p>An invitation to join %s was created. You need to join this organization so you get added to the repository you'll request.</p>", innersourceLink)
	} else if !isOpensourceOrgMember {
		body = body + fmt.Sprintf("<p>An invitation to join %s was created. You need to join this organization so you won't get tagged as an outside collaborator.</p>", opensourceLink)
	} else {
		return
	}

	body = body + "<br>OSPO"

	m := email.TypEmailMessage{
		Subject: "Github Organizations Invitation",
		Body:    body,
		To:      userEmail,
	}

	email.SendEmail(m)
	fmt.Printf("GitHub Organization Invitation on %s was sent.", time.Now())
}

func orgInvitationLink(org string) string {
	url := fmt.Sprintf("https://github.com/orgs/%s/invitation", org)
	return fmt.Sprintf("<a href='%s'>%s</a>", url, org)
}
