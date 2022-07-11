package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	models "main/models"
	db "main/pkg/ghmgmtdb"
	ghmgmtdb "main/pkg/ghmgmtdb"
	githubAPI "main/pkg/github"
	session "main/pkg/session"
	template "main/pkg/template"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	//"github.com/gorilla/mux"
)

func ProjectsNewHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":

		//users := db.GetUsersWithGithub()
		req := mux.Vars(r)
		id := req["id"]

		users := db.GetUsersWithGithub()
		data := map[string]interface{}{
			"Id":    id,
			"users": users,
		}
		template.UseTemplate(&w, r, "projects/new", data)
	case "POST":
		sessionaz, _ := session.Store.Get(r, "auth-session")
		iprofile := sessionaz.Values["profile"]
		profile := iprofile.(map[string]interface{})
		username := profile["preferred_username"]
		r.ParseForm()

		var body models.TypNewProjectReqBody

		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		checkDB := make(chan bool)
		checkGH := make(chan bool)

		var existsDb bool
		var existsGH bool
		dashedProjName := strings.ReplaceAll(body.Name, " ", "-")
		go func() { checkDB <- ghmgmtdb.Projects_IsExisting(body) }()
		go func() { b, _ := githubAPI.Repo_IsExisting(dashedProjName); checkGH <- b }()

		existsDb = <-checkDB
		existsGH = <-checkGH
		if existsDb || existsGH {
			if existsDb {
				httpResponseError(w, http.StatusBadRequest, "The project name is existing in the database.")
			}
			if existsGH {
				httpResponseError(w, http.StatusBadRequest, "The project name is existing in Github.")
			}
		} else {
			_, err = githubAPI.CreatePrivateGitHubRepository(body)
			if err != nil {
				fmt.Println(err)
				httpResponseError(w, http.StatusInternalServerError, "There is a problem creating the GitHub repository.")
			}

			id := ghmgmtdb.PRProjectsInsert(body, username.(string))

			go RequestApproval(id)
			w.WriteHeader(http.StatusOK)
		}
	}
}
func ProjectsHandler(w http.ResponseWriter, r *http.Request) {
	req := mux.Vars(r)
	id := req["id"]

	switch r.Method {
	case "GET":

		users := db.GetUsersWithGithub()
		data := map[string]interface{}{
			"Id":    id,
			"users": users,
		}
		template.UseTemplate(&w, r, "projects/new", data)
	case "POST":

		sessionaz, _ := session.Store.Get(r, "auth-session")
		iprofile := sessionaz.Values["profile"]
		profile := iprofile.(map[string]interface{})
		username := profile["preferred_username"]
		r.ParseForm()

		var body models.TypNewProjectReqBody

		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			fmt.Println(err.Error())
			return
		}

		ghmgmtdb.PRProjectsUpdate(body, username.(string))

		w.WriteHeader(http.StatusOK)
	}

}

func RequestApproval(id int64) {
	projectApprovals := ghmgmtdb.PopulateProjectsApproval(id)

	for _, v := range projectApprovals {
		err := ApprovalSystemRequest(v)
		handleError(err)
	}

}

func ReprocessRequestApproval() {
	projectApprovals := ghmgmtdb.GetFailedProjectApprovalRequests()

	for _, v := range projectApprovals {
		go ApprovalSystemRequest(v)

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

			ghmgmtdb.ProjectsApprovalUpdateGUID(data.Id, res.ItemId)
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

func httpResponseError(w http.ResponseWriter, code int, errorMessage string) {
	msg := TypErrorJsonReturn{
		Error: errorMessage,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	jsonResponse, err := json.Marshal(msg)
	handleError(err)
	w.Write(jsonResponse)
}

type TypErrorJsonReturn struct {
	Error string `json:"error"`
}
