package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	db "main/pkg/ghmgmtdb"
	ghAPI "main/pkg/github"
	"main/pkg/notification"
	"main/pkg/session"

	"github.com/google/go-github/v50/github"
	"github.com/gorilla/mux"
)

type RepositoryListDto struct {
	Data  []RepoDto `json:"data"`
	Total int       `json:"total"`
}

type RepoDto struct {
	Id                     int      `json:"Id"`
	Name                   string   `json:"Name"`
	Description            string   `json:"Description"`
	IsArchived             bool     `json:"IsArchived"`
	Created                string   `json:"Created"`
	RepositorySource       string   `json:"RepositorySource"`
	TFSProjectReference    string   `json:"TFSProjectReference"`
	Visibility             string   `json:"Visibility"`
	ApprovalStatus         bool     `json:"ApprovalStatus"`
	ApprovalStatusId       int      `json:"ApprovalStatusId"`
	CoOwner                string   `json:"CoOwner"`
	ConfirmAvaIP           bool     `json:"ConfirmAvaIP"`
	ConfirmEnabledSecurity bool     `json:"ConfirmEnabledSecurity"`
	ECATTID                int      `json:"ECATTID"`
	CreatedBy              string   `json:"CreatedBy"`
	Modified               string   `json:"Modified"`
	ModifiedBy             string   `json:"ModifiedBy"`
	Topics                 []string `json:"RepoTopics"`
}

type CollaboratorDto struct {
	Id                    int64  `json:"id"`
	Name                  string `json:"name"`
	Email                 string `json:"email"`
	Role                  string `json:"role"`
	GitHubUsername        string `json:"ghUsername"`
	IsOutsideCollaborator bool   `json:"isOutsideCollaborator"`
}

type ProjectApprovalSystemPostResponseDto struct {
	ItemId string `json:"itemId"`
}

type ProjectRequest struct {
	Id                         string `json:"id"`
	Newcontribution            string `json:"newcontribution"`
	OSSsponsor                 string `json:"osssponsor"`
	Offeringsassets            string `json:"avanadeofferingsassets"`
	Willbecommercialversion    string `json:"willbecommercialversion"`
	OSSContributionInformation string `json:"osscontributionInformation"`
}

type ProjectApprovalSystemPostDto struct {
	ApplicationId       string
	ApplicationModuleId string
	Email               string
	Subject             string
	Body                string
	RequesterEmail      string
}

type RequestMakePublicDto struct {
	Id                         string `json:"id"`
	Newcontribution            string `json:"newcontribution"`
	OSSsponsor                 string `json:"osssponsor"`
	Offeringsassets            string `json:"avanadeofferingsassets"`
	Willbecommercialversion    string `json:"willbecommercialversion"`
	OSSContributionInformation string `json:"osscontributionInformation"`
}

func RequestRepository(w http.ResponseWriter, r *http.Request) {
	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := profile["preferred_username"]
	r.ParseForm()

	var body db.Project

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !IsRepoNameValid(body.Name) {
		HttpResponseError(w, http.StatusBadRequest, "Invalid repository name.")
		return
	}

	checkDB := make(chan bool)
	checkGH := make(chan bool)

	var existsDb bool
	var existsGH bool
	dashedProjName := strings.ReplaceAll(body.Name, " ", "-")
	go func() { checkDB <- db.ProjectsIsExisting(body.Name) }()
	go func() { b, _ := ghAPI.IsRepoExisting(dashedProjName); checkGH <- b }()

	existsDb = <-checkDB
	existsGH = <-checkGH
	if existsDb || existsGH {
		if existsDb {
			HttpResponseError(w, http.StatusBadRequest, "The project name is existing in the database.")
			return
		} else if existsGH {
			HttpResponseError(w, http.StatusBadRequest, "The project name is existing in Github.")
			return
		}
	} else {
		isEnterpriseOrg, err := ghAPI.IsEnterpriseOrg()
		if err != nil {
			HttpResponseError(w, http.StatusBadRequest, "There is a problem checking if the organization is enterprise or not.")
			return
		}

		repo, errRepo := ghAPI.CreatePrivateGitHubRepository(body.Name, body.Description, username.(string))
		if errRepo != nil {
			log.Println(errRepo.Error())
			HttpResponseError(w, http.StatusInternalServerError, "There is a problem creating the GitHub repository.")
			return
		}
		body.GithubId = repo.GetID()
		body.TFSProjectReference = repo.GetHTMLURL()
		body.Visibility = 1

		innersource := os.Getenv("GH_ORG_INNERSOURCE")
		if isEnterpriseOrg {
			err := ghAPI.SetProjectVisibility(repo.GetName(), "internal", innersource)
			if err != nil {
				return
			}
			body.Visibility = 2
		}

		repoId := db.PRProjectsInsert(body, username.(string))

		// Add  requestor and coowner as repo admins
		err = AddCollaboratorToRequestedRepo(username.(string), body.Name, repoId)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = AddCollaboratorToRequestedRepo(body.Coowner, body.Name, repoId)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		recipients := []string{
			username.(string),
			body.Coowner,
		}

		messageBody := notification.RepositoryHasBeenCreatedMessageBody{
			Recipients:       recipients,
			GitHubAppLink:    os.Getenv("GH_APP_LINK"),
			OrganizationName: innersource,
			RepoLink:         repo.GetName(),
			RepoName:         repo.GetHTMLURL(),
		}
		err = messageBody.Send()
		if err != nil {
			log.Println(err.Error())
		}

		w.WriteHeader(http.StatusOK)
	}
}

func UpdateRepositoryById(w http.ResponseWriter, r *http.Request) {
	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := profile["preferred_username"]
	r.ParseForm()

	var body db.Project

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	db.PRProjectsUpdate(body, username.(string))

	w.WriteHeader(http.StatusOK)
}

func UpdateRepositoryEcattIdById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := profile["preferred_username"]

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var body RepoDto
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	db.UpdateProjectEcattIdById(id, body.ECATTID, username.(string))

	w.WriteHeader(http.StatusOK)
}

func GetUserProjects(w http.ResponseWriter, r *http.Request) {
	// Get email address of the user
	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := profile["preferred_username"]

	// Get project list
	params := make(map[string]interface{})
	params["UserPrincipalName"] = username

	projects, err := db.ProjectsSelectByUserPrincipalName(params)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s, err := json.Marshal(projects)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var list []RepoDto
	err = json.Unmarshal(s, &list)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i := 0; i < len(list); i++ {
		if projects[i]["Topics"] != nil {
			list[i].Topics = strings.Split(projects[i]["Topics"].(string), ",")
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(list)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func GetUsersWithGithub(w http.ResponseWriter, r *http.Request) {

	users := db.GetUsersWithGithub()

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(users)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func GetRequestStatusByProject(w http.ResponseWriter, r *http.Request) {
	req := mux.Vars(r)
	id := req["id"]

	// Get project list
	params := make(map[string]interface{})
	params["Id"] = id

	projects, err := db.ProjectApprovalsSelectById(params)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(projects)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func GetRepoCollaboratorsByRepoId(w http.ResponseWriter, r *http.Request) {
	req := mux.Vars(r)
	id, err := strconv.ParseInt(req["id"], 10, 64)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get repository
	data := db.GetProjectById(id)
	s, err := json.Marshal(data)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var repoList []RepoDto
	err = json.Unmarshal(s, &repoList)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var result []CollaboratorDto

	if len(repoList) > 0 {
		repo := repoList[0]

		if repo.RepositorySource == "GitHub" {
			if repo.TFSProjectReference == "" {
				http.Error(w, "Repository not found.", http.StatusNotFound)
				return
			}
			repoUrl := strings.Replace(repo.TFSProjectReference, "https://", "", -1)
			repoUrlSub := strings.Split(repoUrl, "/")

			token := os.Getenv("GH_TOKEN")

			collaborators := ghAPI.RepositoriesListCollaborators(token, repoUrlSub[1], repo.Name, "", "direct")
			outsideCollaborators := ghAPI.RepositoriesListCollaborators(token, repoUrlSub[1], repo.Name, "", "outside")
			var outsideCollaboratorsUsernames []string
			for _, x := range outsideCollaborators {
				outsideCollaboratorsUsernames = append(outsideCollaboratorsUsernames, *x.Login)
			}

			for _, collaborator := range collaborators {

				// Identify if user is an outside collaborator
				isOutsideCollab := false
				for _, x := range outsideCollaboratorsUsernames {
					if *collaborator.Login == x {
						isOutsideCollab = true
						break
					}
				}

				collabResult := CollaboratorDto{
					Id:                    collaborator.GetID(),
					Role:                  *collaborator.RoleName,
					GitHubUsername:        *collaborator.Login,
					IsOutsideCollaborator: isOutsideCollab,
				}

				//Get user name and email address
				users, err := db.GetUserByGitHubId(strconv.FormatInt(*collaborator.ID, 10))
				if err != nil {
					log.Println(err.Error())
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				if len(users) > 0 {
					collabResult.Name = users[0]["Name"].(string)
					collabResult.Email = users[0]["UserPrincipalName"].(string)
				}

				result = append(result, collabResult)
			}
		} else {
			var collabResult CollaboratorDto

			repoOwner, err := db.GetRepoOwnersRecordByRepoId(int64(repo.Id))
			if err != nil {
				log.Println(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if len(repoOwner) > 0 {
				users, err := db.GetUserByUserPrincipal(repoOwner[0].UserPrincipalName)
				if err != nil {
					log.Println(err.Error())
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				if len(users) > 0 {
					collabResult.Name = users[0]["Name"].(string)
				}

				collabResult.Email = repoOwner[0].UserPrincipalName
				result = append(result, collabResult)
			}

		}

	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(result)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func ArchiveProject(w http.ResponseWriter, r *http.Request) {
	req := mux.Vars(r)
	project := req["project"]
	projectId := req["projectId"]
	state := req["state"]
	archive := req["archive"]

	organization := os.Getenv("GH_ORG_INNERSOURCE")
	if state == "Public" {
		organization = os.Getenv("GH_ORG_OPENSOURCE")
	}

	if archive == "1" {
		err := ghAPI.ArchiveProject(project, archive == "1", organization)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		isArchived, err := ghAPI.IsArchived(project, organization)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if isArchived {
			http.Error(w, "Repository is still archived. Unarchive the repo on GitHub and try again.", http.StatusBadRequest)
			return
		}
	}
	id, err := strconv.ParseInt(projectId, 10, 64)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	db.UpdateProjectIsArchived(id, archive == "1")
	w.WriteHeader(http.StatusOK)
}

func GetAllRepositories(w http.ResponseWriter, r *http.Request) {

	params := r.URL.Query()
	search := params["search"][0]
	offset, err := strconv.Atoi(params["offset"][0])
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get repository list
	data := db.ReposSelectByOffsetAndFilter(offset, search)
	s, _ := json.Marshal(data)
	var list []RepoDto
	err = json.Unmarshal(s, &list)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i := 0; i < len(list); i++ {
		if data[i]["Topics"] != nil {
			list[i].Topics = strings.Split(data[i]["Topics"].(string), ",")
		}
	}

	result := RepositoryListDto{
		Data:  list,
		Total: db.ReposTotalCountBySearchTerm(search),
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(result)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func SetVisibility(w http.ResponseWriter, r *http.Request) {

	req := mux.Vars(r)
	project := req["project"]
	projectId := req["projectId"]
	currentState := req["currentState"]
	desiredState := req["desiredState"]
	visibilityId := 1 //public
	if desiredState == "internal" {
		visibilityId = 2 //internal
	}

	innersource := os.Getenv("GH_ORG_INNERSOURCE")
	opensource := os.Getenv("GH_ORG_OPENSOURCE")

	id, err := strconv.ParseInt(projectId, 10, 64)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if currentState == "Public" {
		// Set repo to desired visibility then move to innersource
		err := ghAPI.SetProjectVisibility(project, desiredState, opensource)
		if err != nil {
			http.Error(w, "Failed to make the repository "+desiredState, http.StatusInternalServerError)
			return
		}

		ghAPI.TransferRepository(project, opensource, innersource)

		time.Sleep(3 * time.Second)
		repoResp, err := ghAPI.GetRepository(project, innersource)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = db.UpdateTFSProjectReferenceById(id, repoResp.GetHTMLURL())
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// Set repo to desired visibility
		err := ghAPI.SetProjectVisibility(project, desiredState, innersource)
		if err != nil {
			http.Error(w, "Failed to make the repository "+desiredState, http.StatusInternalServerError)
			return
		}
	}

	db.UpdateProjectVisibilityId(id, int64(visibilityId))

	w.WriteHeader(http.StatusOK)
}

func RequestMakePublic(w http.ResponseWriter, r *http.Request) {
	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := profile["preferred_username"]
	r.ParseForm()

	var body RequestMakePublicDto

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	projectRequest := db.ProjectRequest{
		Id:                         body.Id,
		Newcontribution:            body.Newcontribution,
		OSSsponsor:                 body.OSSsponsor,
		Offeringsassets:            body.Offeringsassets,
		Willbecommercialversion:    body.Willbecommercialversion,
		OSSContributionInformation: body.OSSContributionInformation,
	}

	db.PRProjectsUpdateLegalQuestions(projectRequest, username.(string))

	id, err := strconv.ParseInt(projectRequest.Id, 10, 64)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	go RequestApproval(id, username.(string))
}

func IndexOrgRepos(w http.ResponseWriter, r *http.Request) {
	fmt.Println("INDEX ORGANIZATION REPOSITORIES TRIGGERED...")
	var repos []ghAPI.Repo

	orgs := []string{os.Getenv("GH_ORG_INNERSOURCE"), os.Getenv("GH_ORG_OPENSOURCE")}

	for _, org := range orgs {
		reposByOrg, err := ghAPI.GetRepositoriesFromOrganization(org)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if reposByOrg != nil {
			repos = append(repos, reposByOrg...)
		}
	}

	var wg sync.WaitGroup
	maxGoroutines := 50
	guard := make(chan struct{}, maxGoroutines)

	for _, repo := range repos {
		guard <- struct{}{}
		wg.Add(1)
		go func(r ghAPI.Repo) {
			IndexRepo(r)
			<-guard
			wg.Done()
		}(repo)
	}
	wg.Wait()

	w.WriteHeader(http.StatusOK)
	fmt.Println("INDEX ORGANIZATION REPOSITORIES SUCCESSFUL")
}

func ClearOrgRepos(w http.ResponseWriter, r *http.Request) {
	projects, err := db.ProjectsByRepositorySource("GitHub")
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	removedProjects := []string{}

	var wg sync.WaitGroup
	maxGoroutines := 50
	guard := make(chan struct{}, maxGoroutines)

	for _, project := range projects {
		guard <- struct{}{}
		wg.Add(1)
		go func(p map[string]interface{}) {
			projectId := p["Id"].(int64)
			repoName := p["Name"].(string)
			var isGithubIdNil bool
			if p["GithubId"] == nil {
				isGithubIdNil = true
			} else {
				isGithubIdNil = false
			}
			isRemoved := RemoveRepoIfNotExist(int(projectId), repoName, isGithubIdNil)
			if isRemoved {
				removedProjects = append(removedProjects, repoName)
			}
			<-guard
			wg.Done()
		}(project)
	}
	wg.Wait()

	if len(removedProjects) > 0 {
		emailSupport := os.Getenv("EMAIL_SUPPORT")
		EmailAdminDeletedProjects(emailSupport, removedProjects)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Println(" SUCCESSFULLY INDEXED ORGANIZATION REPOSITORIES")
}

func AddCollaborator(w http.ResponseWriter, r *http.Request) {
	req := mux.Vars(r)
	id, _ := strconv.ParseInt(req["id"], 10, 64)
	ghUser := req["ghUser"]
	permission := req["permission"]

	// Get repository
	data := db.GetProjectById(id)
	s, err := json.Marshal(data)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var repoList []RepoDto
	err = json.Unmarshal(s, &repoList)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(repoList) > 0 {
		repo := repoList[0]

		if repo.TFSProjectReference == "" {
			http.Error(w, "Repository doesn't exists on GitHub organization.", http.StatusNotFound)
			return
		}

		repoUrl := strings.Replace(repo.TFSProjectReference, "https://", "", -1)
		repoUrlSub := strings.Split(repoUrl, "/")

		isInnersource := strings.EqualFold(repoUrlSub[1], os.Getenv("GH_ORG_INNERSOURCE"))
		isMember, _, _ := ghAPI.OrganizationsIsMember(os.Getenv("GH_TOKEN"), ghUser)

		if (isInnersource && isMember) || (!isInnersource) {
			_, err := ghAPI.AddCollaborator(repoUrlSub[1], repo.Name, ghUser, permission)
			if err != nil {
				log.Println(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			users, _ := db.GetUserByGitHubUsername(ghUser)
			if permission == "admin" {

				if len(users) > 0 {
					db.RepoOwnersInsert(id, users[0]["UserPrincipalName"].(string))
				}
			} else {
				//if not admin, check is the user is currently an admin, remove if he is
				if len(users) > 0 {
					rec, err := db.RepoOwnersByUserAndProjectId(id, users[0]["UserPrincipalName"].(string))
					if err != nil {
						log.Println(err.Error())
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}

					if len(rec) > 0 {
						err := db.DeleteRepoOwnerRecordByUserAndProjectId(id, users[0]["UserPrincipalName"].(string))
						if err != nil {
							log.Println(err.Error())
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
	s, err := json.Marshal(data)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var repoList []RepoDto
	err = json.Unmarshal(s, &repoList)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(repoList) > 0 {
		repo := repoList[0]

		if repo.TFSProjectReference == "" {
			http.Error(w, "Repository doesn't exists on GitHub organization.", http.StatusNotFound)
			return
		}

		repoUrl := strings.Replace(repo.TFSProjectReference, "https://", "", -1)
		repoUrlSub := strings.Split(repoUrl, "/")

		_, err := ghAPI.RemoveCollaborator(repoUrlSub[1], repo.Name, ghUser, permission)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if permission == "admin" {
			users, _ := db.GetUserByGitHubUsername(ghUser)

			if len(users) > 0 {
				err = db.DeleteRepoOwnerRecordByUserAndProjectId(id, users[0]["UserPrincipalName"].(string))
				if err != nil {
					log.Println(err.Error())
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}

		w.WriteHeader(http.StatusOK)
	}
}

func RepoOwnersCleanup(w http.ResponseWriter, r *http.Request) {
	fmt.Println("REPO OWNERS CLEANUP TRIGGERED...")
	token := os.Getenv("GH_TOKEN")

	// Get all repos from database
	repos, err := db.GetGitHubRepositories()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var wg sync.WaitGroup
	maxGoroutines := 50
	guard := make(chan struct{}, maxGoroutines)

	for _, repo := range repos {
		guard <- struct{}{}
		wg.Add(1)
		go func(r map[string]interface{}) {
			CleanupRepoOwners(r, token)
			<-guard
			wg.Done()
		}(repo)
	}
	wg.Wait()
	w.WriteHeader(http.StatusOK)
	fmt.Println("REPO OWNERS CLEANUP SUCCESSFUL")
}

func GetRepoCollaborators(org string, repo string, role string, affiliations string) []*github.User {

	token := os.Getenv("GH_TOKEN")

	repoCollabs := ghAPI.RepositoriesListCollaborators(token, org, repo, role, affiliations)

	return repoCollabs
}

func CleanupRepoOwners(repo map[string]interface{}, token string) {
	fmt.Println("Checking owners of : " + repo["Name"].(string))

	if repo["TFSProjectReference"] == nil {
		// Get owners of the repo on the database
		owners, err := db.GetRepoOwnersByProjectIdWithGHUsername(repo["Id"].(int64))
		if err != nil {
			log.Println(err.Error())
			return
		}

		// Check if owner is on the admin list
		for _, owner := range owners {
			db.DeleteRepoOwnerRecordByUserAndProjectId(repo["Id"].(int64), owner["UserPrincipalName"].(string))
		}

		return
	}

	// Get admins of the repository
	repoUrl := strings.Replace(repo["TFSProjectReference"].(string), "https://", "", -1)
	repoUrlSub := strings.Split(repoUrl, "/")

	admins := ghAPI.RepositoriesListCollaborators(token, repoUrlSub[1], repo["Name"].(string), "admin", "direct")

	// Get owners of the repo on the database
	owners, err := db.GetRepoOwnersByProjectIdWithGHUsername(repo["Id"].(int64))
	if err != nil {
		log.Println(err.Error())
		return
	}

	// Check if owner is on the admin list
	for _, owner := range owners {
		isAdmin := false
		for _, admin := range admins {
			if owner["GitHubUser"].(string) == *admin.Login {
				isAdmin = true
				break
			}
		}

		// if owner is not on the list of admins
		if !isAdmin {
			db.DeleteRepoOwnerRecordByUserAndProjectId(repo["Id"].(int64), owner["UserPrincipalName"].(string))
		}
	}
}

func RemoveRepoIfNotExist(projectId int, repoName string, isGithubIdNil bool) bool {
	isExist, err := ghAPI.IsRepoExisting(repoName)
	if err != nil {
		fmt.Println(err.Error(), " | REPO NAME : ", repoName)
		return false
	}

	if !isExist || isGithubIdNil {
		err := db.DeleteProjectById(projectId)
		if err != nil {
			fmt.Println(err.Error())
			return false
		}
		return true
	}
	return false
}

func IndexRepo(repo ghAPI.Repo) {
	fmt.Println("Indexing " + repo.Name + "...")

	visibilityId := 3
	if repo.Visibility == "private" {
		visibilityId = 1
	} else if repo.Visibility == "internal" {
		visibilityId = 2
	}

	param := map[string]interface{}{
		"GithubId":            repo.GithubId,
		"Name":                repo.Name,
		"Description":         repo.Description,
		"IsArchived":          repo.IsArchived,
		"VisibilityId":        visibilityId,
		"TFSProjectReference": repo.TFSProjectReference,
		"Created":             repo.Created.Format("2006-01-02 15:04:05"),
	}

	isExisting := db.ProjectsIsExistingByGithubId(repo.GithubId)
	var projectId int

	if isExisting {
		project := db.GetProjectByGithubId(repo.GithubId)
		param["Id"] = project[0]["Id"]
		projectId = int(project[0]["Id"].(int64))

		err := db.ProjectUpdateByImport(param)
		if err != nil {
			log.Println(err.Error())
			return
		}

		// Get direct admin collaborators
		repoUrl := strings.Replace(repo.TFSProjectReference, "https://", "", -1)
		repoUrlSub := strings.Split(repoUrl, "/")

		token := os.Getenv("GH_TOKEN")

		collaborators := ghAPI.RepositoriesListCollaborators(token, repoUrlSub[1], repo.Name, "admin", "direct")
		// Get userprincipal from database
		for _, admin := range collaborators {
			users, err := db.GetUserByGitHubId(strconv.FormatInt(*admin.ID, 10))
			if err != nil {
				log.Println(err.Error())
				return
			}

			//Insert to repoowners table
			if len(users) > 0 {
				if len(users) > 0 {
					db.RepoOwnersInsert(project[0]["Id"].(int64), users[0]["UserPrincipalName"].(string))
				}
			}
		}
	} else {
		err := db.ProjectInsertByImport(param)
		if err != nil {
			log.Println(err.Error())
			return
		}

		project := db.GetProjectByGithubId(repo.GithubId)
		param["Id"] = project[0]["Id"]
		projectId = int(project[0]["Id"].(int64))

		// Get direct admin collaborators
		repoUrl := strings.Replace(repo.TFSProjectReference, "https://", "", -1)
		repoUrlSub := strings.Split(repoUrl, "/")

		token := os.Getenv("GH_TOKEN")

		collaborators := ghAPI.RepositoriesListCollaborators(token, repoUrlSub[1], repo.Name, "admin", "direct")
		// Get userprincipal from database
		for _, admin := range collaborators {
			users, err := db.GetUserByGitHubId(strconv.FormatInt(*admin.ID, 10))
			if err != nil {
				log.Println(err.Error())
				return
			}

			//Insert to repoowners table
			if len(users) > 0 {
				if len(users) > 0 {
					db.RepoOwnersInsert(project[0]["Id"].(int64), users[0]["UserPrincipalName"].(string))
				}
			}
		}
	}
	if len(repo.Topics) > 0 {
		err := db.DeleteProjectTopics(projectId)
		if err != nil {
			log.Println(err.Error())
			return
		}
		for i := 0; i < len(repo.Topics); i++ {
			err := db.InsertProjectTopics(projectId, repo.Topics[i])
			if err != nil {
				log.Println(err.Error())
				return
			}
		}
	}
}

func RequestApproval(id int64, email string) {

	projectApprovals := db.PopulateProjectsApproval(id, email)

	for _, v := range projectApprovals {
		if v.RequestStatus == "New" {
			err := ApprovalSystemRequest(v)
			if err != nil {
				log.Println("ID:" + strconv.FormatInt(v.Id, 10) + " " + err.Error())
				return
			}
		}
	}
}

func ApprovalSystemRequest(data db.ProjectApproval) error {

	url := os.Getenv("APPROVAL_SYSTEM_APP_URL")
	if url != "" {
		url = url + "/request"
		ch := make(chan *http.Response)
		// var res *http.Response

		bodyTemplate := `<p>Hi |ApproverUserPrincipalName|!</p>
		<p>|RequesterName| is requesting for a new project and is now pending for |ApprovalType| review.</p>
		<p>Below are the details:</p>
		<table>
			<tr>
				<td style="font-weight: bold;">Project Name<td>
				<td style="font-size:larger">|ProjectName|<td>
			</tr>
			<tr>
				<td style="font-weight: bold;">Requested by<td>
				<td style="font-size:larger">|Requester|<td>
			</tr>
			<tr>
				<td style="font-weight: bold;">Description<td>
				<td style="font-size:larger">|ProjectDescription|<td>
			</tr>
		</table>
		<table>
			<tr>
				<td style="font-weight: bold;">Is this a new contribution with no prior code development? (i.e., no existing Avanade IP, no third-party/OSS code, etc.)<td>
				<td style="font-size:larger">|Newcontribution|<td>
			</tr>
			<tr>
				<td style="font-weight: bold;">Who is sponsoring thapprovalsyscois OSS contribution?<td>
				<td style="font-size:larger">|OSSsponsor|<td>
			</tr>
			<tr>
				<td style="font-weight: bold;">Will Avanade use this contribution in client accounts and/or as part of an Avanade offerings/assets?<td>
				<td style="font-size:larger">|Avanadeofferingsassets|<td>
			</tr>
			<tr>
				<td style="font-weight: bold;">Will there be a commercial version of this contribution<td>
				<td style="font-size:larger">|Willbecommercialversion|<td>
			</tr>
				<tr>
				<td style="font-weight: bold;">Additional OSS Contribution Information (e.g. planned maintenance/support, etc.)?<td>
				<td style="font-size:larger">|OSSContributionInformation|<td>
			</tr>
		</table>
		<p>For more information, send an email to <a href="mailto:|RequesterUserPrincipalName|">|RequesterUserPrincipalName|</a></p>
		`
		replacer := strings.NewReplacer("|ApproverUserPrincipalName|", data.ApproverUserPrincipalName,
			"|RequesterName|", data.RequesterName,
			"|ApprovalType|", data.ApprovalType,
			"|ProjectName|", data.ProjectName,
			"|Requester|", data.RequesterName,
			"|ProjectDescription|", data.ProjectDescription,
			"|RequesterUserPrincipalName|", data.RequesterUserPrincipalName,

			"|Newcontribution|", data.Newcontribution,
			"|OSSsponsor|", data.OSSsponsor,
			"|Avanadeofferingsassets|", data.Offeringsassets,
			"|Willbecommercialversion|", data.Willbecommercialversion,
			"|OSSContributionInformation|", data.OSSContributionInformation,
		)
		body := replacer.Replace(bodyTemplate)
		postParams := ProjectApprovalSystemPostDto{
			ApplicationId:       os.Getenv("APPROVAL_SYSTEM_APP_ID"),
			ApplicationModuleId: os.Getenv("APPROVAL_SYSTEM_APP_MODULE_PROJECTS"),
			Email:               data.ApproverUserPrincipalName,
			Subject:             fmt.Sprintf("[GH-Management] New Project For Review - %v", data.ProjectName),
			Body:                body,
			RequesterEmail:      data.RequesterUserPrincipalName,
		}

		go getHttpPostResponseStatus(url, postParams, ch)
		r := <-ch
		if r != nil {
			var res ProjectApprovalSystemPostResponseDto
			err := json.NewDecoder(r.Body).Decode(&res)
			if err != nil {
				return err
			}

			db.ProjectsApprovalUpdateGUID(data.Id, res.ItemId)
		}
	}
	return nil
}

func getHttpPostResponseStatus(url string, data interface{}, ch chan *http.Response) {
	jsonReq, err := json.Marshal(data)
	if err != nil {
		log.Println(err.Error())
		ch <- nil
	}
	res, err := http.Post(url, "application/json; charset=utf-8", bytes.NewBuffer(jsonReq))
	if err != nil {
		ch <- nil
	}
	ch <- res
}

func ReprocessRequestApproval() {
	projectApprovals := db.GetFailedProjectApprovalRequests()

	for _, v := range projectApprovals {
		go ApprovalSystemRequest(v)
	}
}

func IsRepoNameValid(value string) bool {
	regex := regexp.MustCompile(`^([a-zA-Z0-9\-\_]|\.{3,}|\.{1,}[a-zA-Z0-9\-\_])([a-zA-Z0-9\-\_\.]+)?`)
	matched := regex.FindAllString(value, 1)

	if matched == nil {
		return false
	}

	return matched[0] == value
}

func AddCollaboratorToRequestedRepo(user string, repo string, repoId int64) error {
	innersource := os.Getenv("GH_ORG_INNERSOURCE")
	gHUser := db.Users_Get_GHUser(user)
	isInnersourceMember, _, _ := ghAPI.OrganizationsIsMember(os.Getenv("GH_TOKEN"), gHUser)
	if isInnersourceMember {
		_, err := ghAPI.AddCollaborator(innersource, repo, gHUser, "admin")
		if err != nil {
			return err
		}
		err = db.RepoOwnersInsert(repoId, user)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetPopularTopics(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	offset, err := strconv.Atoi(params["offset"][0])
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rowCount, err := strconv.Atoi(params["rowCount"][0])
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	result, err := db.GetPopularTopics(offset, rowCount)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(result)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func HttpResponseError(w http.ResponseWriter, code int, errorMessage string) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(code)
	response, err := json.Marshal(errorMessage)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}
