package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	models "main/models"
	ghmgmt "main/pkg/ghmgmtdb"
	gh "main/pkg/github"
	githubAPI "main/pkg/github"
	session "main/pkg/session"
	"main/pkg/sql"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

func GetUserProjects(w http.ResponseWriter, r *http.Request) {
	// Get email address of the user
	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := profile["preferred_username"]

	// Connect to database
	dbConnectionParam := sql.ConnectionParam{
		ConnectionString: os.Getenv("GHMGMTDB_CONNECTION_STRING"),
	}

	db, err := sql.Init(dbConnectionParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Get project list
	params := make(map[string]interface{})
	params["UserPrincipalName"] = username
	projects, err := db.ExecuteStoredProcedureWithResult("PR_Projects_Select_ByUserPrincipalName", params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s, _ := json.Marshal(projects)
	var list []Repo
	err = json.Unmarshal(s, &list)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(list)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func GetUsersWithGithub(w http.ResponseWriter, r *http.Request) {

	users := ghmgmt.GetUsersWithGithub()

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func GetRequestStatusByProject(w http.ResponseWriter, r *http.Request) {
	req := mux.Vars(r)
	id := req["id"]

	// Connect to database
	dbConnectionParam := sql.ConnectionParam{
		ConnectionString: os.Getenv("GHMGMTDB_CONNECTION_STRING"),
	}

	db, err := sql.Init(dbConnectionParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Get project list
	params := make(map[string]interface{})
	params["Id"] = id
	projects, err := db.ExecuteStoredProcedureWithResult("PR_ProjectApprovals_Select_ById", params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(projects)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func GetRepoCollaboratorsByRepoId(w http.ResponseWriter, r *http.Request) {
	req := mux.Vars(r)
	id, _ := strconv.ParseInt(req["id"], 10, 64)

	// Get repository
	data := ghmgmt.GetProjectById(id)
	s, _ := json.Marshal(data)
	var repoList []Repo
	err := json.Unmarshal(s, &repoList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var result []Collaborator

	if len(repoList) > 0 {
		repo := repoList[0]

		if repo.RepositorySource == "GitHub" {
			repoUrl := strings.Replace(repo.TFSProjectReference, "https://", "", -1)
			repoUrlSub := strings.Split(repoUrl, "/")

			token := os.Getenv("GH_TOKEN")

			collaborators := gh.RepositoriesListCollaborators(token, repoUrlSub[1], repo.Name, "", "direct")
			outsideCollaborators := gh.RepositoriesListCollaborators(token, repoUrlSub[1], repo.Name, "", "outside")
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

				collabResult := Collaborator{
					Id:                    collaborator.GetID(),
					Role:                  *collaborator.RoleName,
					GitHubUsername:        *collaborator.Login,
					IsOutsideCollaborator: isOutsideCollab,
				}

				//Get user name and email address
				users, err := ghmgmt.GetUserByGitHubId(strconv.FormatInt(*collaborator.ID, 10))
				if err != nil {
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
			var collabResult Collaborator

			repoOwner, err := ghmgmt.GetRepoOwnersRecordByRepoId(int64(repo.Id))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if len(repoOwner) > 0 {
				users, err := ghmgmt.GetUserByUserPrincipal(repoOwner[0].UserPrincipalName)
				if err != nil {
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
		err := gh.ArchiveProject(project, archive == "1", organization)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		isArchived, err := gh.IsArchived(project, organization)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if isArchived {
			http.Error(w, "Repository is still archived. Unarchive the repo on GitHub and try again.", http.StatusBadRequest)
			return
		}
	}
	id, _ := strconv.ParseInt(projectId, 10, 64)
	ghmgmt.UpdateProjectIsArchived(id, archive == "1")
	w.WriteHeader(http.StatusOK)
}

func GetAllRepositories(w http.ResponseWriter, r *http.Request) {

	params := r.URL.Query()
	search := params["search"][0]
	offset, _ := strconv.Atoi(params["offset"][0])

	// Get repository list
	data := ghmgmt.Repos_Select_ByOffsetAndFilter(offset, search)
	s, _ := json.Marshal(data)
	var list []Repo
	err := json.Unmarshal(s, &list)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	result := RepositoryList{
		Data:  list,
		Total: ghmgmt.Repos_TotalCount_BySearchTerm(search),
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(result)
	if err != nil {
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

	if currentState == "Public" {
		// Set repo to desired visibility then move to innersource
		err := gh.SetProjectVisibility(project, desiredState, opensource)
		if err != nil {
			http.Error(w, "Failed to make the repository "+desiredState, http.StatusInternalServerError)
			return
		}

		gh.TransferRepository(project, opensource, innersource)

		time.Sleep(3 * time.Second)
		repoResp, _ := gh.GetRepository(project, innersource)
		ghmgmt.UpdateTFSProjectReferenceById(id, repoResp.GetHTMLURL())
	} else {
		// Set repo to desired visibility
		err := gh.SetProjectVisibility(project, desiredState, innersource)
		if err != nil {
			http.Error(w, "Failed to make the repository "+desiredState, http.StatusInternalServerError)
			return
		}
	}

	// Update database
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ghmgmt.UpdateProjectVisibilityId(id, int64(visibilityId))

	w.WriteHeader(http.StatusOK)
}

func RequestMakePublic(w http.ResponseWriter, r *http.Request) {
	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := profile["preferred_username"]
	r.ParseForm()

	var body models.TypeMakeProjectPublicReqBody

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ghmgmt.PRProjectsUpdateLegalQuestions(body, username.(string))

	id, err := strconv.ParseInt(body.Id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	RequestApproval(id)
}

func ImportReposToDatabase(w http.ResponseWriter, r *http.Request) {
	var allRepos []gh.Repo

	organizations := []string{os.Getenv("GH_ORG_INNERSOURCE"), os.Getenv("GH_ORG_OPENSOURCE")}

	for _, org := range organizations {
		repos, err := gh.GetRepositoriesFromOrganization(org)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if repos != nil {
			allRepos = append(allRepos, repos...)
		}
	}

	var wg sync.WaitGroup

	for _, repo := range allRepos {
		wg.Add(1)

		go func(repo gh.Repo) {
			defer wg.Done()
			e := ghmgmt.Projects_IsExisting(models.TypNewProjectReqBody{Name: repo.Name})

			if !e {
				visibilityId := 3
				if repo.Visibility == "private" {
					visibilityId = 1
				} else if repo.Visibility == "internal" {
					visibilityId = 2
				}

				param := map[string]interface{}{
					"GithubId":     repo.GithubId,
					"Name":         repo.Name,
					"Description":  repo.Description,
					"IsArchived":   repo.IsArchived,
					"VisibilityId": visibilityId,
					"Created":      repo.Created.Format("2006-01-02 15:04:05"),
				}

				err := ghmgmt.ProjectInsertByImport(param)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}

		}(repo)
	}
}

func InitIndexOrgRepos(w http.ResponseWriter, r *http.Request) {
	var repos []gh.Repo

	orgs := []string{os.Getenv("GH_ORG_INNERSOURCE"), os.Getenv("GH_ORG_OPENSOURCE")}

	for _, org := range orgs {
		reposByOrg, err := gh.GetRepositoriesFromOrganization(org)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if reposByOrg != nil {
			repos = append(repos, reposByOrg...)
		}
	}

	var wg sync.WaitGroup

	for _, repo := range repos {
		wg.Add(1)
		go func(repo gh.Repo) {
			defer wg.Done()

			visibilityId := 3
			if repo.Visibility == "private" {
				visibilityId = 1
			} else if repo.Visibility == "internal" {
				visibilityId = 2
			}

			param := map[string]interface{}{
				"GithubId":     repo.GithubId,
				"Name":         repo.Name,
				"Description":  repo.Description,
				"IsArchived":   repo.IsArchived,
				"VisibilityId": visibilityId,
				"Created":      repo.Created.Format("2006-01-02 15:04:05"),
			}

			isExisting := ghmgmt.Projects_IsExisting(models.TypNewProjectReqBody{Name: repo.Name})

			var err error

			if isExisting {
				param["Id"] = ghmgmt.GetProjectByName(repo.Name)[0]["Id"]
				err = ghmgmt.ProjectUpdateByImport(param)
			} else {
				err = ghmgmt.ProjectInsertByImport(param)
			}

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}(repo)
	}

	wg.Wait()
	w.WriteHeader(http.StatusOK)
}

func IndexOrgRepos(w http.ResponseWriter, r *http.Request) {
	var repos []gh.Repo

	orgs := []string{os.Getenv("GH_ORG_INNERSOURCE"), os.Getenv("GH_ORG_OPENSOURCE")}

	for _, org := range orgs {
		reposByOrg, err := gh.GetRepositoriesFromOrganization(org)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if reposByOrg != nil {
			repos = append(repos, reposByOrg...)
		}
	}

	var wg sync.WaitGroup

	for _, repo := range repos {
		wg.Add(1)
		go func(repo gh.Repo) {
			defer wg.Done()

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

			isExisting := ghmgmt.Projects_IsExisting_By_GithubId(models.TypNewProjectReqBody{GithubId: repo.GithubId})

			var err error

			if isExisting {
				userdata := ghmgmt.GetProjectByGithubId(repo.GithubId)
				param["Id"] = userdata[0]["Id"]

				err = ghmgmt.ProjectUpdateByImport(param)
				if userdata[0]["CreatedBy"] != nil {
					RepoOwners, _ := ghmgmt.RepoOwnersByUserAndProjectId(param["Id"].(int64), userdata[0]["CreatedBy"].(string))
					if len(RepoOwners) < 1 {
						error := ghmgmt.RepoOwnersInsert(param["Id"].(int64), userdata[0]["CreatedBy"].(string))
						if error != nil {

							fmt.Println(error)
						}
					}
				}
			} else {
				err = ghmgmt.ProjectInsertByImport(param)
			}

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}(repo)
	}

	wg.Wait()
	w.WriteHeader(http.StatusOK)
}

func RequestApproval(id int64) {
	projectApprovals := ghmgmt.PopulateProjectsApproval(id)

	for _, v := range projectApprovals {
		if v.RequestStatus == "New" {
			err := ApprovalSystemRequest(v)
			handleError(err)
		}
	}
}

func ApprovalSystemRequest(data models.TypProjectApprovals) error {

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
			"|Avanadeofferingsassets|", data.Avanadeofferingsassets,
			"|Willbecommercialversion|", data.Willbecommercialversion,
			"|OSSContributionInformation|", data.OSSContributionInformation,
		)
		body := replacer.Replace(bodyTemplate)
		postParams := models.TypApprovalSystemPost{
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
			var res models.TypApprovalSystemPostResponse
			err := json.NewDecoder(r.Body).Decode(&res)
			if err != nil {
				return err
			}

			ghmgmt.ProjectsApprovalUpdateGUID(data.Id, res.ItemId)
		}
	}
	return nil
}

func getHttpPostResponseStatus(url string, data interface{}, ch chan *http.Response) {
	jsonReq, err := json.Marshal(data)
	res, err := http.Post(url, "application/json; charset=utf-8", bytes.NewBuffer(jsonReq))
	if err != nil {
		ch <- nil
	}
	ch <- res
}

func handleError(err error) {
	if err != nil {
		fmt.Printf("ERROR: %+v", err)
	}
}

func ReprocessRequestApproval() {
	projectApprovals := ghmgmt.GetFailedProjectApprovalRequests()

	for _, v := range projectApprovals {
		go ApprovalSystemRequest(v)
	}
}

func RepoOwnersCleanup(w http.ResponseWriter, r *http.Request) {
	token := os.Getenv("GH_TOKEN")
	RepoUsers, _ := ghmgmt.SelectAllRepoNameAndOwners()
	organizations := [...]string{os.Getenv("GH_ORG_INNERSOURCE"), os.Getenv("GH_ORG_OPENSOURCE")}
	var wg sync.WaitGroup
	for _, RepoUser := range RepoUsers {
		wg.Add(1)
		func(RepoUser models.TypRepoOwner) {
			fmt.Println("Project ID : " + strconv.FormatInt(RepoUser.Id, 10))
			isAdmin := false
			validRepo := false
			for _, organization := range organizations {
				RepoAdmins_ := githubAPI.RepositoriesListCollaborators(token, organization, RepoUser.RepoName, "admin", "direct")
				if len(RepoAdmins_) > 0 {
					for _, list := range RepoAdmins_ {
						validRepo = true
						email, _ := ghmgmt.UsersGetEmail(*list.Login)
						if RepoUser.UserPrincipalName == email {
							isAdmin = true
						}
					}
				}
			}

			if validRepo && !isAdmin {
				err := ghmgmt.DeleteRepoOwnerRecordByUserAndProjectId(RepoUser.Id, RepoUser.UserPrincipalName)

				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}
			defer wg.Done()

		}(RepoUser)
	}
	wg.Wait()
}

type RepositoryList struct {
	Data  []Repo `json:"data"`
	Total int    `json:"total"`
}

type Repo struct {
	Id                     int    `json:"Id"`
	Name                   string `json:"Name"`
	Description            string `json:"Description"`
	IsArchived             bool   `json:"IsArchived"`
	Created                string `json:"Created"`
	RepositorySource       string `json:"RepositorySource"`
	TFSProjectReference    string `json:"TFSProjectReference"`
	Visibility             string `json:"Visibility"`
	ApprovalStatus         bool   `json:"ApprovalStatus"`
	ApprovalStatusId       int    `json:"ApprovalStatusId"`
	CoOwner                string `json:CoOwner`
	ConfirmAvaIP           bool   `json:ConfirmAvaIP`
	ConfirmEnabledSecurity bool   `json:ConfirmEnabledSecurity`
	CreatedBy              string `json:CreatedBy`
	Modified               string `json:Modified`
	ModifiedBy             string `json: ModifiedBy`
}

type Collaborator struct {
	Id                    int64  `json:"id"`
	Name                  string `json:"name"`
	Email                 string `json:"email"`
	Role                  string `json:"role"`
	GitHubUsername        string `json:"ghUsername`
	IsOutsideCollaborator bool   `json:"isOutsideCollaborator"`
}
