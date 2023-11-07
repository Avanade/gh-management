package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	auth "main/pkg/authentication"
	"main/pkg/email"
	db "main/pkg/ghmgmtdb"
	ghAPI "main/pkg/github"
	"main/pkg/msgraph"
	"main/pkg/notification"
	"main/pkg/session"

	"golang.org/x/oauth2"
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
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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

	result, err := db.UpdateUserGithub(userPrincipalName, ghId, ghUser, 0)
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

	lastGithubLogin := result["LastGithubLogin"].(time.Time)

	if !DateEqual(lastGithubLogin) && result["IsValid"].(bool) {
		CheckMembership(userPrincipalName, ghUser)
	}

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

	result, err := db.UpdateUserGithub(userPrincipalName, ghId, ghUser, 1)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["ghIsValid"] = result["IsValid"].(bool)

	CheckMembership(userPrincipalName, ghUser)

	err = session.Save(r, w)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func CheckMembership(userPrincipalName, ghusername string) {
	token := os.Getenv("GH_TOKEN")

	innerSourceOrgName := os.Getenv("GH_ORG_INNERSOURCE")
	openSourceOrgName := os.Getenv("GH_ORG_OPENSOURCE")

	isInnerSourceMember, err := ghAPI.IsOrganizationMember(token, innerSourceOrgName, ghusername)
	if err != nil {
		isInnerSourceMember = true
		log.Println(err.Error())
	} else {
		if !isInnerSourceMember {
			ghAPI.OrganizationInvitation(token, ghusername, innerSourceOrgName)
			NotificationAccepOrgInvitation(userPrincipalName, innerSourceOrgName)
		}
	}

	isOpenSourceMember, err := ghAPI.IsOrganizationMember(token, openSourceOrgName, ghusername)
	if err != nil {
		isOpenSourceMember = true
		log.Println(err.Error())
	} else {
		if !isOpenSourceMember {
			ghAPI.OrganizationInvitation(token, ghusername, openSourceOrgName)
			NotificationAccepOrgInvitation(userPrincipalName, openSourceOrgName)
		}
	}
	EmailAcceptOrgInvitation(userPrincipalName, ghusername, isInnerSourceMember, isOpenSourceMember)
}

func NotificationAccepOrgInvitation(userEmail, org string) {
	messageBody := notification.OrganizationInvitationMessageBody{
		Recipients: []string{
			userEmail,
		},
		InvitationLink:   fmt.Sprintf("https://github.com/orgs/%s/invitation", org),
		OrganizationLink: fmt.Sprintf("https://github.com/%s", org),
		OrganizationName: org,
	}
	err := messageBody.Send()
	if err != nil {
		log.Println(err.Error())
	}
}

func EmailAcceptOrgInvitation(userEmail, ghUsername string, isInnersourceOrgMember, isOpensourceOrgMember bool) {
	opensourceLink := OrgInvitationLink(os.Getenv("GH_ORG_OPENSOURCE"))
	innersourceLink := OrgInvitationLink(os.Getenv("GH_ORG_INNERSOURCE"))

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

	m := email.EmailMessage{
		Subject: "Github Organizations Invitation",
		Body:    body,
		To:      userEmail,
	}

	email.SendEmail(m)
	fmt.Printf("GitHub Organization Invitation on %s was sent.", time.Now())
}

func OrgInvitationLink(org string) string {
	url := fmt.Sprintf("https://github.com/orgs/%s/invitation", org)
	return fmt.Sprintf("<a href='%s'>%s</a>", url, org)
}

func DateEqual(date time.Time) bool {
	y1, m1, d1 := time.Now().Date()
	y2, m2, d2 := date.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
