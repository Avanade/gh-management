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
	var OutsideCollabsUsers []string
	token := os.Getenv("GH_TOKEN")
	repos, err := ghAPI.GetRepositoriesFromOrganization(org)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Outsidecollabs := ghAPI.ListOutsideCollaborators(token, org)
	for _, list := range Outsidecollabs {
		OutsideCollabsUsers = append(OutsideCollabsUsers, *list.Login)
	}
	var RepoOutsideCollabsList []string
	for _, collab := range repos {
		var RepoCollabsUserNames []string

		var Adminmember []string
		RepoOutsideCollabsList = nil

		RepoCollabs := ghAPI.RepositoriesListCollaborators(token, org, collab.Name, "", "direct")
		for _, list := range RepoCollabs {

			RepoCollabsUserNames = append(RepoCollabsUserNames, *list.Login)
			if *list.RoleName == "admin" {
				Adminmember = append(Adminmember, *list.Login)

			}
		}

		for _, list := range RepoCollabsUserNames {
			for _, Outsidelist := range OutsideCollabsUsers {
				if list == Outsidelist {
					RepoOutsideCollabsList = append(RepoOutsideCollabsList, Outsidelist)
				}
			}
		}
		if len(RepoOutsideCollabsList) > 0 {

			for _, admin := range Adminmember {
				email, err := db.UsersGetEmail(admin)
				if err != nil {
					log.Println(err.Error())
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				EmailAdmin(admin, email, collab.Name, RepoOutsideCollabsList)
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
	users := ghAPI.OrgListMembers(token, organization)
	for _, list := range users {
		email, err := db.UsersGetEmail(*list.Login)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(email) > 0 {
			activeuser, err := msgraph.ActiveUsers(email)
			if err != nil {
				log.Println(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if activeuser == nil {
				ghAPI.RemoveOrganizationsMember(token, organization, *list.Login)

			}
		} else {
			ghAPI.RemoveOrganizationsMember(token, organization, *list.Login)

		}

	}

	// Convert users who are not employees to an outside collaborator
	organizationsOpen := os.Getenv("GH_ORG_OPENSOURCE")

	usersOpenorg := ghAPI.OrgListMembers(token, organizationsOpen)
	for _, list := range usersOpenorg {
		email, err := db.UsersGetEmail(*list.Login)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(email) > 0 {
			activeuser, _ := msgraph.ActiveUsers(email)
			if activeuser == nil {
				ghAPI.ConvertMemberToOutsideCollaborator(token, organizationsOpen, *list.Login)
				ConvertedOutsidecollabsList = append(ConvertedOutsidecollabsList, *list.Login)
			}
		} else {
			ghAPI.ConvertMemberToOutsideCollaborator(token, organizationsOpen, *list.Login)
			ConvertedOutsidecollabsList = append(ConvertedOutsidecollabsList, *list.Login)
		}
	}

	if len(ConvertedOutsidecollabsList) > 0 {
		// Email list of new outside collaborators to ospo
		EmailAdminConvertToColaborator(EmailSupport, ConvertedOutsidecollabsList)

		// Email repo admins with converted users
		repos, _ := ghAPI.GetRepositoriesFromOrganization(organizationsOpen)
		for _, repo := range repos {

			RepoAdmins := GetRepoCollaborators(organizationsOpen, repo.Name, "admin", "direct")
			Repocollabs := GetRepoCollaborators(organizationsOpen, repo.Name, "", "direct")
			var ConvertedInRepo []string
			for _, collab1 := range ConvertedOutsidecollabsList {
				for _, collab2 := range Repocollabs {
					if collab1 == *collab2.Login {
						ConvertedInRepo = append(ConvertedInRepo, collab1)
					}
				}
			}

			for _, collab := range RepoAdmins {
				collabemail, _ := db.UsersGetEmail(*collab.Login)

				if len(ConvertedInRepo) > 0 {
					EmailRepoAdminConvertToColaborator(collabemail, repo.Name, ConvertedInRepo)
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
	EmailSupport := os.Getenv("EMAIL_SUPPORT")
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
			EmailOspoOwnerDeficient(EmailSupport, org, repoOnwerDeficient)
		}
	}
	fmt.Println("REPOOWNERSSCAN SUCCESSFUL")
}

// Repo Collaborators Scan
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

// Deleted repositories from index
func EmailAdminDeletedProjects(to string, repos []string) {
	repoList := "</p> <table  >"
	for _, repo := range repos {
		repoList = repoList + " <tr> <td>" + repo + " </td></tr>"
	}
	repoList = repoList + " </table  > <p>"

	body := fmt.Sprintf("The following repositories were removed from the database as they no longer exist on %s and %s GitHub organizations: %s", os.Getenv("GH_ORG_INNERSOURCE"), os.Getenv("GH_ORG_OPENSOURCE"), repoList)

	m := email.TypEmailMessage{
		Subject: "List of Deleted Repo",
		Body:    body,
		To:      to,
	}

	email.SendEmail(m)
}

// List of users converted into outside collaborators to Repo Owner
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

// List of users converted into outside collaborators to OSPO
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

// List of repos with less than 2 owners to OSPO
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
		Subject: "Repository Owners Scan",
		Body:    body,
		To:      Email,
	}

	email.SendEmail(m)
	fmt.Printf(" less than 2 owner    %s was sent.", e)
}

// List of repos with less than 2 owners to repo owner
func EmailcoownerDeficient(Email string, Org string, reponame string) {
	var body string
	var link string
	link = "https://github.com/" + Org + "/" + reponame + "/settings/access"
	link = "<a href=\"" + link + "\"> here </a>"

	body = fmt.Sprintf("<p>Hello %s ,  </p>  \n<p>This is to inform you that you are the only admin on %s  GitHub repository. We recommend at least 2 admins on each repository. Click %s to add a co-owner.</p> \n <p>OSPO</p>", Email, reponame, link)

	m := email.TypEmailMessage{
		Subject: "Repository Owners Scan",
		Body:    body,
		To:      Email,
	}

	email.SendEmail(m)
}
