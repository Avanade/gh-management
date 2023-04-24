package routes

import (
	githubAPI "main/pkg/github"
	"net/http"
	"os"
	"strings"
	"time"

	"fmt"
	"main/pkg/email"
	db "main/pkg/ghmgmtdb"
	"main/pkg/msgraph"
)

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
				EmailAdmin(admin, email, collab.Name, OutsideRepocollabsList)
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
		EmailAdminConvertToColaborator(EmailSupport, ConvertedOutsidecollabsList)

		// Email repo admins with converted users
		repos, _ := githubAPI.GetRepositoriesFromOrganization(organizationsOpen)
		for _, repo := range repos {

			RepoAdmins := GetRepoAdmin(organizationsOpen, repo.Name)
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
					EmailRepoAdminConvertToColaborator(collabemail, repo.Name, ConvertedInRepo)
				}
			}

		}
	}

}

func RepoOwnerScan(w http.ResponseWriter, r *http.Request) {

	organizationsOpen := [...]string{os.Getenv("GH_ORG_OPENSOURCE"), os.Getenv("GH_ORG_INNERSOURCE")}
	var repoOnwerDeficient []string
	var email string
	EmailSupport := os.Getenv("EMAIL_SUPPORT")
	for _, org := range organizationsOpen {

		repoOnwerDeficient = nil
		repos, _ := githubAPI.GetRepositoriesFromOrganization(org)

		for _, repo := range repos {
			owners := GetRepoAdmin(org, repo.Name)
			if len(owners) < 2 {
				repoOnwerDeficient = append(repoOnwerDeficient, repo.Name)
				for _, owner := range owners {
					email, _ = db.UsersGetEmail(owner)

					EmailcoownerDeficient(email, org, repo.Name)
				}
			}

		}

		if len(repoOnwerDeficient) > 0 {
			EmailOspoOwnerDeficient(EmailSupport, org, repoOnwerDeficient)
		}
	}

}

func EmailAdmin(admin string, adminemail string, reponame string, outisideCollab []string) {
	e := time.Now()

	link := "https://github.com/" + os.Getenv("GH_ORG_OPENSOURCE") + "/" + reponame
	link = "<a href=\"" + link + "\">" + reponame + "</a>"
	Collablist := "</p> <table  >"
	for _, collab := range outisideCollab {
		Collablist = Collablist + " <tr> <td>" + collab + " </td></tr>"
	}
	Collablist = Collablist + " </table  > <p>"
	body := fmt.Sprintf("<p>Hello %s ,  </p>  \n<p>This is to inform you that your Github repository <b> %s </b> has %d outside collaborator/s. </p> %s  This email was sent to the admins of the repository.  </p> \n <p>OSPO</p>", admin, link, len(outisideCollab), Collablist)

	m := email.TypEmailMessage{
		Subject: "GitHub Repo Collaborators Scan",
		Body:    body,
		To:      adminemail,
	}

	email.SendEmail(m)
	fmt.Printf(" GitHub Repo Collaborators Scan on %s was sent.", e)
}

func EmailAdminConvertToColaborator(Email string, outisideCollab []string) {
	e := time.Now()
	var body string
	Collablist := "</p> <table  >"
	for _, collab := range outisideCollab {
		Collablist = Collablist + " <tr> <td>" + collab + " </td></tr>"
	}
	Collablist = Collablist + " </table  > <p>"
	if len(outisideCollab) == 1 {
		body = fmt.Sprintf("<p>Hello %s ,  </p>  \n<p>This is to inform you that %d GitHub user on Avanade was converted as an outside collaborator. </p> %s  ", Email, len(outisideCollab), Collablist)
	} else {

		body = fmt.Sprintf("<p>Hello %s ,  </p>  \n<p>This is to inform you that %d GitHub user on Avanade was converted to an outside collaborator. </p> %s  ", Email, len(outisideCollab), Collablist)
	}

	m := email.TypEmailMessage{
		Subject: "GitHub Organization Scan",
		Body:    body,
		To:      Email,
	}

	email.SendEmail(m)
	fmt.Printf("GitHub User was converted into an outside  on %s was sent.", e)
}

func EmailRepoAdminConvertToColaborator(Email string, reponame string, outisideCollab []string) {
	e := time.Now()
	var body string
	link := "https://github.com/" + os.Getenv("GH_ORG_OPENSOURCE") + "/" + reponame
	link = "<a href=\"" + link + "\">" + reponame + "</a>"
	Collablist := "</p> <table  >"
	for _, collab := range outisideCollab {
		Collablist = Collablist + " <tr> <td>" + collab + " </td></tr>"
	}

	Collablist = Collablist + " </table  > <p>"
	if len(outisideCollab) == 1 {
		body = fmt.Sprintf("<p>Hello %s ,  </p>  \n<p>This is to inform you that <b> %d </b> GitHub user on your GitHub repo %s was converted to an outside collaborator. </p> %s This email was sent to the admins of the repository. </p> \n <p>OSPO</p>", Email, len(outisideCollab), link, Collablist)

	} else {

		body = fmt.Sprintf("<p>Hello %s ,  </p>  \n<p>This is to inform you that <b> %d </b> GitHub users on your GitHub repo %s were converted to outside collaborators. </p> %s This email was sent to the admins of the repository. </p> \n <p>OSPO</p>", Email, len(outisideCollab), link, Collablist)
	}

	m := email.TypEmailMessage{
		Subject: "GitHub Organization Scan",
		Body:    body,
		To:      Email,
	}

	email.SendEmail(m)
	fmt.Printf("GitHub User was converted into an outside  on %s was sent.", e)
}

func GetRepoAdmin(org string, repo string) []string {
	var Adminmember []string
	var RepocollabsList []string
	var OrgOwners []string
	token := os.Getenv("GH_TOKEN")
	ORG_OWNERS := os.Getenv("ORG_OWNERS")
	OrgOwners = strings.Fields(ORG_OWNERS)
	Repocollabs := githubAPI.RepositoriesListCollaborators(token, org, repo)

	for _, list := range Repocollabs {

		RepocollabsList = append(RepocollabsList, *list.Login)
		if *list.RoleName == "admin" {
			if !stringInArray(*list.Login, OrgOwners) {
				Adminmember = append(Adminmember, *list.Login)

			}
		}
	}
	return Adminmember
}

func EmailOspoOwnerDeficient(Email string, org string, reponame []string) {
	e := time.Now()
	var body string
	var link string

	reponamelist := "</p> <table  >"
	for _, repo := range reponame {
		link = "https://github.com/" + org + "/" + repo + "/settings/access"
		link = "<a href=\"" + link + "\">" + repo + "</a>"
		reponamelist = reponamelist + " <tr> <td>" + link + " </td></tr>"
	}

	reponamelist = reponamelist + " </table  > <p>"
	if len(reponame) == 1 {
		body = fmt.Sprintf("<p>Hello %s ,  </p>  \n<p>This is to inform you that <b> %d </b> repository on %s needs to add a co-owner.</p> %s   </p>  ", Email, len(reponame), org, reponamelist)

	} else {
		body = fmt.Sprintf("<p>Hello %s ,  </p>  \n<p>This is to inform you that <b> %d </b> repositories on %s need to add a co-owner.</p> %s   </p>  ", Email, len(reponame), org, reponamelist)
	}
	m := email.TypEmailMessage{
		Subject: "GitHub Organization Scan",
		Body:    body,
		To:      Email,
	}

	email.SendEmail(m)
	fmt.Printf(" less than 2 owner    %s was sent.", e)
}

func EmailcoownerDeficient(Email string, Org string, reponame string) {
	e := time.Now()
	var body string
	var link string
	link = "https://github.com/" + Org + "/" + reponame + "/settings/access"
	link = "<a href=\"" + link + "\"> here </a>"

	body = fmt.Sprintf("<p>Hello %s ,  </p>  \n<p>This is to inform you that you are the only admin on %s  GitHub repository. We recommend at least 2 admins on each repository. Click %s to add a co-owner.</p> \n <p>OSPO</p>", Email, reponame, link)

	m := email.TypEmailMessage{
		Subject: "GitHub Organization Scan",
		Body:    body,
		To:      Email,
	}

	email.SendEmail(m)
	fmt.Printf(" less than 2 owner  on %s was sent.", e)
}
func stringInArray(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
