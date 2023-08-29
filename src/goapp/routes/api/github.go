package routes

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"main/pkg/email"
	db "main/pkg/ghmgmtdb"
	ghAPI "main/pkg/github"
	"main/pkg/msgraph"
	"main/pkg/notification"
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
	org := os.Getenv("GH_ORG_OPENSOURCE")
	var outsideCollabsUsers []string
	token := os.Getenv("GH_TOKEN")
	repos, err := ghAPI.GetRepositoriesFromOrganization(org)
	if err != nil {
		log.Println(err.Error())
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
					log.Println(err.Error())
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				EmailAdmin(admin, email, collab.Name, repoOutsideCollabsList)
			}

		}

	}
}

func ClearOrgMembers(w http.ResponseWriter, r *http.Request) {
	token := os.Getenv("GH_TOKEN")

	// Remove GitHub users from innersource who are not employees
	organization := os.Getenv("GH_ORG_INNERSOURCE")
	emailSupport := os.Getenv("EMAIL_SUPPORT")
	var convertedOutsideCollabsList []string
	users := ghAPI.OrgListMembers(token, organization)
	for _, list := range users {
		email, err := db.UsersGetEmail(*list.Login)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(email) > 0 {
			activeUser, err := msgraph.ActiveUsers(email)
			if err != nil {
				log.Println(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if activeUser == nil {
				ghAPI.RemoveOrganizationsMember(token, organization, *list.Login)

			}
		} else {
			ghAPI.RemoveOrganizationsMember(token, organization, *list.Login)

		}

	}

	// Convert users who are not employees to an outside collaborator
	organizationsOpen := os.Getenv("GH_ORG_OPENSOURCE")

	usersOpenOrg := ghAPI.OrgListMembers(token, organizationsOpen)
	for _, list := range usersOpenOrg {
		email, err := db.UsersGetEmail(*list.Login)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(email) > 0 {
			activeUser, _ := msgraph.ActiveUsers(email)
			if activeUser == nil {
				ghAPI.ConvertMemberToOutsideCollaborator(token, organizationsOpen, *list.Login)
				convertedOutsideCollabsList = append(convertedOutsideCollabsList, *list.Login)
			}
		} else {
			ghAPI.ConvertMemberToOutsideCollaborator(token, organizationsOpen, *list.Login)
			convertedOutsideCollabsList = append(convertedOutsideCollabsList, *list.Login)
		}
	}

	if len(convertedOutsideCollabsList) > 0 {
		// to list of new outside collaborators to ospo
		EmailAdminConvertToColaborator(emailSupport, convertedOutsideCollabsList)

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
				collabEmail, _ := db.UsersGetEmail(*collab.Login)

				if len(convertedInRepo) > 0 {
					EmailRepoAdminConvertToColaborator(collabEmail, repo.Name, convertedInRepo)
				}
			}

		}
	}
}

func RepoOwnerScan(w http.ResponseWriter, r *http.Request) {
	fmt.Println("REPOOWNERSSCAN TRIGGERED...")
	organizationsOpen := [...]string{os.Getenv("GH_ORG_OPENSOURCE"), os.Getenv("GH_ORG_INNERSOURCE")}
	var repoOnwerDeficient []string
	var email string
	emailSupport := os.Getenv("EMAIL_SUPPORT")
	for _, org := range organizationsOpen {

		repoOnwerDeficient = nil
		repos, err := ghAPI.GetRepositoriesFromOrganization(org)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, repo := range repos {
			fmt.Println("Checking number of owners of " + repo.Name)
			owners := GetRepoCollaborators(org, repo.Name, "admin", "direct")
			if len(owners) < 2 {
				fmt.Println(repo.Name + " has less than 2 owners")
				repoOnwerDeficient = append(repoOnwerDeficient, repo.Name)
				for _, owner := range owners {
					email, err = db.UsersGetEmail(*owner.Login)
					if err != nil {
						log.Println(err.Error())
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
			EmailOspoOwnerDeficient(emailSupport, org, repoOnwerDeficient)
		}
	}
	fmt.Println("REPOOWNERSSCAN SUCCESSFUL")
}

func ExpiringInvitation(w http.ResponseWriter, r *http.Request) {
	token := os.Getenv("GH_TOKEN")
	innersourceName := os.Getenv("GH_ORG_INNERSOURCE")
	opensourceName := os.Getenv("GH_ORG_OPENSOURCE")

	sendNotification(token, innersourceName)
	sendNotification(token, opensourceName)
}

// Send notifications to those who has pending org invitation that is about to expire tom.
func sendNotification(token, org string) {
	invitations := ghAPI.ListPendingOrgInvitations(token, org)
	for _, v := range invitations {
		expiresIn, _ := time.ParseDuration("144h")

		if v.CreatedAt.Add(expiresIn).Before(time.Now()) {
			user, err := db.GetUserByGitHubUsername(fmt.Sprint(v.GetLogin()))
			if err != nil {
				log.Println(err.Error())
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
				log.Println(err.Error())
			}
		}
	}
}

// Repo Collaborators Scan
func EmailAdmin(admin string, adminemail string, reponame string, outisideCollab []string) {
	e := time.Now()

	link := "https://github.com/" + os.Getenv("GH_ORG_OPENSOURCE") + "/" + reponame
	link = "<a href=\"" + link + "\">" + reponame + "</a>"
	collabList := "</p> <table  >"
	for _, collab := range outisideCollab {
		collabList = collabList + " <tr> <td>" + collab + " </td></tr>"
	}
	collabList = collabList + " </table  > <p>"
	body := fmt.Sprintf("<p>Hello %s ,  </p>  \n<p>This is to inform you that your Github repository <b> %s </b> has %d outside collaborator/s. </p> %s  This email was sent to the admins of the repository.  </p> \n <p>OSPO</p>", admin, link, len(outisideCollab), collabList)

	m := email.EmailMessage{
		Subject: "GitHub Repo Collaborators Scan",
		Body:    body,
		To:      adminemail,
	}

	email.SendEmail(m)
	fmt.Printf(" GitHub Repo Collaborators Scan on %s was sent.", e)
}

// Deleted repositories from index
func EmailAdminDeletedProjects(to string, repos []string) {
	repoList := "</p> <table  >"
	for _, repo := range repos {
		repoList = repoList + " <tr> <td>" + repo + " </td></tr>"
	}
	repoList = repoList + " </table  > <p>"

	body := fmt.Sprintf("The following repositories were removed from the database as they no longer exist on %s and %s GitHub organizations: %s", os.Getenv("GH_ORG_INNERSOURCE"), os.Getenv("GH_ORG_OPENSOURCE"), repoList)

	m := email.EmailMessage{
		Subject: "List of Deleted Repo",
		Body:    body,
		To:      to,
	}

	email.SendEmail(m)
}

// List of users converted into outside collaborators to Repo Owner
func EmailAdminConvertToColaborator(to string, outisideCollab []string) {
	e := time.Now()
	var body string
	collabList := "</p> <table  >"
	for _, collab := range outisideCollab {
		collabList = collabList + " <tr> <td>" + collab + " </td></tr>"
	}
	collabList = collabList + " </table  > <p>"
	if len(outisideCollab) == 1 {
		body = fmt.Sprintf("<p>Hello %s ,  </p>  \n<p>This is to inform you that %d GitHub user on %s was converted as an outside collaborator. </p> %s  ", to, len(outisideCollab), os.Getenv("GH_ORG_OPENSOURCE"), collabList)
	} else {

		body = fmt.Sprintf("<p>Hello %s ,  </p>  \n<p>This is to inform you that %d GitHub user on %s was converted to an outside collaborator. </p> %s  ", to, len(outisideCollab), os.Getenv("GH_ORG_OPENSOURCE"), collabList)
	}

	m := email.EmailMessage{
		Subject: "GitHub Organization Scan",
		Body:    body,
		To:      to,
	}

	email.SendEmail(m)
	fmt.Printf("GitHub User was converted into an outside  on %s was sent.", e)
}

// List of users converted into outside collaborators to OSPO
func EmailRepoAdminConvertToColaborator(to string, repoName string, outisideCollab []string) {
	e := time.Now()
	var body string
	link := "https://github.com/" + os.Getenv("GH_ORG_OPENSOURCE") + "/" + repoName
	link = "<a href=\"" + link + "\">" + repoName + "</a>"
	collabList := "</p> <table  >"
	for _, collab := range outisideCollab {
		collabList = collabList + " <tr> <td>" + collab + " </td></tr>"
	}

	collabList = collabList + " </table  > <p>"
	if len(outisideCollab) == 1 {
		body = fmt.Sprintf("<p>Hello %s ,  </p>  \n<p>This is to inform you that <b> %d </b> GitHub user on your GitHub repo %s was converted to an outside collaborator. </p> %s This email was sent to the admins of the repository. </p> \n <p>OSPO</p>", to, len(outisideCollab), link, collabList)

	} else {

		body = fmt.Sprintf("<p>Hello %s ,  </p>  \n<p>This is to inform you that <b> %d </b> GitHub users on your GitHub repo %s were converted to outside collaborators. </p> %s This email was sent to the admins of the repository. </p> \n <p>OSPO</p>", to, len(outisideCollab), link, collabList)
	}

	m := email.EmailMessage{
		Subject: "GitHub Organization Scan",
		Body:    body,
		To:      to,
	}

	email.SendEmail(m)
	fmt.Printf("GitHub User was converted into an outside  on %s was sent.", e)
}

// List of repos with less than 2 owners to OSPO
func EmailOspoOwnerDeficient(to string, org string, repoName []string) {
	e := time.Now()
	var body string
	var link string

	repoNameList := "</p> <table  >"
	for _, repo := range repoName {
		link = "https://github.com/" + org + "/" + repo + "/settings/access"
		link = "<a href=\"" + link + "\">" + repo + "</a>"
		repoNameList = repoNameList + " <tr> <td>" + link + " </td></tr>"
	}

	repoNameList = repoNameList + " </table  > <p>"
	if len(repoName) == 1 {
		body = fmt.Sprintf("<p>Hello %s ,  </p>  \n<p>This is to inform you that <b> %d </b> repository on %s needs to add a co-owner.</p> %s   </p>  ", to, len(repoName), org, repoNameList)

	} else {
		body = fmt.Sprintf("<p>Hello %s ,  </p>  \n<p>This is to inform you that <b> %d </b> repositories on %s need to add a co-owner.</p> %s   </p>  ", to, len(repoName), org, repoNameList)
	}
	m := email.EmailMessage{
		Subject: "Repository Owners Scan",
		Body:    body,
		To:      to,
	}

	email.SendEmail(m)
	fmt.Printf(" less than 2 owner    %s was sent.", e)
}

// List of repos with less than 2 owners to repo owner
func EmailcoownerDeficient(to string, Org string, reponame string) {
	var body string
	var link string
	link = "https://github.com/" + Org + "/" + reponame + "/settings/access"
	link = "<a href=\"" + link + "\"> here </a>"

	body = fmt.Sprintf("<p>Hello %s ,  </p>  \n<p>This is to inform you that you are the only admin on %s  GitHub repository. We recommend at least 2 admins on each repository. Click %s to add a co-owner.</p> \n <p>OSPO</p>", to, reponame, link)

	m := email.EmailMessage{
		Subject: "Repository Owners Scan",
		Body:    body,
		To:      to,
	}

	email.SendEmail(m)
}
