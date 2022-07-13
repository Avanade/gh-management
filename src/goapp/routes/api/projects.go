package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	models "main/models"
	ghmgmt "main/pkg/ghmgmtdb"
	gh "main/pkg/github"
	session "main/pkg/session"
	"main/pkg/sql"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

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

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(projects)
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

func GetAvanadeProjects(w http.ResponseWriter, r *http.Request) {
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

	sort.Slice(allRepos[:], func(i, j int) bool {
		return strings.ToLower(allRepos[i].Name) < strings.ToLower(allRepos[j].Name)
	})

	// var wg = &sync.WaitGroup{}

	// for i, project := range allRepos {
	// 	wg.Add(1)
	// 	go func(i int, p gh.Repo) {
	// 		rec := ghmgmt.GetProjectByName(p.Name)
	// 		if len(rec) == 0 {
	// 			p.IsArchived = false
	// 		} else {
	// 			allRepos[i].IsArchived = rec[0]["IsArchived"].(bool)
	// 		}
	// 		wg.Done()
	// 	}(i, project)
	// }

	// wg.Wait()

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(allRepos)
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
	isArchived := req["isArchived"]
	visibilityId := 1 //public
	if desiredState == "internal" {
		visibilityId = 2 //internal
	}

	innersource := os.Getenv("GH_ORG_INNERSOURCE")
	opensource := os.Getenv("GH_ORG_OPENSOURCE")

	if currentState == "Public" {
		if isArchived == "1" {
			gh.ArchiveProject(project, false, opensource)
		}

		// Set repo to desired visibility then move to innersource
		err := gh.SetProjectVisibility(project, desiredState, opensource)
		if err != nil {
			http.Error(w, "Failed to make the repository "+desiredState, http.StatusInternalServerError)
			return
		}

		gh.TransferRepository(project, opensource, innersource)
	} else {
		if isArchived == "1" {
			gh.ArchiveProject(project, false, innersource)
		}
		// Set repo to desired visibility
		err := gh.SetProjectVisibility(project, desiredState, innersource)
		if err != nil {
			http.Error(w, "Failed to make the repository "+desiredState, http.StatusInternalServerError)
			return
		}
	}

	// Update database
	id, err := strconv.ParseInt(projectId, 10, 64)
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

	go RequestApproval(id)
}

func RequestApproval(id int64) {
	projectApprovals := ghmgmt.PopulateProjectsApproval(id)

	for _, v := range projectApprovals {
		err := ApprovalSystemRequest(v)
		handleError(err)
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
				<td style="font-weight: bold;">CoOwner<td>
				<td style="font-size:larger">|CoownerName|<td>
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
				<td style="font-weight: bold;">Who is sponsoring this OSS contribution?<td>
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
			"|CoownerName|", data.CoownerName,
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
