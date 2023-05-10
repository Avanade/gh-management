package routes

import (
	"encoding/json"
	"fmt"
	"main/pkg/email"
	db "main/pkg/ghmgmtdb"
	githubAPI "main/pkg/github"
	"main/pkg/msgraph"
	session "main/pkg/session"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/v50/github"
	"github.com/gorilla/mux"
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
	var OutsideCollabsUsers []string
	token := os.Getenv("GH_TOKEN")
	repos, _ := githubAPI.GetRepositoriesFromOrganization(org)
	Outsidecollabs := githubAPI.ListOutsideCollaborators(token, org)
	for _, list := range Outsidecollabs {
		OutsideCollabsUsers = append(OutsideCollabsUsers, *list.Login)
	}
	var RepoOutsideCollabsList []string
	for _, collab := range repos {
		var RepoCollabsUserNames []string

		var Adminmember []string
		RepoOutsideCollabsList = nil

		RepoCollabs := githubAPI.RepositoriesListCollaborators(token, org, collab.Name, "", "direct")
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

	organizationsOpen := [...]string{os.Getenv("GH_ORG_OPENSOURCE"), os.Getenv("GH_ORG_INNERSOURCE")}
	var repoOnwerDeficient []string
	var email string
	EmailSupport := os.Getenv("EMAIL_SUPPORT")
	for _, org := range organizationsOpen {

		repoOnwerDeficient = nil
		repos, _ := githubAPI.GetRepositoriesFromOrganization(org)

		for _, repo := range repos {
			owners := GetRepoCollaborators(org, repo.Name, "", "direct")
			if len(owners) < 2 {
				repoOnwerDeficient = append(repoOnwerDeficient, repo.Name)
				for _, owner := range owners {
					email, _ = db.UsersGetEmail(*owner.Login)

					EmailcoownerDeficient(email, org, repo.Name)
				}
			}

		}

		if len(repoOnwerDeficient) > 0 {
			EmailOspoOwnerDeficient(EmailSupport, org, repoOnwerDeficient)
		}
	}

}

func AddCollaborator(w http.ResponseWriter, r *http.Request) {
	req := mux.Vars(r)
	id, _ := strconv.ParseInt(req["id"], 10, 64)
	ghUser := req["ghUser"]
	permission := req["permission"]

	// Get repository
	data := db.GetProjectById(id)
	s, _ := json.Marshal(data)
	var repoList []Repo
	err := json.Unmarshal(s, &repoList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(repoList) > 0 {
		repo := repoList[0]

		repoUrl := strings.Replace(repo.TFSProjectReference, "https://", "", -1)
		repoUrlSub := strings.Split(repoUrl, "/")

		isInnersource := strings.EqualFold(repoUrlSub[1], os.Getenv("GH_ORG_INNERSOURCE"))
		isMember, _, _ := githubAPI.OrganizationsIsMember(os.Getenv("GH_TOKEN"), ghUser)

		if (isInnersource && isMember) || (!isInnersource) {
			_, err := githubAPI.AddCollaborator(repoUrlSub[1], repo.Name, ghUser, permission)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			users, _ := db.GetUserByGitHubUsername(ghUser)
			if permission == "admin" {

				if len(users) > 0 {
					err = db.RepoOwnersInsert(id, users[0]["UserPrincipalName"].(string))
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
				}
			} else {
				//if not admin, check is the user is currently an admin, remove if he is
				if len(users) > 0 {
					rec, err := db.RepoOwnersByUserAndProjectId(id, users[0]["UserPrincipalName"].(string))
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}

					if len(rec) > 0 {
						err := db.DeleteRepoOwnerRecordByUserAndProjectId(id, users[0]["UserPrincipalName"].(string))
						if err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
							return
						}
					}
				}
			}

			w.WriteHeader(http.StatusOK)
		} else {
			http.Error(w, "Can't invite a user that is not a member of the innersource organization.", http.StatusInternalServerError)
			return
		}
	}

}

func RemoveCollaborator(w http.ResponseWriter, r *http.Request) {
	req := mux.Vars(r)
	id, _ := strconv.ParseInt(req["id"], 10, 64)
	ghUser := req["ghUser"]
	permission := req["permission"]

	// Get repository
	data := db.GetProjectById(id)
	s, _ := json.Marshal(data)
	var repoList []Repo
	err := json.Unmarshal(s, &repoList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(repoList) > 0 {
		repo := repoList[0]

		repoUrl := strings.Replace(repo.TFSProjectReference, "https://", "", -1)
		repoUrlSub := strings.Split(repoUrl, "/")

		_, err := githubAPI.RemoveCollaborator(repoUrlSub[1], repo.Name, ghUser, permission)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if permission == "admin" {
			users, _ := db.GetUserByGitHubUsername(ghUser)

			if len(users) > 0 {
				err = db.DeleteRepoOwnerRecordByUserAndProjectId(id, users[0]["UserPrincipalName"].(string))
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}

		w.WriteHeader(http.StatusOK)
	}

}

func GetTeamsMembers(w http.ResponseWriter, r *http.Request) {

	sessionaz, _ := session.Store.Get(r, "auth-session")
	iaccess_token := sessionaz.Values["access_token"]
	token := fmt.Sprint(iaccess_token)

	token = "eyJ0eXAiOiJKV1QiLCJub25jZSI6IjRsbGxBbkNYa3BOZjBxd195YnQ5S0p3bFBiMzUwVy1TYlBoZTVFVDlPaTgiLCJhbGciOiJSUzI1NiIsIng1dCI6Ii1LSTNROW5OUjdiUm9meG1lWm9YcWJIWkdldyIsImtpZCI6Ii1LSTNROW5OUjdiUm9meG1lWm9YcWJIWkdldyJ9.eyJhdWQiOiIwMDAwMDAwMy0wMDAwLTAwMDAtYzAwMC0wMDAwMDAwMDAwMDAiLCJpc3MiOiJodHRwczovL3N0cy53aW5kb3dzLm5ldC9lMDc5M2QzOS0wOTM5LTQ5NmQtYjEyOS0xOThlZGQ5MTZmZWIvIiwiaWF0IjoxNjgzNzEyMzIwLCJuYmYiOjE2ODM3MTIzMjAsImV4cCI6MTY4Mzc5OTAyMCwiYWNjdCI6MCwiYWNyIjoiMSIsImFjcnMiOlsidXJuOnVzZXI6cmVnaXN0ZXJzZWN1cml0eWluZm8iLCJjMiIsImMxIl0sImFpbyI6IkFkUUFLLzhUQUFBQThXUTQvRDQ4YWpGUUNHRlo1VmdURnlOUTV4Tks5aktkcFNZU2M5ZWpOa0FIbVMwOCtJZ3lvRzhDU09OWVJrRm1qakdFZHFZT2p2M0xIWko1bElXVCtYTDVqS1pIQ2l1ejBGMHJpaGUwUGUxcTJRT0Fnc0pQWXFJOGZVVERkZTM1UWV6NFNvUEE0UndNN3hVSU1ZWTRHbUxwTlNqV2sya2JramE3ZVNNN1VtYkhFRkpJcFM2dDlPT2FUczA2RG5GakRUa2owMElaWEdRMTNjbmV1MFFKTTFIMHMrN2pJUmllYnlrczVaTDBXUnFLT3lsYlJUS09ibllMRTRNbmNadGIxb2lxeXN2clVWUjZtSEZVNWtUejBBPT0iLCJhbXIiOlsicnNhIiwibWZhIl0sImFwcF9kaXNwbGF5bmFtZSI6IkdyYXBoIEV4cGxvcmVyIiwiYXBwaWQiOiJkZThiYzhiNS1kOWY5LTQ4YjEtYThhZC1iNzQ4ZGE3MjUwNjQiLCJhcHBpZGFjciI6IjAiLCJjb250cm9scyI6WyJhcHBfcmVzIl0sImNvbnRyb2xzX2F1ZHMiOlsiMDAwMDAwMDMtMDAwMC0wMDAwLWMwMDAtMDAwMDAwMDAwMDAwIiwiMDAwMDAwMDMtMDAwMC0wZmYxLWNlMDAtMDAwMDAwMDAwMDAwIl0sImRldmljZWlkIjoiMDZlOTQ2NTktYWViMC00YzY1LTlkMTQtNDU3ZGU3NzkyMTY3IiwiZmFtaWx5X25hbWUiOiJEZWxhbWlkYSIsImdpdmVuX25hbWUiOiJEZW5uaXMiLCJpZHR5cCI6InVzZXIiLCJpcGFkZHIiOiIxMzYuMTU4Ljc5LjEiLCJuYW1lIjoiRGVsYW1pZGEsIERlbm5pcyIsIm9pZCI6ImUxYTJlODY2LTRhMDAtNGNlNy1hN2FkLWNlZjcyMWI4ZTg2ZSIsIm9ucHJlbV9zaWQiOiJTLTEtNS0yMS0zMjkwNjgxNTItMTQ1NDQ3MTE2NS0xNDE3MDAxMzMzLTExMTg3MzY4IiwicGxhdGYiOiIzIiwicHVpZCI6IjEwMDMyMDAxNzQzNkU2RjUiLCJyaCI6IjAuQVhzQU9UMTU0RGtKYlVteEtSbU8zWkZ2NndNQUFBQUFBQUFBd0FBQUFBQUFBQUI3QUNzLiIsInNjcCI6IkFwcGxpY2F0aW9uLlJlYWQuQWxsIEF1ZGl0TG9nLlJlYWQuQWxsIENhbGVuZGFycy5SZWFkV3JpdGUgQ29udGFjdHMuUmVhZFdyaXRlIERpcmVjdG9yeS5BY2Nlc3NBc1VzZXIuQWxsIERpcmVjdG9yeS5SZWFkV3JpdGUuQWxsIEZpbGVzLlJlYWRXcml0ZS5BbGwgR3JvdXAuUmVhZFdyaXRlLkFsbCBJZGVudGl0eVJpc2tFdmVudC5SZWFkLkFsbCBJZGVudGl0eVJpc2t5VXNlci5SZWFkLkFsbCBJZGVudGl0eVJpc2t5VXNlci5SZWFkV3JpdGUuQWxsIE1haWwuUmVhZFdyaXRlIE1lbWJlci5SZWFkLkhpZGRlbiBOb3Rlcy5SZWFkV3JpdGUuQWxsIG9wZW5pZCBQZW9wbGUuUmVhZCBQb2xpY3kuUmVhZFdyaXRlLkNvbmRpdGlvbmFsQWNjZXNzIHByb2ZpbGUgUmVwb3J0cy5SZWFkLkFsbCBTZWN1cml0eUV2ZW50cy5SZWFkLkFsbCBTaXRlcy5SZWFkV3JpdGUuQWxsIFRhc2tzLlJlYWRXcml0ZSBVc2VyLlJlYWQgVXNlci5SZWFkLkFsbCBVc2VyLlJlYWRCYXNpYy5BbGwgVXNlci5SZWFkV3JpdGUgVXNlci5SZWFkV3JpdGUuQWxsIGVtYWlsIiwic2lnbmluX3N0YXRlIjpbImR2Y19tbmdkIiwiZHZjX2NtcCIsImR2Y19kbWpkIiwia21zaSJdLCJzdWIiOiJQbjdmV2p3b0tOOTlkSmJPcmtVY1lkRjRlbnFPczVPdkNEQk5RQXpEUWpvIiwidGVuYW50X3JlZ2lvbl9zY29wZSI6IldXIiwidGlkIjoiZTA3OTNkMzktMDkzOS00OTZkLWIxMjktMTk4ZWRkOTE2ZmViIiwidW5pcXVlX25hbWUiOiJkZW5uaXMuZGVsYW1pZGFAYWNjZW50dXJlLmNvbSIsInVwbiI6ImRlbm5pcy5kZWxhbWlkYUBhY2NlbnR1cmUuY29tIiwidXRpIjoiemhtZldaMVEwRU9TRU5yOFRiRktBQSIsInZlciI6IjEuMCIsIndpZHMiOlsiYjc5ZmJmNGQtM2VmOS00Njg5LTgxNDMtNzZiMTk0ZTg1NTA5Il0sInhtc19jYyI6WyJDUDEiXSwieG1zX3NzbSI6IjEiLCJ4bXNfc3QiOnsic3ViIjoiLVZ5TmVuSjh6OGpfbzR4cXZvMnlZdUJXcjZ3TFFRdnFwSlh6LTM5WmhVayJ9LCJ4bXNfdGNkdCI6MTM5NjA0OTM2NH0.N1ksH9ed2UChrPcGVvBdKnkinleEKD_xXAzepaMQhpp_jKstqYA0V5HgbUVveSCgJOZNMTa7W27FtVEgek2UpC-2CUnNXJxZlRCBPmnc4nMxz3RvgswNniShT6TIZzA8YGfmjOnvLbWIX0bQa4O0t5PlWd4Y4_LKcYT1ECwTA4g7MRl2I_WQurmyWVpPQxlD8Dal_FeYmsUbNulFq69nGcbxtuVwaoksM-jCGugX7bbrMZbOhUeQHrd-ehuTFwJ_kgnEsIw7oNcOBwZFRl61LLrkAjyL0gURzurD4wMn_UImEWNcBzpdfBsCLvsuU62xXlVZlToP0E2nNPqWqOJGVg"
	var ChannelId string = "28e03758-7837-493a-97b2-f2e5923b61cc"
	activeuser, _ := msgraph.GetTeamsMembers(ChannelId, token)
	json.NewEncoder(w).Encode(activeuser)
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

func GetRepoCollaborators(org string, repo string, role string, affiliations string) []*github.User {

	token := os.Getenv("GH_TOKEN")

	Repocollabs := githubAPI.RepositoriesListCollaborators(token, org, repo, role, affiliations)

	return Repocollabs
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
		Subject: "Repository Owners Scan",
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
		Subject: "Repository Owners Scan",
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
