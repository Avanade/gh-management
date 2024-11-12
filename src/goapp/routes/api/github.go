package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"main/pkg/appinsights_wrapper"
	"main/pkg/email"
	ev "main/pkg/envvar"
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
	token := os.Getenv("GH_TOKEN")
	repos, err := ghAPI.GetRepositoriesFromOrganization(org)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	orgOutsideCollaborators := ghAPI.ListOutsideCollaborators(token, org)

	var orgOutsideCollaboratorUsernames []string

	for _, orgOutsideCollaborator := range orgOutsideCollaborators {
		orgOutsideCollaboratorUsernames = append(orgOutsideCollaboratorUsernames, *orgOutsideCollaborator.Login)
	}

	for _, repo := range repos {
		var repoCollabUsernames []string
		var repoAdminUsernames []string

		repoCollaborators := ghAPI.RepositoriesListCollaborators(token, org, repo.Name, "", "direct")
		for _, repoCollaborator := range repoCollaborators {

			repoCollabUsernames = append(repoCollabUsernames, *repoCollaborator.Login)
			if *repoCollaborator.RoleName == "admin" {
				repoAdminUsernames = append(repoAdminUsernames, *repoCollaborator.Login)
			}
		}

		var repoOutsideCollaboratorUsernames []string
		for _, repoCollabUsername := range repoCollabUsernames {
			for _, orgOutsideCollaboratorUsername := range orgOutsideCollaboratorUsernames {
				if repoCollabUsername == orgOutsideCollaboratorUsername {
					repoOutsideCollaboratorUsernames = append(repoOutsideCollaboratorUsernames, orgOutsideCollaboratorUsername)
				}
			}
		}

		if len(repoOutsideCollaboratorUsernames) > 0 {
			logger.TrackTrace(fmt.Sprintf("%s has %d outside collaborators", repo.Name, len(repoOutsideCollaboratorUsernames)), contracts.Information)

			for _, repoAdminUsername := range repoAdminUsernames {
				repoAdminEmail, err := db.UsersGetEmail(repoAdminUsername)
				if err != nil {
					continue
				}

				if repoAdminEmail != "" {
					emailAdmin(repoAdminUsername, repoAdminEmail, repo.Name, repoOutsideCollaboratorUsernames, logger)
				}
			}

		}
	}
}

func ClearOrgMembers(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	start := time.Now()
	// FETCH ENTERPRISE
	enterpriseToken := os.Getenv("GH_ENTERPRISE_TOKEN")
	enterpriseName := os.Getenv("GH_ENTERPRISE_NAME")
	enterpriseMembers, err := ghAPI.GetMembersByEnterprise(enterpriseName, enterpriseToken)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.LogTrace(fmt.Sprintf("Fetched %d enterprise members", len(enterpriseMembers.Members)), contracts.Information)

	// PROCESS INNERSOURCE & REGIONAL ORGS
	err = ProcessCleanupEnterpriseOrgs(enterpriseMembers, logger)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// PROCESS OPENSOURCE ORG
	err = ProcessCleanupOpensourceOrg(enterpriseMembers, logger)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Time elapsed: ", time.Since(start))
	w.WriteHeader(http.StatusAccepted)
}

func RepoOwnerScan(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	orgs := []string{os.Getenv("GH_ORG_OPENSOURCE"), os.Getenv("GH_ORG_INNERSOURCE")}

	isEnabled := db.NullBool{Value: true}
	regOrgs, err := db.SelectRegionalOrganization(&isEnabled)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, regOrg := range regOrgs {
		if !regOrg.IsIndexRepoEnabled {
			continue
		}
		orgs = append(orgs, regOrg.Name)
	}

	// Temporarily log Organization
	if len(regOrgs) > 0 {
		utilRepoOwnerScan := appinsights.NewTraceTelemetry("util repo owner scan", contracts.Information)
		regOrgsJson, err := json.Marshal(regOrgs)
		if err != nil {
			fmt.Println(err)
			return
		}
		utilRepoOwnerScan.Properties["Orgs"] = string(regOrgsJson)
		logger.Track(utilRepoOwnerScan)

		if ev.GetEnvVar("ENABLED_REPO_OWNER_SCAN", "false") != "true" {
			return
		}
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
	Id       int64
	Username string
}

func ClearOrgMembersInnersource(token, org string, logger *appinsights_wrapper.TelemetryClient) {
	var notFoundDB []string
	var notFoundAD []string
	var disabledAccountAD []string

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
				notFoundAD = append(notFoundAD, fmt.Sprint(user.GetLogin(), " - ", email))
				if ev.GetEnvVar("ENABLED_REMOVE_COLLABORATORS", "false") == "true" {
					ghAPI.RemoveOrganizationsMember(token, org, user.GetLogin()) // Remove user from organization
				}
			}
			if !isAccountEnabled {
				disabledAccountAD = append(disabledAccountAD, fmt.Sprint(user.GetLogin(), " - ", email))
				if ev.GetEnvVar("ENABLED_REMOVE_COLLABORATORS", "false") == "true" {
					ghAPI.RemoveOrganizationsMember(token, org, user.GetLogin()) // Remove user from organization
				}
			}
		} else {
			notFoundDB = append(notFoundDB, user.GetLogin())
			if ev.GetEnvVar("ENABLED_REMOVE_COLLABORATORS", "false") == "true" {
				ghAPI.RemoveOrganizationsMember(token, org, user.GetLogin()) // Remove user from organization
			}
		}
	}
	removedMembersTC := appinsights.NewTraceTelemetry(fmt.Sprintf("INNERSOURCE ORG : %s", org), contracts.Information)

	removedMembersTC.Properties["NotFoundOnAD"] = strings.Join(notFoundAD, ",")
	removedMembersTC.Properties["NotFoundOnDB"] = strings.Join(notFoundDB, ",")
	removedMembersTC.Properties["DisabledADAccount"] = strings.Join(disabledAccountAD, ",")
	logger.Track(removedMembersTC)
}

type MemberWithCommunityOrgs struct {
	Id       int64
	Username string
	Email    string
	Orgs     []string
}

// Process Innsersource & Regional Orgs
func ProcessCleanupEnterpriseOrgs(enterpriseMembers *ghAPI.GetMembersByEnterpriseResult, logger *appinsights_wrapper.TelemetryClient) error {
	removeMemberTC := appinsights.NewTraceTelemetry("Cleanup Enterprise Orgs (Innersource & Regional Orgs)", contracts.Information)

	// Fetch enterprise orgs
	enterpriseOrgs := []string{os.Getenv("GH_ORG_INNERSOURCE")}
	isEnabled := db.NullBool{Value: true}
	regOrgs, err := db.SelectRegionalOrganization(&isEnabled)
	if err != nil {
		logger.LogException(err)
		return err
	}
	removeMemberTC.Properties["1st Step"] = fmt.Sprint("Fetched regional organizations. regOrgs.Length : ", len(regOrgs))

	for _, regOrg := range regOrgs {
		enterpriseOrgs = append(enterpriseOrgs, regOrg.Name)
	}
	removeMemberTC.Properties["2nd Step"] = fmt.Sprint("Filter enterprise organizations. enterpriseOrgs.Length : ", len(enterpriseOrgs))

	// Filter enterprise members if they are in the community organizations
	var communityMembers []Member
	communityMembersSet := make(map[string]struct{})

	var wgFCM sync.WaitGroup
	var muFCM sync.Mutex
	concurrencyLimitFCM := make(chan struct{}, 50) // Limit to 50 concurrent goroutines

	for _, enterpriseOrg := range enterpriseOrgs {
		wgFCM.Add(1)
		concurrencyLimitFCM <- struct{}{} // Acquire a slot

		go func(enterpriseOrg string) {
			defer wgFCM.Done()
			defer func() { <-concurrencyLimitFCM }() // Release the slot

			token := os.Getenv("GH_TOKEN")
			// Fetch all members of the community organization
			members, err := ghAPI.OrgListMembers(token, enterpriseOrg, "all")
			if err != nil {
				logger.LogException(err)
				return
			}

			for _, member := range members {
				for _, enterpriseMember := range enterpriseMembers.Members {
					if member.GetLogin() == enterpriseMember.Login {
						// Check if exist in communityMembers
						if _, exists := communityMembersSet[enterpriseMember.Login]; !exists {
							muFCM.Lock()
							communityMembers = append(communityMembers, Member{
								NodeId:     enterpriseMember.Id,
								DatabaseId: enterpriseMember.DatabaseId,
								Username:   enterpriseMember.Login,
								Email:      enterpriseMember.EnterpriseEmail})
							muFCM.Unlock()
							communityMembersSet[member.GetLogin()] = struct{}{}
							break
						}
					}
				}
			}
		}(enterpriseOrg)
	}
	wgFCM.Wait()
	removeMemberTC.Properties["3rd Step"] = fmt.Sprint("Filter github enterprise members if they are a member of any community organizations. communityMembers.Length : ", len(communityMembers))

	// Remove members that are not in the active directory
	var removeMembers []string
	var wgRAD sync.WaitGroup
	var muRAD sync.Mutex
	concurrencyLimitRAD := make(chan struct{}, 50) // Limit to 50 concurrent goroutines
	for _, member := range communityMembers {
		wgRAD.Add(1)
		concurrencyLimitRAD <- struct{}{} // Acquire a slot

		go func(member Member) {
			defer wgRAD.Done()
			defer func() { <-concurrencyLimitRAD }() // Release the slot

			isUserExist, isAccountEnabled, err := msgraph.IsUserExist(member.Email)
			if err != nil {
				logger.LogException(err)
				return
			}
			if !isUserExist || !isAccountEnabled {
				muRAD.Lock()
				removeMembers = append(removeMembers, fmt.Sprintln(member.Username, " - ", member.Email))
				muRAD.Unlock()
				if ev.GetEnvVar("ENABLED_REMOVE_COLLABORATORS", "false") == "true" {
					token := os.Getenv("GH_TOKEN")
					enterpriseId := os.Getenv("GH_ENTERPRISE_ID")
					err := ghAPI.RemoveEnterpriseMember(token, enterpriseId, member.NodeId)
					if err != nil {
						logger.LogException(err)
					}
				}
			}
		}(member)
	}
	wgRAD.Wait()

	removeMemberChunks := splitStringArray(removeMembers, 50)
	for indexRemoveMembersChunk, removeMembersChunk := range removeMemberChunks {
		removeMemberTC.Properties[fmt.Sprint("RemovedMembers_", indexRemoveMembersChunk)] = strings.Join(removeMembersChunk, ",")
	}

	// Log removed members
	logger.Track(removeMemberTC)
	return nil
}

// ***** Convert members to outside collaborators if not exist in AD *****
func ProcessCleanupOpensourceOrg(enterpriseMembers *ghAPI.GetMembersByEnterpriseResult, logger *appinsights_wrapper.TelemetryClient) error {
	convertedToOutsideCollaboratorsTC := appinsights.NewTraceTelemetry("Cleanup Opensource Orgs", contracts.Information)
	opensourceOrg := os.Getenv("GH_ORG_OPENSOURCE")

	// Fetch opensource members
	opensourceMembers, err := ghAPI.OrgListMembers(os.Getenv("GH_TOKEN"), opensourceOrg, "all")
	if err != nil {
		return err
	}
	convertedToOutsideCollaboratorsTC.Properties["1st Step"] = fmt.Sprint("Fetched opensource members. opensourceMembers.Length : ", len(opensourceMembers))

	// Filter opensource members if they are in the opensource organization
	var openSourceMembers []Member
	for _, member := range opensourceMembers {
		for _, enterpriseMember := range enterpriseMembers.Members {
			if member.GetLogin() == enterpriseMember.Login {
				openSourceMembers = append(openSourceMembers, Member{NodeId: enterpriseMember.Id, DatabaseId: enterpriseMember.DatabaseId, Username: enterpriseMember.Login, Email: enterpriseMember.EnterpriseEmail})
			}
		}
	}
	convertedToOutsideCollaboratorsTC.Properties["2nd Step"] = fmt.Sprint("Filter opensource members if they are a member of the opensource organization. openSourceMembers.Length : ", len(openSourceMembers))

	// Remove members that are not in the active directory
	var removeMembers []string
	var wgRAD sync.WaitGroup
	var muRAD sync.Mutex
	concurrencyLimitRAD := make(chan struct{}, 50) // Limit to 50 concurrent goroutines
	for _, member := range openSourceMembers {
		wgRAD.Add(1)
		concurrencyLimitRAD <- struct{}{} // Acquire a slot

		go func(member Member) {
			defer wgRAD.Done()
			defer func() { <-concurrencyLimitRAD }() // Release the slot

			isUserExist, isAccountEnabled, err := msgraph.IsUserExist(member.Email)
			if err != nil {
				logger.LogException(err)
				return
			}
			if !isUserExist || !isAccountEnabled {
				muRAD.Lock()
				removeMembers = append(removeMembers, fmt.Sprintln(member.Username, " - ", member.Email))
				muRAD.Unlock()
				if ev.GetEnvVar("ENABLED_REMOVE_COLLABORATORS", "false") == "true" {
					token := os.Getenv("GH_TOKEN")
					ghAPI.ConvertMemberToOutsideCollaborator(token, opensourceOrg, member.Username)
				}
			}
		}(member)
	}
	wgRAD.Wait()
	convertedToOutsideCollaboratorsTC.Properties["3rd Step"] = fmt.Sprint("Convert members to outside collaborators if not exist in AD. removeMembers.Length : ", len(removeMembers))

	removeMember := splitStringArray(removeMembers, 50)
	for indexRemoveMembersChunk, removeMembersChunk := range removeMember {
		convertedToOutsideCollaboratorsTC.Properties[fmt.Sprint("RemovedUsers_OSO", indexRemoveMembersChunk)] = strings.Join(removeMembersChunk, ",")
	}
	logger.Track(convertedToOutsideCollaboratorsTC)
	return nil
}
