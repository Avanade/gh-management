package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"main/pkg/appinsights_wrapper"
	"main/pkg/email"
	db "main/pkg/ghmgmtdb"
	ghAPI "main/pkg/github"
	"main/pkg/msgraph"
	"main/pkg/notification"

	"github.com/microsoft/ApplicationInsights-Go/appinsights"
	"github.com/microsoft/ApplicationInsights-Go/appinsights/contracts"
)

func CheckAvaInnerSource(w http.ResponseWriter, r *http.Request) {

	org := os.Getenv("GH_ORG_INNERSOURCE")
	token := os.Getenv("GH_TOKEN")

	collabs := ghAPI.ListOutsideCollaborators(token, org)
	for _, collab := range collabs {
		ghAPI.RemoveOutsideCollaborator(token, org, *collab.Login)
	}
}

func CheckAvaOpenSource(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	org := os.Getenv("GH_ORG_OPENSOURCE")
	var outsideCollabsUsers []string
	token := os.Getenv("GH_TOKEN")
	repos, err := ghAPI.GetRepositoriesFromOrganization(org)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	outsidecollabs := ghAPI.ListOutsideCollaborators(token, org)
	for _, list := range outsidecollabs {
		outsideCollabsUsers = append(outsideCollabsUsers, *list.Login)
	}
	var repoOutsideCollabsList []string
	for _, collab := range repos {
		var repoCollabsUserNames []string

		var adminmember []string
		repoOutsideCollabsList = nil

		repoCollabs := ghAPI.RepositoriesListCollaborators(token, org, collab.Name, "", "direct")
		for _, list := range repoCollabs {

			repoCollabsUserNames = append(repoCollabsUserNames, *list.Login)
			if *list.RoleName == "admin" {
				adminmember = append(adminmember, *list.Login)

			}
		}

		for _, list := range repoCollabsUserNames {
			for _, outsidelist := range outsideCollabsUsers {
				if list == outsidelist {
					repoOutsideCollabsList = append(repoOutsideCollabsList, outsidelist)
				}
			}
		}
		if len(repoOutsideCollabsList) > 0 {

			for _, admin := range adminmember {
				email, err := db.UsersGetEmail(admin)
				if err != nil {
					logger.LogException(err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				if email != "" {
					emailAdmin(admin, email, collab.Name, repoOutsideCollabsList, logger)
				}
			}

		}

	}
}

func ClearOrgMembers(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	token := os.Getenv("GH_TOKEN")

	// Remove GitHub users from innersource who are not employees
	innersourceOrgs := []string{os.Getenv("GH_ORG_INNERSOURCE")}

	regOrgs, err := db.GetAllRegionalOrganizations()
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, regOrg := range regOrgs {
		innersourceOrgs = append(innersourceOrgs, regOrg["Name"].(string))
	}

	for _, innersourceOrg := range innersourceOrgs {
		ClearOrgMembersInnersource(token, innersourceOrg, logger)
	}

	// Convert users who are not employees to an outside collaborator
	var convertedOutsideCollabsList []string
	organizationsOpen := os.Getenv("GH_ORG_OPENSOURCE")

	usersOpenOrg, err := ghAPI.OrgListMembers(token, organizationsOpen, "all")
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, user := range usersOpenOrg {
		email, err := db.GetUserEmailByGithubId(fmt.Sprint(user.GetID()))
		if err != nil {
			logger.LogException(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if email != "" {
			isUserExist, isAccountEnabled, err := msgraph.IsUserExist(email)
			if err != nil {
				logger.LogException(err)
				continue
			}
			if !isUserExist {
				logger.TrackTrace(fmt.Sprint("GitHub ID: ", user.GetID(), " not found on AD | EXTERNAL"), contracts.Information)
				// ghAPI.ConvertMemberToOutsideCollaborator(token, organizationsOpen, *user.Login)
				convertedOutsideCollabsList = append(convertedOutsideCollabsList, *user.Login)
			}
			if !isAccountEnabled {
				logger.TrackTrace(fmt.Sprint("GitHub ID: ", user.GetID(), " found on AD but account disabled | EXTERNAL"), contracts.Information)
			}
		} else {
			logger.TrackTrace(fmt.Sprint("GitHub ID: ", user.GetID(), " not found | EXTERNAL"), contracts.Information)
			// ghAPI.ConvertMemberToOutsideCollaborator(token, organizationsOpen, *user.Login)
			convertedOutsideCollabsList = append(convertedOutsideCollabsList, *user.Login)
		}
	}

	if len(convertedOutsideCollabsList) > 0 {
		emailSupport := os.Getenv("EMAIL_SUPPORT")
		// to list of new outside collaborators to ospo
		// EmailAdminConvertToColaborator(emailSupport, convertedOutsideCollabsList, logger)
		emailConvertedCollaboratorTC := appinsights.NewTraceTelemetry(fmt.Sprintf("SUPPORT EMAIL : %s", emailSupport), contracts.Information)

		convertedOutsideCollabsListJson, err := json.Marshal(convertedOutsideCollabsList)
		if err != nil {
			fmt.Println(err)
			return
		}

		emailConvertedCollaboratorTC.Properties["ConvertedOutsideCollabsList"] = string(convertedOutsideCollabsListJson)
		logger.Track(emailConvertedCollaboratorTC)

		// to repo admins with converted users
		repos, _ := ghAPI.GetRepositoriesFromOrganization(organizationsOpen)
		for _, repo := range repos {

			repoAdmins := GetRepoCollaborators(organizationsOpen, repo.Name, "admin", "direct")
			repoCollabs := GetRepoCollaborators(organizationsOpen, repo.Name, "", "direct")
			var convertedInRepo []string
			for _, convertedOutsideCollab := range convertedOutsideCollabsList {
				for _, repoCollab := range repoCollabs {
					if convertedOutsideCollab == *repoCollab.Login {
						convertedInRepo = append(convertedInRepo, convertedOutsideCollab)
					}
				}
			}

			if len(convertedInRepo) > 0 {
				for _, collab := range repoAdmins {
					collabEmail, _ := db.GetUserEmailByGithubId(fmt.Sprint(collab.GetID()))

					if collabEmail != "" {
						// EmailRepoAdminConvertToColaborator(collabEmail, repo.Name, convertedInRepo, logger)
						emailAdminConvertedCollaboratorTC := appinsights.NewTraceTelemetry(fmt.Sprintf("ADMIN EMAIL : %s", collabEmail), contracts.Information)

						convertInRepoJson, err := json.Marshal(convertedInRepo)
						if err != nil {
							fmt.Println(err)
							return
						}

						emailAdminConvertedCollaboratorTC.Properties["RepoName"] = repo.Name
						emailAdminConvertedCollaboratorTC.Properties["ConvertedInRepo"] = string(convertInRepoJson)
						logger.Track(emailAdminConvertedCollaboratorTC)
					}
				}
			}

		}
	}
}

func RepoOwnerScan(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	orgs := []string{os.Getenv("GH_ORG_OPENSOURCE"), os.Getenv("GH_ORG_INNERSOURCE")}

	regOrgs, err := db.GetAllRegionalOrganizations()
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, regOrg := range regOrgs {
		orgs = append(orgs, regOrg["Name"].(string))
	}

	var repoOnwerDeficient []string
	var email string
	emailSupport := os.Getenv("EMAIL_SUPPORT")
	for _, org := range orgs {

		repoOnwerDeficient = nil
		repos, err := ghAPI.GetRepositoriesFromOrganization(org)
		if err != nil {
			logger.LogException(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, repo := range repos {
			logger.TrackTrace("Checking number of owners of "+repo.Name, contracts.Information)
			owners := GetRepoCollaborators(org, repo.Name, "admin", "direct")
			if len(owners) < 2 {
				logger.TrackTrace(repo.Name+" has less than 2 owners", contracts.Information)
				repoOnwerDeficient = append(repoOnwerDeficient, repo.Name)
				for _, owner := range owners {
					email, err = db.UsersGetEmail(*owner.Login)
					if err != nil {
						logger.LogException(err)
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}

					if email != "" {
						EmailcoownerDeficient(email, org, repo.Name)
					}
				}
			}
		}

		if len(repoOnwerDeficient) > 0 {
			EmailOspoOwnerDeficient(emailSupport, org, repoOnwerDeficient, logger)
		}
	}
}

func ExpiringInvitation(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	token := os.Getenv("GH_TOKEN")
	innersourceName := os.Getenv("GH_ORG_INNERSOURCE")
	opensourceName := os.Getenv("GH_ORG_OPENSOURCE")

	sendNotification(token, innersourceName, logger)
	sendNotification(token, opensourceName, logger)
}

// Send notifications to those who has pending org invitation that is about to expire tom.
func sendNotification(token, org string, logger *appinsights_wrapper.TelemetryClient) {
	invitations := ghAPI.ListPendingOrgInvitations(token, org)
	for _, v := range invitations {
		expiresIn, _ := time.ParseDuration("144h")

		if v.CreatedAt.Add(expiresIn).Before(time.Now()) {
			user, err := db.GetUserByGitHubUsername(fmt.Sprint(v.GetLogin()))
			if err != nil {
				logger.LogException(err)
			}
			messageBody := notification.OrganizationInvitationExpireMessageBody{
				Recipients: []string{
					user[0]["UserPrincipalName"].(string),
				},
				InvitationLink:   fmt.Sprintf("https://github.com/orgs/%s/invitation", org),
				OrganizationLink: fmt.Sprintf("https://github.com/%s", org),
				OrganizationName: org,
			}
			err = messageBody.Send()
			if err != nil {
				logger.LogException(err)
			}
		}
	}
}

// Repo Collaborators Scan
func emailAdmin(admin string, adminemail string, reponame string, outsideCollab []string, logger *appinsights_wrapper.TelemetryClient) {
	e := time.Now()

	link := "https://github.com/" + os.Getenv("GH_ORG_OPENSOURCE") + "/" + reponame
	link = "<a href=\"" + link + "\">" + reponame + "</a>"
	collabList := "</p> <table  >"
	for _, collab := range outsideCollab {
		collabList = collabList + " <tr> <td>" + collab + " </td></tr>"
	}
	collabList = collabList + " </table  > <p>"

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
						<td class="center-table" align="center">
							<table style="width: 100%; max-width: 700px;" class="margin-auto">
								<tr>
									<td style="padding-top: 20px">
										<p>Hello |Admin| ,  </p>
										<p>This is to inform you that your Github repository <b> |Link| </b> has |NumberOfOutsideCollaborators| outside collaborator/s. </p>
										|CollabList|
										<p>This email was sent to the admins of the repository.  </p>
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
		"|Admin|", admin,
		"|Link|", link,
		"|NumberOfOutsideCollaborators|", strconv.Itoa(len(outsideCollab)),
		"|CollabList|", collabList,
	)
	body := replacer.Replace(bodyTemplate)

	m := email.Message{
		Subject: "GitHub Repo Collaborators Scan",
		Body: email.Body{
			Content: body,
			Type:    email.HtmlMessageType,
		},
		ToRecipients: []email.Recipient{
			{
				Email: adminemail,
			},
		},
	}

	err := email.SendEmail(m, true)
	if err != nil {
		logger.LogException(err)
		return
	}
	logger.TrackTrace(fmt.Sprintf(" GitHub Repo Collaborators Scan on %s was sent.", e), contracts.Information)
}

type RemovedMember struct {
	Id          int64  `json:"id"`
	Username    string `json:"username"`
	Information string `json:"information"`
}

func ClearOrgMembersInnersource(token, org string, logger *appinsights_wrapper.TelemetryClient) {
	var removedMembers []RemovedMember

	users, _ := ghAPI.OrgListMembers(token, org, "all")
	for _, user := range users {
		email, err := db.GetUserEmailByGithubId(fmt.Sprint(user.GetID()))
		if err != nil {
			logger.LogException(err)
			continue
		}
		if email != "" {
			isUserExist, isAccountEnabled, err := msgraph.IsUserExist(email)
			if err != nil {
				logger.LogException(err)
				continue
			}
			if !isUserExist {
				information := fmt.Sprint("GitHub ID: ", user.GetID(), " not found on AD | INTERNAL")
				removedMembers = append(removedMembers, RemovedMember{
					Id:          user.GetID(),
					Username:    user.GetLogin(),
					Information: information,
				})
				logger.TrackTrace(information, contracts.Information)
				// ghAPI.RemoveOrganizationsMember(token, organization, *user.Login)
			}
			if !isAccountEnabled {
				information := fmt.Sprint("GitHub ID: ", user.GetID(), " found on AD but account disabled | INTERNAL")
				removedMembers = append(removedMembers, RemovedMember{
					Id:          user.GetID(),
					Username:    user.GetLogin(),
					Information: information,
				})
				logger.TrackTrace(information, contracts.Information)
			}
		} else {
			information := fmt.Sprint("GitHub ID: ", user.GetID(), " not found | INTERNAL")
			removedMembers = append(removedMembers, RemovedMember{
				Id:          user.GetID(),
				Username:    user.GetLogin(),
				Information: information,
			})
			logger.TrackTrace(information, contracts.Information)
			// ghAPI.RemoveOrganizationsMember(token, organization, *user.Login)
		}
	}
	removedMembersTC := appinsights.NewTraceTelemetry(fmt.Sprintf("INNERSOURCE ORG : %s", org), contracts.Information)

	removedMembersJson, err := json.Marshal(removedMembers)
	if err != nil {
		fmt.Println(err)
		return
	}

	removedMembersTC.Properties["RemovedMembers"] = string(removedMembersJson)
	logger.Track(removedMembersTC)
}
