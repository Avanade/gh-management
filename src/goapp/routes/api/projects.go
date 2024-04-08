package routes

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"main/pkg/appinsights_wrapper"
	"main/pkg/email"
	db "main/pkg/ghmgmtdb"
	ghAPI "main/pkg/github"
	"main/pkg/notification"
	"main/pkg/session"

	"github.com/google/go-github/v50/github"
	"github.com/gorilla/mux"
	"github.com/microsoft/ApplicationInsights-Go/appinsights/contracts"
)

type RepositoryListDto struct {
	Data  []RepoDto `json:"data"`
	Total int       `json:"total"`
}

type RepoDto struct {
	Id                     int      `json:"Id"`
	Name                   string   `json:"Name"`
	AssetCode              string   `json:"AssetCode"`
	Organization           string   `json:"Organization"`
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
	ApplicationId       string   `json:"applicationId"`
	ApplicationModuleId string   `json:"applicationModuleId"`
	Emails              []string `json:"emails"`
	Subject             string   `json:"subject"`
	Body                string   `json:"body"`
	RequesterEmail      string   `json:"requesterEmail"`
}

type RequestMakePublicDto struct {
	Id                         string `json:"id"`
	Newcontribution            string `json:"newcontribution"`
	OSSsponsor                 int    `json:"osssponsor"`
	Offeringsassets            string `json:"avanadeofferingsassets"`
	Willbecommercialversion    string `json:"willbecommercialversion"`
	OSSContributionInformation string `json:"osscontributionInformation"`
}

type TransferRepoDto struct {
	NewOrg string `json:"newOrg"`
}

func CreateRepository(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := profile["preferred_username"]
	r.ParseForm()

	var body db.Project

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !IsRepoNameValid(body.Name) {
		logger.TrackTrace("Invalid repository name.", contracts.Error)
		HttpResponseError(w, http.StatusBadRequest, "Invalid repository name.", logger)
		return
	}

	checkDB := make(chan bool)
	checkGH := make(chan bool)

	var existsDb bool
	var existsGH bool
	dashedProjName := strings.ReplaceAll(body.Name, " ", "-")
	description := strings.ReplaceAll(body.Description, "\n", " ")

	go func() { checkDB <- db.ProjectsIsExisting(body.Name) }()
	go func() { b, _ := ghAPI.IsRepoExisting(dashedProjName); checkGH <- b }()

	existsDb = <-checkDB
	existsGH = <-checkGH
	if existsDb || existsGH {
		if existsDb {
			logger.TrackTrace("Project name exists on the database.", contracts.Error)
			HttpResponseError(w, http.StatusBadRequest, "The project name is existing in the database.", logger)
			return
		} else if existsGH {
			logger.TrackTrace("Project name exists on GitHub org.", contracts.Error)
			HttpResponseError(w, http.StatusBadRequest, "The project name is existing in Github.", logger)
			return
		}
	} else {
		isEnterpriseOrg, err := ghAPI.IsEnterpriseOrg()
		if err != nil {
			logger.LogException(err)
			HttpResponseError(w, http.StatusBadRequest, "There is a problem checking if the organization is enterprise or not.", logger)
			return
		}

		logger.LogTrace("Creating repository...", contracts.Information)
		repo, err := ghAPI.CreatePrivateGitHubRepository(body.Name, description, username.(string))

		if err != nil {
			logger.LogException(err)
			HttpResponseError(w, http.StatusInternalServerError, "There is a problem creating the GitHub repository.", logger)
			return
		}

		logger.LogTrace(repo.GetName(), contracts.Information) // TEMP LOG - END TEMP LOG

		body.AssetCode = body.Name
		body.GithubId = repo.GetID()
		body.TFSProjectReference = repo.GetHTMLURL()

		innersource := os.Getenv("GH_ORG_INNERSOURCE")
		body.Organization = innersource
		if isEnterpriseOrg && body.Visibility == 2 {
			logger.LogTrace("Making the repository as internal...", contracts.Information)
			_, err := ghAPI.SetProjectVisibility(repo.GetName(), "internal", innersource)
			if err != nil {
				logger.LogException(err)
				HttpResponseError(w, http.StatusInternalServerError, err.Error(), logger)
				return
			}
		}

		logger.LogTrace("Adding repository to database...", contracts.Information)
		repoId := db.PRProjectsInsert(body, username.(string))

		// Add  requestor and coowner as repo admins
		for x := 1; x <= 3; x++ {
			time.Sleep(2 * time.Second)
			logger.LogTrace(fmt.Sprintf("Attempt %d: Adding requestor as a collaborator...", x), contracts.Information)
			resp, err := AddCollaboratorToRequestedRepo(username.(string), body.Name, repoId, logger)
			if err != nil {
				logger.LogException(err)
				continue
			}
			if resp.StatusCode != 403 {
				break
			}
		}

		for x := 1; x <= 3; x++ {
			time.Sleep(2 * time.Second)
			logger.LogTrace(fmt.Sprintf("Attempt %d: Adding coowner as a collaborator...", x), contracts.Information)
			resp, err := AddCollaboratorToRequestedRepo(body.Coowner, body.Name, repoId, logger)
			if err != nil {
				logger.LogException(err)
				continue
			}
			if resp.StatusCode != 403 {
				break
			}
		}

		recipients := []string{
			username.(string),
			body.Coowner,
		}

		messageBody := notification.RepositoryHasBeenCreatedMessageBody{
			Recipients:       recipients,
			GitHubAppLink:    os.Getenv("GH_APP_LINK"),
			OrganizationName: innersource,
			RepoLink:         repo.GetHTMLURL(),
			RepoName:         repo.GetName(),
		}
		err = messageBody.Send()
		if err != nil {
			logger.LogException(err)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	db.PRProjectsUpdate(body, username.(string))

	w.WriteHeader(http.StatusOK)
}

func UpdateRepositoryEcattIdById(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	vars := mux.Vars(r)

	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := profile["preferred_username"]

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var body RepoDto
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = db.UpdateProjectEcattIdById(id, body.ECATTID, username.(string))
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetMyRepositories(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

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
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s, err := json.Marshal(projects)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var list []RepoDto
	err = json.Unmarshal(s, &list)
	if err != nil {
		logger.LogException(err)
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
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func GetUsersWithGithub(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	users := db.GetUsersWithGithub()

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(users)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func GetRequestStatusByRepoId(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	req := mux.Vars(r)
	id := req["id"]

	// Get project list
	params := make(map[string]interface{})
	params["Id"] = id

	projects, err := db.ProjectApprovalsSelectById(params)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(projects)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func GetRepositoryReadmeById(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	req := mux.Vars(r)
	repoName := req["repoName"]
	visibility := req["visibility"]

	readme, _ := ghAPI.GetRepositoryReadmeById(repoName, visibility)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(readme)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func GetRepoCollaboratorsByRepoId(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	req := mux.Vars(r)
	id, err := strconv.ParseInt(req["id"], 10, 64)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get repository
	data := db.GetProjectById(id)
	s, err := json.Marshal(data)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var repoList []RepoDto
	err = json.Unmarshal(s, &repoList)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var result []CollaboratorDto

	if len(repoList) > 0 {
		repo := repoList[0]

		if repo.RepositorySource == "GitHub" {
			if repo.TFSProjectReference == "" {
				logger.LogTrace("Repository not found", contracts.Error)
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
					logger.LogException(err)
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
				logger.LogException(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if len(repoOwner) > 0 {
				users, err := db.GetUserByUserPrincipal(repoOwner[0].UserPrincipalName)
				if err != nil {
					logger.LogException(err)
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
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func ArchiveProject(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

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
			logger.LogException(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		isArchived, err := ghAPI.IsArchived(project, organization)
		if err != nil {
			logger.LogException(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if isArchived {
			logger.LogTrace("Repository is still archived. Unarchive the repo on GitHub and try again.", contracts.Error)
			http.Error(w, "Repository is still archived. Unarchive the repo on GitHub and try again.", http.StatusBadRequest)
			return
		}
	}
	id, err := strconv.ParseInt(projectId, 10, 64)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = db.UpdateProjectIsArchived(id, archive == "1")
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetRepositories(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	params := r.URL.Query()
	search := params["search"][0]
	offset, err := strconv.Atoi(params["offset"][0])
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get repository list
	data := db.ReposSelectByOffsetAndFilter(offset, search)
	s, err := json.Marshal(data)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var list []RepoDto
	err = json.Unmarshal(s, &list)
	if err != nil {
		logger.LogException(err)
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
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func GetTotalRepositoriesOwnedByUsers(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	vars := mux.Vars(r)
	user := vars["user"]

	visibility := 0

	params := r.URL.Query()

	// 0/NONE - ALL | 1 - PRIVATE | 2 - INTERNAL | 3 - PUBLIC
	if params.Has("visibility") {
		visibility, _ = strconv.Atoi(params["visibility"][0])
	}

	organization := os.Getenv("GH_ORG_INNERSOURCE")
	// public - GH_ORG_OPENSOURCE | private/none - GH_ORG_INNERSOURCE
	if params.Has("orgtype") {
		if params["orgtype"][0] == "public" {
			organization = os.Getenv("GH_ORG_OPENSOURCE")
		}
	}

	if user == "me" {
		sessionaz, _ := session.Store.Get(r, "auth-session")
		iprofile := sessionaz.Values["profile"]
		profile := iprofile.(map[string]interface{})
		user = fmt.Sprint(profile["preferred_username"])
	}

	total, err := db.CountOwnedRepoByVisibility(user, organization, visibility)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(struct {
		Total int `json:"total"`
	}{
		Total: total,
	})
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func GetRepositoriesById(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	req := mux.Vars(r)

	id, err := strconv.ParseInt(req["id"], 10, 64)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := db.GetProjectById(id)
	s, err := json.Marshal(data)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var repo []RepoDto
	err = json.Unmarshal(s, &repo)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if data[0]["Topics"] != nil {
		repo[0].Topics = strings.Split(data[0]["Topics"].(string), ",")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(repo[0])
}

func SetVisibility(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

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
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if currentState == "Public" {
		// Set repo to desired visibility then move to innersource
		_, err := ghAPI.SetProjectVisibility(project, desiredState, opensource)
		if err != nil {
			logger.LogException(err)
			http.Error(w, "Failed to make the repository "+desiredState, http.StatusInternalServerError)
			return
		}

		_, err = ghAPI.TransferRepository(project, opensource, innersource)
		if err != nil {
			logger.LogException(err)
			http.Error(w, "Failed to make the repository "+desiredState, http.StatusInternalServerError)
			return
		}

		time.Sleep(3 * time.Second)
		repoResp, err := ghAPI.GetRepository(project, innersource)
		if err != nil {
			logger.LogException(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = db.UpdateTFSProjectReferenceById(id, repoResp.GetHTMLURL(), *repoResp.GetOwner().Login)
		if err != nil {
			logger.LogException(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// Set repo to desired visibility
		_, err := ghAPI.SetProjectVisibility(project, desiredState, innersource)
		if err != nil {
			logger.LogException(err)
			http.Error(w, "Failed to make the repository "+desiredState, http.StatusInternalServerError)
			return
		}
	}

	err = db.UpdateProjectVisibilityId(id, int64(visibilityId))
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func TransferRepository(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	vars := mux.Vars(r)

	projectId, err := strconv.Atoi(vars["projectId"])
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var body TransferRepoDto
	r.ParseForm()
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var newOwner string

	if body.NewOrg == "innersource" {
		newOwner = os.Getenv("GH_ORG_INNERSOURCE")
	} else if body.NewOrg == "opensource" {
		newOwner = os.Getenv("GH_ORG_OPENSOURCE")
	} else {
		newOwner = body.NewOrg
	}

	project := db.GetProjectById(int64(projectId))
	name := project[0]["Name"].(string)
	owner := project[0]["Organization"].(string)

	ValidateOrgMembers(owner, name, newOwner, logger)

	_, err = ghAPI.TransferRepository(name, owner, newOwner)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	time.Sleep(3 * time.Second)
	repository, err := ghAPI.GetRepository(name, newOwner)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = db.UpdateTFSProjectReferenceById(int64(projectId), repository.GetHTMLURL(), *repository.GetOwner().Login)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func RequestMakePublic(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := profile["preferred_username"]
	r.ParseForm()

	var body RequestMakePublicDto

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		logger.LogException(err)
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
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	RequestApproval(id, username.(string), logger)
}

func IndexOrgRepos(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	var repos []ghAPI.Repo

	orgs := []string{os.Getenv("GH_ORG_INNERSOURCE"), os.Getenv("GH_ORG_OPENSOURCE")}

	regOrgs, err := db.GetAllRegionalOrganizations()
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, regOrg := range regOrgs {
		orgs = append(orgs, regOrg["Name"].(string))
	}

	for _, org := range orgs {
		reposByOrg, err := ghAPI.GetRepositoriesFromOrganization(org)
		if err != nil {
			logger.LogException(err)
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
			indexRepo(r, logger)
			<-guard
			wg.Done()
		}(repo)
	}
	wg.Wait()

	w.WriteHeader(http.StatusOK)
	logger.TrackTrace("Index organization repositories successful", contracts.Information)
}

func ClearOrgRepos(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	projects, err := db.ProjectsByRepositorySource("GitHub")
	if err != nil {
		logger.LogException(err)
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
			isRemoved := RemoveRepoIfNotExist(int(projectId), repoName, isGithubIdNil, logger)
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
		emailAdminDeletedProjects(emailSupport, removedProjects, logger)
	}

	w.WriteHeader(http.StatusOK)
	logger.TrackTrace("Clear org repos successful", contracts.Information)
}

func AddCollaborator(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	req := mux.Vars(r)
	id, _ := strconv.ParseInt(req["id"], 10, 64)
	ghUser := req["ghUser"]
	permission := req["permission"]

	// Get repository
	data := db.GetProjectById(id)
	s, err := json.Marshal(data)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var repoList []RepoDto
	err = json.Unmarshal(s, &repoList)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(repoList) > 0 {
		repo := repoList[0]

		if repo.TFSProjectReference == "" {
			logger.LogTrace("Repository doesn't exists on GitHub organization.", contracts.Error)
			http.Error(w, "Repository doesn't exists on GitHub organization.", http.StatusNotFound)
			return
		}

		repoUrl := strings.Replace(repo.TFSProjectReference, "https://", "", -1)
		repoUrlSub := strings.Split(repoUrl, "/")

		isInnersource := strings.EqualFold(repoUrlSub[1], os.Getenv("GH_ORG_INNERSOURCE"))

		isMember, err := ghAPI.IsOrganizationMember(os.Getenv("GH_TOKEN"), os.Getenv("GH_ORG_INNERSOURCE"), ghUser)
		if err != nil {
			logger.LogException(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if (isInnersource && isMember) || (!isInnersource) {
			_, err := ghAPI.AddCollaborator(repoUrlSub[1], repo.Name, ghUser, permission)
			if err != nil {
				logger.LogException(err)
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
						logger.LogException(err)
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}

					if len(rec) > 0 {
						err := db.DeleteRepoOwnerRecordByUserAndProjectId(id, users[0]["UserPrincipalName"].(string))
						if err != nil {
							logger.LogException(err)
							http.Error(w, err.Error(), http.StatusInternalServerError)
							return
						}
					}
				}
			}

			w.WriteHeader(http.StatusOK)
		} else {
			logger.LogTrace("Can't invite a user that is not a member of the innersource organization.", contracts.Error)
			http.Error(w, "Can't invite a user that is not a member of the innersource organization.", http.StatusInternalServerError)
			return
		}
	}
}

func RemoveCollaborator(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	req := mux.Vars(r)
	id, _ := strconv.ParseInt(req["id"], 10, 64)
	ghUser := req["ghUser"]
	permission := req["permission"]

	// Get repository
	data := db.GetProjectById(id)
	s, err := json.Marshal(data)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var repoList []RepoDto
	err = json.Unmarshal(s, &repoList)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(repoList) > 0 {
		repo := repoList[0]

		if repo.TFSProjectReference == "" {
			logger.LogTrace("Repository doesn't exists on GitHub organization.", contracts.Error)
			http.Error(w, "Repository doesn't exists on GitHub organization.", http.StatusNotFound)
			return
		}

		repoUrl := strings.Replace(repo.TFSProjectReference, "https://", "", -1)
		repoUrlSub := strings.Split(repoUrl, "/")

		_, err := ghAPI.RemoveCollaborator(repoUrlSub[1], repo.Name, ghUser, permission)
		if err != nil {
			logger.LogException(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if permission == "admin" {
			users, _ := db.GetUserByGitHubUsername(ghUser)

			if len(users) > 0 {
				err = db.DeleteRepoOwnerRecordByUserAndProjectId(id, users[0]["UserPrincipalName"].(string))
				if err != nil {
					logger.LogException(err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}

		w.WriteHeader(http.StatusOK)
	}
}

func RepoOwnersCleanup(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	logger.TrackTrace("REPO OWNERS CLEANUP TRIGGERED", contracts.Information)
	token := os.Getenv("GH_TOKEN")

	// Get all repos from database
	repos, err := db.GetGitHubRepositories()
	if err != nil {
		logger.LogException(err)
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
			cleanupRepoOwners(r, token, logger)
			<-guard
			wg.Done()
		}(repo)
	}
	wg.Wait()
	w.WriteHeader(http.StatusOK)
	logger.TrackTrace("REPO OWNERS CLEANUP SUCCESSFUL", contracts.Information)
}

func RecurringApproval(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	const IN_REVIEW = 2

	projectApprovals, err := db.GetProjectApprovalsByStatusId(IN_REVIEW)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, projectApproval := range projectApprovals {
		created := projectApproval["Created"].(time.Time)

		daysSinceCreation := time.Since(created).Hours() / 24

		daysSinceCreationFloor := math.Floor(daysSinceCreation)

		projectApprovalApprovers, err := db.GetApprovalRequestApproversByApprovalRequestId(int(projectApproval["ProjectApprovalId"].(int64)))
		if err != nil {
			logger.LogException(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			continue
		}

		var recipients []string

		for _, projectApprovalApprover := range projectApprovalApprovers {
			recipients = append(recipients, projectApprovalApprover.ApproverEmail)
		}

		if int(daysSinceCreationFloor)%7 == 0 && daysSinceCreationFloor != 0 {
			messageBody := notification.RepositoryPublicApprovalRemainderMessageBody{
				Recipients:   recipients,
				ApprovalLink: fmt.Sprintf("%s/response/%s/%s/%s/1", os.Getenv("APPROVAL_SYSTEM_APP_URL"), os.Getenv("APPROVAL_SYSTEM_APP_ID"), os.Getenv("APPROVAL_SYSTEM_APP_MODULE_PROJECTS"), projectApproval["ItemId"].(string)),
				ApprovalType: projectApproval["ApprovalType"].(string),
				RepoLink:     projectApproval["RepoLink"].(string),
				RepoName:     projectApproval["RepoName"].(string),
				UserName:     projectApproval["Requester"].(string),
			}
			err = messageBody.Send()
			if err != nil {
				logger.LogException(err)
			}
		}
	}
}

func GetRepoCollaborators(org string, repo string, role string, affiliations string) []*github.User {

	token := os.Getenv("GH_TOKEN")

	repoCollabs := ghAPI.RepositoriesListCollaborators(token, org, repo, role, affiliations)

	return repoCollabs
}

func cleanupRepoOwners(repo map[string]interface{}, token string, logger *appinsights_wrapper.TelemetryClient) {
	logger.TrackTrace("Checking owners of : "+repo["Name"].(string), contracts.Information)

	if repo["TFSProjectReference"] == nil {
		// Get owners of the repo on the database
		owners, err := db.GetRepoOwnersByProjectIdWithGHUsername(repo["Id"].(int64))
		if err != nil {
			logger.LogException(err)
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
		logger.LogException(err)
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

func RemoveRepoIfNotExist(projectId int, repoName string, isGithubIdNil bool, logger *appinsights_wrapper.TelemetryClient) bool {
	isExist, err := ghAPI.IsRepoExisting(repoName)
	if err != nil {
		logger.LogTrace(fmt.Sprintf(err.Error(), " | REPO NAME : ", repoName), contracts.Error)
		return false
	}

	if !isExist || isGithubIdNil {
		err := db.DeleteProjectById(projectId)
		if err != nil {
			logger.LogException(err)
			return false
		}
		return true
	}
	return false
}

func indexRepo(repo ghAPI.Repo, logger *appinsights_wrapper.TelemetryClient) {
	logger.TrackTrace("Indexing "+repo.Name+"...", contracts.Information)

	visibilityId := 3
	if repo.Visibility == "private" {
		visibilityId = 1
	} else if repo.Visibility == "internal" {
		visibilityId = 2
	}

	param := map[string]interface{}{
		"GithubId":            repo.GithubId,
		"Name":                repo.Name,
		"AssetCode":           repo.Name,
		"Description":         repo.Description,
		"Organization":        repo.Org,
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
			logger.LogException(err)
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
				logger.LogException(err)
				return
			}

			//Insert to repoowners table
			if len(users) > 0 {
				if len(users) > 0 {
					err = db.RepoOwnersInsert(project[0]["Id"].(int64), users[0]["UserPrincipalName"].(string))
					if err != nil {
						logger.LogException(err)
					}
				}
			}
		}
	} else {
		err := db.ProjectInsertByImport(param)
		if err != nil {
			logger.LogException(err)
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
				logger.LogException(err)
				return
			}

			//Insert to repoowners table
			if len(users) > 0 {
				if len(users) > 0 {
					err = db.RepoOwnersInsert(project[0]["Id"].(int64), users[0]["UserPrincipalName"].(string))
					if err != nil {
						logger.LogException(err)
						return
					}
				}
			}
		}
	}
	if len(repo.Topics) > 0 {
		err := db.DeleteProjectTopics(projectId)
		if err != nil {
			logger.LogException(err)
			return
		}
		for i := 0; i < len(repo.Topics); i++ {
			err := db.InsertProjectTopics(projectId, repo.Topics[i])
			if err != nil {
				logger.LogException(err)
				return
			}
		}
	}
}

func RequestApproval(projectId int64, email string, logger *appinsights_wrapper.TelemetryClient) {
	projectApprovals, err := db.RequestProjectApprovals(projectId, email)
	if err != nil {
		logger.LogException(err)
		return
	}

	for _, v := range projectApprovals {
		if len(v.Approvers) == 0 {
			continue
		}
		if v.RequestStatus == "New" {
			err := ApprovalSystemRequest(v, logger)
			if err != nil {
				logger.LogException("ID:" + strconv.FormatInt(int64(v.Id), 10) + " " + err.Error())
				return
			}
		}
	}
}

func ApprovalSystemRequest(data db.ProjectApprovalApprovers, logger *appinsights_wrapper.TelemetryClient) error {
	url := os.Getenv("APPROVAL_SYSTEM_APP_URL")
	if url != "" {
		url = url + "/api/request"
		ch := make(chan *http.Response)
		// var res *http.Response

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
									Request for |ApprovalType| Review
									</td>
								</tr>
							</table>
						</th>
					</tr>
					<tr>
						<td class="center-table"  align="center">
							<table style="width: 100%; max-width: 700px;" class="margin-auto">
								<tr>
									<td style="padding-top: 20px">
										<p>Hi,</p>
										<p>|RequesterName| is requesting for a repository to be public and is now pending for |ApprovalType| review.</p>
										<p>Below are the details:</p>
									</td>
								</tr>
							</table>
						</td>
					</tr>
					<tr>
						<th class="center-table">
							<table style="width: 100%; max-width: 700px; margin: auto">
								<tr class="border-top">
									<td style="font-size: 14px; padding-top: 15px; font-weight: 600;">
										Repository Name
									</td>
									<td style="font-size: 14px; padding-top: 15px; font-weight: 400;">
										|ProjectName|
									</td>
								</tr>
								<tr class="border-top">
									<td style="font-size: 14px; padding-top: 15px; font-weight: 600;">
										Requested by
									</td>
									<td style="font-size: 14px; padding-top: 15px; font-weight: 400;">
										|Requester|
									</td>
								</tr>
								<tr class="border-top">
									<td style="font-size: 14px; padding-top: 15px; font-weight: 600;">
										Description
									</td>
									<td style="font-size: 14px; padding-top: 15px; font-weight: 400;">
										|ProjectDescription|
									</td>
								</tr>
								<tr class="border-top">
									<td style="font-size: 14px; padding-top: 15px; font-weight: 600;">
										Is this a new contribution with no prior code development? <br>
										(i.e., no existing Avanade IP, no third-party/OSS code, etc.)
									</td>
									<td style="font-size: 14px; padding-top: 15px; font-weight: 400;">
										|Newcontribution|
									</td>
								</tr>
								<tr class="border-top">
									<td style="font-size: 14px; padding-top: 15px; font-weight: 600;">
										Who is sponsoring thapprovalsyscois OSS contribution?
									</td>
									<td style="font-size: 14px; padding-top: 15px; font-weight: 400;">
										|OSSsponsor|
									</td>
								</tr>
								<tr class="border-top">
									<td style="font-size: 14px; padding-top: 15px; font-weight: 600;">
										Will Avanade use this contribution in client accounts <br>
										and/or as part of an Avanade offerings/assets?
									</td>
									<td style="font-size: 14px; padding-top: 15px; font-weight: 400;">
										|Avanadeofferingsassets|
									</td>
								</tr>
								<tr class="border-top">
									<td style="font-size: 14px; padding-top: 15px; font-weight: 600;">
										Will there be a commercial version of this contribution?
									</td>
									<td style="font-size: 14px; padding-top: 15px; font-weight: 400;">
										|Willbecommercialversion|
									</td>
								</tr>
								<tr class="border-top">
									<td style="font-size: 14px; padding-top: 15px; font-weight: 600;">
										Additional OSS Contribution Information <br>
										(e.g. planned maintenance/support, etc.)?
									</td>
									<td style="font-size: 14px; padding-top: 15px; font-weight: 400;">
										|OSSContributionInformation|
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
									For more information, send an email to <a href="mailto:|RequesterUserPrincipalName|">|RequesterUserPrincipalName|</a>
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
			"|RequesterName|", data.RequesterName,
			"|ApprovalType|", data.ApprovalType,
			"|ProjectName|", data.ProjectName,
			"|Requester|", data.RequesterName,
			"|ProjectDescription|", data.ProjectDescription,
			"|RequesterUserPrincipalName|", data.RequesterUserPrincipalName,
			"|Newcontribution|", data.NewContribution,
			"|OSSsponsor|", data.OSSsponsor,
			"|Avanadeofferingsassets|", data.OfferingsAssets,
			"|Willbecommercialversion|", data.WillBeCommercialVersion,
			"|OSSContributionInformation|", data.OSSContributionInformation,
		)
		body := replacer.Replace(bodyTemplate)
		postParams := ProjectApprovalSystemPostDto{
			ApplicationId:       os.Getenv("APPROVAL_SYSTEM_APP_ID"),
			ApplicationModuleId: os.Getenv("APPROVAL_SYSTEM_APP_MODULE_PROJECTS"),
			Subject:             fmt.Sprintf("Request for %v Review - %v", data.ApprovalType, data.ProjectName),
			Body:                body,
			Emails:              data.Approvers,
			RequesterEmail:      data.RequesterUserPrincipalName,
		}

		go getHttpPostResponseStatus(url, postParams, ch, logger)
		r := <-ch
		if r != nil {
			var res ProjectApprovalSystemPostResponseDto
			err := json.NewDecoder(r.Body).Decode(&res)
			if err != nil {
				logger.LogException(err)
				return err
			}

			messageBody := notification.RepositoryPublicApprovalMessageBody{
				Recipients:   []string{},
				ApprovalLink: fmt.Sprintf("%s/response/%s/%s/%s/1", os.Getenv("APPROVAL_SYSTEM_APP_URL"), postParams.ApplicationId, postParams.ApplicationModuleId, res.ItemId),
				ApprovalType: data.ApprovalType,
				RepoLink:     fmt.Sprintf("https://github.com/" + os.Getenv("GH_ORG_INNERSOURCE") + "/" + data.ProjectName),
				RepoName:     data.ProjectName,
				UserName:     data.RequesterName,
			}
			err = messageBody.Send()
			if err != nil {
				logger.LogException(err)
			}

			db.ProjectsApprovalUpdateGUID(data.Id, res.ItemId)
		}
	}
	return nil
}

func getHttpPostResponseStatus(url string, data interface{}, ch chan *http.Response, logger *appinsights_wrapper.TelemetryClient) {
	jsonReq, err := json.Marshal(data)
	if err != nil {
		logger.LogException(err)
		ch <- nil
	}
	res, err := http.Post(url, "application/json; charset=utf-8", bytes.NewBuffer(jsonReq))
	if err != nil {
		ch <- nil
	}
	ch <- res
}

func ReprocessRequestApproval() {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	projectApprovals, err := db.ReprocessFailedProjectApprovals()
	if err != nil {
		logger.LogException(err)
		return
	}

	for _, v := range projectApprovals {
		if len(v.Approvers) == 0 {
			continue
		}
		err := ApprovalSystemRequest(v, logger)
		if err != nil {
			logger.LogException(err)
			return
		}
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

func AddCollaboratorToRequestedRepo(user string, repo string, repoId int64, logger *appinsights_wrapper.TelemetryClient) (*github.Response, error) {
	innersource := os.Getenv("GH_ORG_INNERSOURCE")
	ghUser := db.Users_Get_GHUser(user)

	isInnersourceMember, err := ghAPI.IsOrganizationMember(os.Getenv("GH_TOKEN"), os.Getenv("GH_ORG_INNERSOURCE"), ghUser)
	if err != nil {
		logger.LogException(err)
		return nil, err
	}

	var resp *github.Response
	if isInnersourceMember {
		resp, err = ghAPI.AddCollaborator(innersource, repo, ghUser, "admin")
		if err != nil {
			logger.LogException(err)
			return resp, err
		}
		err = db.RepoOwnersInsert(repoId, user)
		if err != nil {
			logger.LogException(err)
			return resp, err
		}
	}
	return resp, nil
}

func GetPopularTopics(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	params := r.URL.Query()
	offset, err := strconv.Atoi(params["offset"][0])
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rowCount, err := strconv.Atoi(params["rowCount"][0])
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	result, err := db.GetPopularTopics(offset, rowCount)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResp, err := json.Marshal(result)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

func DownloadProjectApprovalsByUsername(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	vars := mux.Vars(r)

	username := vars["username"]

	projectApprovals, err := db.GetProjectApprovalsByUsername(username)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var data [][]string

	t := reflect.TypeOf(db.ApprovalRequest{})

	columns := make([]string, t.NumField())
	for i := range columns {
		columns[i] = t.Field(i).Name
	}

	data = append(data, columns)

	for _, projectApproval := range projectApprovals {
		v := reflect.ValueOf(projectApproval)

		row := make([]string, v.NumField())

		for i := 0; i < v.NumField(); i++ {
			vi := v.Field(i).Interface()
			if vi == "<nil>" {
				row[i] = ""
				continue
			}
			row[i] = fmt.Sprintf("%v", v.Field(i).Interface())
		}

		fmt.Println(row)

		data = append(data, row)
	}

	w.Header().Set("Content-Disposition", "attachment; filename=project_approval_requests.csv")
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Transfer-Encoding", "chunked")

	writer := csv.NewWriter(w)
	defer writer.Flush()

	for _, value := range data {
		if err := writer.Write(value); err != nil {
			logger.LogException(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func HttpResponseError(w http.ResponseWriter, code int, errorMessage string, logger *appinsights_wrapper.TelemetryClient) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(code)
	response, err := json.Marshal(errorMessage)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

// Deleted repositories from index
func emailAdminDeletedProjects(to string, repos []string, logger *appinsights_wrapper.TelemetryClient) {
	repoList := "</p> <table  >"
	for _, repo := range repos {
		repoList = repoList + " <tr> <td>" + repo + " </td></tr>"
	}
	repoList = repoList + " </table  > <p>"

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

										<p>The following repositories were removed from the database as they no longer exist:</p>
										|RepoList|
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
		"|RepoList|", repoList,
	)

	body := replacer.Replace(bodyTemplate)

	m := email.Message{
		Subject: "List of Deleted Repo",
		Body: email.Body{
			Content: body,
			Type:    email.HtmlMessageType,
		},
		ToRecipients: []email.Recipient{
			{
				Email: to,
			},
		},
	}

	err := email.SendEmail(m, false)
	if err != nil {
		logger.LogException(err)
	}
}

// List of users converted into outside collaborators to Repo Owner
func EmailAdminConvertToColaborator(to string, outisideCollab []string, logger *appinsights_wrapper.TelemetryClient) {
	e := time.Now()
	collabList := "</p> <table  >"
	for _, collab := range outisideCollab {
		collabList = collabList + " <tr> <td>" + collab + " </td></tr>"
	}
	collabList = collabList + " </table  > <p>"

	var message string
	if len(collabList) == 1 {
		message = fmt.Sprintf("<p>This is to inform you that %d GitHub user on %s was converted as an outside collaborator.</p>", len(outisideCollab), os.Getenv("GH_ORG_OPENSOURCE"))
	} else {
		message = fmt.Sprintf("<p>This is to inform you that %d GitHub users on %s were converted as outside collaborators.</p>", len(outisideCollab), os.Getenv("GH_ORG_OPENSOURCE"))
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
						<td class="center-table" align="center">
							<table style="width: 100%; max-width: 700px;" class="margin-auto">
								<tr>
									<td style="padding-top: 20px">
										<p>Hello |Admin| ,  </p>
										|Message|
										|ConvertedList|
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
		"|Admin|", to,
		"|Message", message,
		"|ConvertedList|", collabList,
	)

	body := replacer.Replace(bodyTemplate)

	m := email.Message{
		Subject: "GitHub Organization Scan",
		Body: email.Body{
			Content: body,
			Type:    email.HtmlMessageType,
		},
		ToRecipients: []email.Recipient{
			{
				Email: to,
			},
		},
	}

	email.SendEmail(m, false)
	logger.TrackTrace(fmt.Sprintf("GitHub User was converted into an outside  on %s was sent.", e), contracts.Information)
}

// List of users converted into outside collaborators to OSPO
func EmailRepoAdminConvertToColaborator(to string, repoName string, outisideCollab []string, logger *appinsights_wrapper.TelemetryClient) {
	e := time.Now()
	link := "https://github.com/" + os.Getenv("GH_ORG_OPENSOURCE") + "/" + repoName
	link = "<a href=\"" + link + "\">" + repoName + "</a>"
	collabList := "</p> <table  >"
	for _, collab := range outisideCollab {
		collabList = collabList + " <tr> <td>" + collab + " </td></tr>"
	}

	var message string
	collabList = collabList + " </table  > <p>"
	if len(outisideCollab) == 1 {
		message = fmt.Sprintf("<p>This is to inform you that <b> %d </b> GitHub user on your GitHub repo %s was converted to an outside collaborator. </p>", len(outisideCollab), link)
	} else {
		message = fmt.Sprintf("<p>This is to inform you that <b> %d </b> GitHub users on your GitHub repo %s were converted to outside collaborators. </p>", len(outisideCollab), link)
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
						<td class="center-table" align="center">
							<table style="width: 100%; max-width: 700px;" class="margin-auto">
								<tr>
									<td style="padding-top: 20px">
										<p>Hello |Admin| ,  </p>
										|Message|
										|ConvertedList|
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
		"|Admin|", to,
		"|Message", message,
		"|ConvertedList|", collabList,
	)

	body := replacer.Replace(bodyTemplate)

	m := email.Message{
		Subject: "GitHub Organization Scan",
		Body: email.Body{
			Content: body,
			Type:    email.HtmlMessageType,
		},
		ToRecipients: []email.Recipient{
			{
				Email: to,
			},
		},
	}

	email.SendEmail(m, true)
	logger.TrackTrace(fmt.Sprintf("GitHub User was converted into an outside  on %s was sent.", e), contracts.Information)
}

// List of repos with less than 2 owners to OSPO
func EmailOspoOwnerDeficient(to string, org string, repoName []string, logger *appinsights_wrapper.TelemetryClient) {
	e := time.Now()
	var link string

	repoNameList := "</p> <table  >"
	for _, repo := range repoName {
		link = "https://github.com/" + org + "/" + repo + "/settings/access"
		link = "<a href=\"" + link + "\">" + repo + "</a>"
		repoNameList = repoNameList + " <tr> <td>" + link + " </td></tr>"
	}

	repoNameList = repoNameList + " </table  > <p>"

	var message string
	if len(repoName) == 1 {
		message = fmt.Sprintf("<p>This is to inform you that <b> %d </b> repository on %s needs to add a co-owner.", len(repoName), org)

	} else {
		message = fmt.Sprintf("<p>This is to inform you that <b> %d </b> repositories on %s need to add a co-owner.", len(repoName), org)
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
									Activity - Request for Help
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
										|Message|
										|RepoList|
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
		"|Admin|", to,
		"|Message", message,
		"|RepoList|", repoNameList,
	)

	body := replacer.Replace(bodyTemplate)

	m := email.Message{
		Subject: "Repository Owners Scan",
		Body: email.Body{
			Content: body,
			Type:    email.HtmlMessageType,
		},
		ToRecipients: []email.Recipient{
			{
				Email: to,
			},
		},
	}

	email.SendEmail(m, false)
	logger.TrackTrace(fmt.Sprintf(" less than 2 owner    %s was sent to OSPO.", e), contracts.Information)
}

// List of repos with less than 2 owners to repo owner
func EmailcoownerDeficient(to string, Org string, reponame string) {
	var body string
	var link string
	link = "https://github.com/" + Org + "/" + reponame + "/settings/access"
	link = "<a href=\"" + link + "\"> here </a>"

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
										<p>This is to inform you that you are the only admin on |RepoName|  GitHub repository. We recommend at least 2 admins on each repository.
										<p>Click |Link| to add a co-owner.  </p>
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
		"|Admin|", to,
		"|RepoName", reponame,
		"|Link|", link,
	)

	body = replacer.Replace(bodyTemplate)

	m := email.Message{
		Subject: "Repository Owners Scan",
		Body: email.Body{
			Content: body,
			Type:    email.HtmlMessageType,
		},
		ToRecipients: []email.Recipient{
			{
				Email: to,
			},
		},
	}

	email.SendEmail(m, true)
}

func ValidateOrgMembers(org, repo, newOrg string, logger *appinsights_wrapper.TelemetryClient) (isSuccessful bool) {
	isSuccessful = true
	// GET ALL MEMBERS OF THE REPO
	collaborators := ghAPI.RepositoriesListCollaborators(os.Getenv("GH_TOKEN"), org, repo, "", "")

	// CHECK EACH COLLABORATORS OF THE REPO IF THEY ARE MEMBER OF THE NEW ORG
	// IF NOT REMOVE THEM FROM REPO
	for _, collaborator := range collaborators {
		isMember, err := ghAPI.IsOrganizationMember(os.Getenv("GH_TOKEN"), newOrg, collaborator.GetLogin())
		if err != nil {
			isSuccessful = false
			logger.LogException(err)
			continue
		}

		if !isMember {
			_, err := ghAPI.RemoveCollaborator(org, repo, collaborator.GetLogin(), "")
			if err != nil {
				isSuccessful = false
				logger.LogException(err)
				continue
			}
		}
	}

	return
}
