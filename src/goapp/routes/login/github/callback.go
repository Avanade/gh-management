package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"main/pkg/appinsights_wrapper"
	auth "main/pkg/authentication"
	"main/pkg/email"
	ev "main/pkg/envvar"
	db "main/pkg/ghmgmtdb"
	ghAPI "main/pkg/github"
	"main/pkg/msgraph"
	"main/pkg/notification"
	"main/pkg/session"

	"github.com/gorilla/sessions"
	"github.com/microsoft/ApplicationInsights-Go/appinsights/contracts"
	"golang.org/x/oauth2"
)

func GithubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	sessionaz, err := session.Store.Get(r, "auth-session")
	if err != nil {
		log.Println(err.Error())
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		http.Redirect(w, r, "/authentication/github/failed", http.StatusSeeOther)
		return
	}
	// Check session and state
	state, err := session.GetState(w, r)
	if err != nil {
		log.Println(err.Error())
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		http.Redirect(w, r, "/authentication/github/failed", http.StatusSeeOther)
		return
	}

	session, err := session.Store.Get(r, "gh-auth-session")
	if err != nil {
		log.Println(err.Error())
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		http.Redirect(w, r, "/authentication/github/failed", http.StatusSeeOther)
		return
	}

	if r.URL.Query().Get("state") != state {
		log.Println("Invalid state paramerer")
		// http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		http.Redirect(w, r, "/authentication/github/failed", http.StatusSeeOther)
		return
	}

	ghauth := auth.GetGitHubOauthConfig(r.Host)

	// Exchange temporary code for access token
	code := r.URL.Query().Get("code")

	ghAccessToken, err := ghauth.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Println(err.Error())
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		http.Redirect(w, r, "/authentication/github/failed", http.StatusSeeOther)
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
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		http.Redirect(w, r, "/authentication/github/failed", http.StatusSeeOther)
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
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		http.Redirect(w, r, "/authentication/github/failed", http.StatusSeeOther)
		return
	}
	session.Values["ghIsValid"] = result["IsValid"].(bool)
	sessionaz.Values["isGHAssociated"] = result["IsValid"].(bool)
	sessionaz.Save(r, w)

	isDirect, _ := msgraph.IsDirectMember(fmt.Sprintf("%s", azProfile["oid"]))
	isEnterpriseMember, _ := msgraph.IsGithubEnterpriseMember(fmt.Sprintf("%s", azProfile["oid"]))

	session.Values["ghIsDirect"] = isDirect
	session.Values["ghIsEnterpriseMember"] = isEnterpriseMember

	lastGithubLogin := result["LastGithubLogin"].(time.Time)

	if !DateEqual(lastGithubLogin) && result["IsValid"].(bool) {
		CheckMembership(userPrincipalName, ghUser)
	}

	session.Options = &sessions.Options{
		Path:     "/",
		Domain:   "",
		MaxAge:   2592000,
		Secure:   true,
		HttpOnly: false,
		SameSite: http.SameSiteNoneMode,
	}
	err = session.Save(r, w)
	if err != nil {
		log.Panicln(err.Error())
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		http.Redirect(w, r, "/authentication/github/failed", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/authentication/github/successful", http.StatusSeeOther)
}

func GithubForceSaveHandler(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	sessionaz, err := session.Store.Get(r, "auth-session")
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Check session and state
	session, err := session.Store.Get(r, "gh-auth-session")
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ghProfile := session.Values["ghProfile"].(string)

	var p map[string]interface{}
	err = json.Unmarshal([]byte(ghProfile), &p)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Save and Validate github account
	azProfile := sessionaz.Values["profile"].(map[string]interface{})
	userPrincipalName := fmt.Sprintf("%s", azProfile["preferred_username"])
	newGhId := strconv.FormatFloat(p["id"].(float64), 'f', 0, 64)
	newGhUser := fmt.Sprintf("%s", p["login"])

	logger.LogTrace(fmt.Sprintf("User %s is trying to reassociate new GitHub account %s", userPrincipalName, newGhUser), contracts.Information)

	if ev.GetEnvVar("ENABLED_REMOVE_ENTERPRISE_MEMBER", "false") == "true" {
		// Get the current associated GitHub account
		currentDbUser, err := db.GetUserByUserPrincipal(userPrincipalName)
		if err != nil {
			logger.LogException(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Get the user by GitHub ID
		user, err := ghAPI.GetUserByLogin(currentDbUser[0]["GitHubUser"].(string), os.Getenv("GH_TOKEN"))
		if err != nil {
			logger.LogException(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		orgs, err := ghAPI.GetOrganizationsByGitHubName(user.Login, os.Getenv("GH_TOKEN"))
		if err != nil {
			logger.LogException(err)
		}

		if orgs != nil && len(orgs.Organizations) > 0 {
			isMembershipChecked := false
			for _, org := range orgs.Organizations {
				if org.Login == os.Getenv("GH_ORG_OPENSOURCE") || org.Login == os.Getenv("GH_ORG_INNERSOURCE") {
					if !isMembershipChecked {
						logger.LogTrace("Checking membership of user in opensource & innersource orgs", contracts.Information)
						CheckMembership(userPrincipalName, newGhUser)
						isMembershipChecked = true
					}
				} else {
					invite := ghAPI.OrganizationInvitation(os.Getenv("GH_TOKEN"), newGhUser, org.Login)
					if invite == nil {
						logger.LogTrace("Error sending invitation to organization", contracts.Error)
					} else {
						logger.LogTrace("Invitation sent to organization", contracts.Information)
					}
				}
				collaboratorRepos, err := ghAPI.GetCollaboratorRepositoriesFromOrganization(os.Getenv("GH_TOKEN"), org.Login, user.Login)
				if err != nil {
					logger.LogException(err)
				}
				for _, colRepo := range collaboratorRepos {
					permission, err := ghAPI.GetPermissionLevel(org.Login, colRepo.Name, user.Login)
					if err != nil {
						logger.LogException(err)
						continue
					}
					_, err = ghAPI.AddCollaborator(colRepo.Org, colRepo.Name, newGhUser, permission)
					if err != nil {
						logger.LogException(err)
						continue
					}
				}
			}
		} else {
			logger.LogTrace(fmt.Sprintf("No organizations found for user: %s", user.Login), contracts.Information)
		}

		enterpriseToken := os.Getenv("GH_ENTERPRISE_TOKEN")
		enterpriseId := os.Getenv("GH_ENTERPRISE_ID")
		err = ghAPI.RemoveEnterpriseMember(enterpriseToken, enterpriseId, user.Id)
		if err != nil {
			logger.LogException(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	result, err := db.UpdateUserGithub(userPrincipalName, newGhId, newGhUser, 1)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["ghIsValid"] = result["IsValid"].(bool)

	session.Options = &sessions.Options{
		Path:     "/",
		Domain:   "",
		MaxAge:   2592000,
		Secure:   true,
		HttpOnly: false,
		SameSite: http.SameSiteNoneMode,
	}
	err = session.Save(r, w)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	orgNames := struct {
		InnerSourceOrgName string `json:"innersourceOrgName"`
	}{
		InnerSourceOrgName: os.Getenv("GH_ORG_INNERSOURCE"),
	}

	jsonResp, err := json.Marshal(orgNames)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write(jsonResp)
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

	var message string

	if !isInnersourceOrgMember && !isOpensourceOrgMember {
		message = fmt.Sprintf("<p>An invitation to join %s and %s was created. You need to join the %s so you get added to the repository you'll request.</p>", innersourceLink, opensourceLink, innersourceLink)
	} else if !isInnersourceOrgMember {
		message = fmt.Sprintf("<p>An invitation to join %s was created. You need to join this organization so you get added to the repository you'll request.</p>", innersourceLink)
	} else if !isOpensourceOrgMember {
		message = fmt.Sprintf("<p>An invitation to join %s was created. You need to join this organization so you won't get tagged as an outside collaborator.</p>", opensourceLink)
	} else {
		return
	}

	bodyTemplate := `
		<html>
			<head>
				<style>
					table,
					th,
					tr,
					td {
					border: 0;
					border-collapse: collapse;
					vertical-align: middle;
					}

					.thead {
					padding: 15px;
					}

					.center-table {
					text-align: -webkit-center;
					}

					.margin-auto {
					margin: auto;
					}

					.border-top {
					border-top: 1px rgb(204, 204, 204) solid;
					border-collapse: separate;
					}
				</style>
			</head>

			<body>
				<table style="width: 100%">
					<tr>
						<th class="center-table">
							<table style="width: 100%; max-width: 700px;" class="margin-auto">
								<tr style="background-color: #ff5800">
									<td class="thead" style="width: 95px">
									<img
										src="data:image/png;base64,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" />
									</td>
									<td class="thead" style="
										font-family: SegoeUI, sans-serif;
										font-size: 14px;
										color: white;
										padding-bottom: 10px;
										">
									Community
									</td>
								</tr>
							</table>
						</th>
					</tr>
					<tr>
						<th class="center-table">
							<table style="width: 100%; max-width: 700px;" class="margin-auto">
								<tr>
									<td style="font-size: 18px; font-weight: 600; padding-top: 20px">
									Invitation to Join GitHub Organization
									</td>
								</tr>
							</table>
						</th>
					</tr>
					<tr>
						<td class="center-table" align="center">
							<table style="width: 100%; max-width: 700px;" class="margin-auto">
								<tr>
									<td style="padding-top: 20px">
										<p>Hello |Recipient|,  </p>
										|Message|
										<p>OSPO</p>
									</td>
								</tr>
							</table>
						</td>
					</tr>
				</table>
				<br>
			</body>

		</html>
	`

	replacer := strings.NewReplacer(
		"|Recipient|", ghUsername,
		"|Message|", message,
	)

	body := replacer.Replace(bodyTemplate)

	m := email.Message{
		Subject: "Github Organizations Invitation",
		Body: email.Body{
			Content: body,
			Type:    email.HtmlMessageType,
		},
		ToRecipients: []email.Recipient{
			{
				Email: userEmail,
			},
		},
	}

	email.SendEmail(m, true)
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
