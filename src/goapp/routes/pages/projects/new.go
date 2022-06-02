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
)

func ProjectsNewHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		users := db.GetUsersWithGithub()
		template.UseTemplate(&w, r, "projects/new", users)
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

		nameCheck := make(chan bool, 2)

		go func() { nameCheck <- ghmgmtdb.Projects_IsExisting(body) }()
		go func() { b, _ := githubAPI.Repo_IsExisting(body.Name); nameCheck <- b }()

		if <-nameCheck || <-nameCheck {
			http.Error(w, "Project already exists.", http.StatusBadRequest)
		} else {
			_, err = githubAPI.CreatePrivateGitHubRepository(body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			id := ghmgmtdb.PRProjectsInsert(body, username.(string))
			RequestApproval(id)
		}

	}
}
func RequestApproval(id int64) {
	projectApprovals := ghmgmtdb.PopulateProjectsApproval(id)

	for _, v := range projectApprovals {
		ApprovalSystemRequest(v)
	}

}

func ApprovalSystemRequest(data models.TypProjectApprovals) {

	url := os.Getenv("APPROVAL_SYSTEM_APP_URL")
	if url != "" {
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
		<p>For more information, send an email to <a href="mailto:|RequesterUserPrincipalName|">|RequesterUserPrincipalName|</a></p>
		`
		replacer := strings.NewReplacer("|ApproverUserPrincipalName|", data.ApproverUserPrincipalName,
			"|RequesterName|", data.RequesterName,
			"|ApprovalType|", data.ApprovalType,
			"|ProjectName|", data.ProjectName,
			"|CoownerName|", data.CoownerName,
			"|ProjectDescription|", data.ProjectDescription,
			"|RequesterUserPrincipalName|", data.RequesterUserPrincipalName,
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

		var res models.TypApprovalSystemPostResponse
		err := json.NewDecoder(r.Body).Decode(&res)
		handleError(err)

		ghmgmtdb.ProjectsApprovalUpdateGUID(data.Id, res.ItemId)
	}

}

func getHttpPostResponseStatus(url string, data interface{}, ch chan *http.Response) {
	jsonReq, err := json.Marshal(data)
	res, err := http.Post(url, "application/json; charset=utf-8", bytes.NewBuffer(jsonReq))
	handleError(err)
	ch <- res
}

func handleError(err error) {
	if err != nil {
		fmt.Printf("ERROR: %+v", err)
	}
}
