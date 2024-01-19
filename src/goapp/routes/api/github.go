package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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
	organization := os.Getenv("GH_ORG_INNERSOURCE")
	emailSupport := os.Getenv("EMAIL_SUPPORT")
	var convertedOutsideCollabsList []string
	users := ghAPI.OrgListMembers(token, organization)
	for _, user := range users {
		email, err := db.GetUserEmailByGithubId(fmt.Sprint(user.GetID()))
		if err != nil {
			logger.LogException(err)
			continue
		}
		if email != "" {
			isUserExist, err := msgraph.IsUserExist(email)
			if err != nil {
				logger.LogException(err)
				continue
			}
			if !isUserExist {
				logger.TrackTrace(fmt.Sprint("GitHub ID: ", user.GetID(), " not found on AD"), contracts.Information)
				// ghAPI.RemoveOrganizationsMember(token, organization, *user.Login)
			}
		} else {
			logger.TrackTrace(fmt.Sprint("GitHub ID: ", user.GetID(), " not found"), contracts.Information)
			// ghAPI.RemoveOrganizationsMember(token, organization, *user.Login)
		}

	}

	// Convert users who are not employees to an outside collaborator
	organizationsOpen := os.Getenv("GH_ORG_OPENSOURCE")

	usersOpenOrg := ghAPI.OrgListMembers(token, organizationsOpen)
	for _, user := range usersOpenOrg {
		email, err := db.GetUserEmailByGithubId(fmt.Sprint(user.GetID()))
		if err != nil {
			logger.LogException(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if email != "" {
			isUserExist, _ := msgraph.IsUserExist(email)
			if !isUserExist {
				logger.TrackTrace(fmt.Sprint("GitHub ID: ", user.GetID(), " not found on AD"), contracts.Information)
				// ghAPI.ConvertMemberToOutsideCollaborator(token, organizationsOpen, *user.Login)
				convertedOutsideCollabsList = append(convertedOutsideCollabsList, email)
			}
		} else {
			logger.TrackTrace(fmt.Sprint("GitHub ID: ", user.GetID(), " not found"), contracts.Information)
			// ghAPI.ConvertMemberToOutsideCollaborator(token, organizationsOpen, *user.Login)
			convertedOutsideCollabsList = append(convertedOutsideCollabsList, email)
		}
	}

	if len(convertedOutsideCollabsList) > 0 {
		// to list of new outside collaborators to ospo
		// EmailAdminConvertToColaborator(emailSupport, convertedOutsideCollabsList, logger)
		emailConvertedCollaboratorTC := appinsights.NewTraceTelemetry(fmt.Sprintf("COLLAB EMAIL : %s", emailSupport), contracts.Information)

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

			for _, collab := range repoAdmins {
				collabEmail, _ := db.GetUserEmailByGithubId(fmt.Sprint(collab.GetID()))

				if len(convertedInRepo) > 0 {
					// EmailRepoAdminConvertToColaborator(collabEmail, repo.Name, convertedInRepo, logger)
					emailAdminConvertedCollaboratorTC := appinsights.NewTraceTelemetry(fmt.Sprintf("COLLAB EMAIL : %s", collabEmail), contracts.Information)

					convertInRepoJson, err := json.Marshal(convertedInRepo)
					if err != nil {
						fmt.Println(err)
						return
					}

					emailAdminConvertedCollaboratorTC.Properties["RepoName"] = repo.Name
					emailAdminConvertedCollaboratorTC.Properties["ConvertedOutsideCollabsList"] = string(convertInRepoJson)
					logger.Track(emailAdminConvertedCollaboratorTC)
				}
			}

		}
	}
}

func RepoOwnerScan(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	organizationsOpen := [...]string{os.Getenv("GH_ORG_OPENSOURCE"), os.Getenv("GH_ORG_INNERSOURCE")}
	var repoOnwerDeficient []string
	var email string
	emailSupport := os.Getenv("EMAIL_SUPPORT")
	for _, org := range organizationsOpen {

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
func emailAdmin(admin string, adminemail string, reponame string, outisideCollab []string, logger *appinsights_wrapper.TelemetryClient) {
	e := time.Now()

	link := "https://github.com/" + os.Getenv("GH_ORG_OPENSOURCE") + "/" + reponame
	link = "<a href=\"" + link + "\">" + reponame + "</a>"
	collabList := "</p> <table  >"
	for _, collab := range outisideCollab {
		collabList = collabList + " <tr> <td>" + collab + " </td></tr>"
	}
	collabList = collabList + " </table  > <p>"
	body := fmt.Sprintf("<p>Hello %s ,  </p>  \n<p>This is to inform you that your Github repository <b> %s </b> has %d outside collaborator/s. </p> %s  This email was sent to the admins of the repository.  </p> \n <p>OSPO</p>", admin, link, len(outisideCollab), collabList)

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
