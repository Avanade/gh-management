package routes

import (
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
	"regexp"
	"strings"

	"github.com/gorilla/mux"
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

		if !isRepoNameValid(body.Name) {
			httpResponseError(w, http.StatusBadRequest, "Invalid repository name.")
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
				return
			} else if existsGH {
				httpResponseError(w, http.StatusBadRequest, "The project name is existing in Github.")
				return
			}
		} else {
			isOrgAllowInternalRepo, err := githubAPI.GetIsOrgAllowInternalRepo()
			if err != nil {
				httpResponseError(w, http.StatusBadRequest, "There is a problem checking if the organization is enterprise or not.")
				return
			}

			repo, errRepo := githubAPI.CreatePrivateGitHubRepository(body, username.(string))
			if errRepo != nil {
				fmt.Println(errRepo)
				httpResponseError(w, http.StatusInternalServerError, "There is a problem creating the GitHub repository.")
				return
			}
			body.GithubId = repo.GetID()
			body.Visibility = 1

			if isOrgAllowInternalRepo {
				innersource := os.Getenv("GH_ORG_INNERSOURCE")
				err := githubAPI.SetProjectVisibility(repo.GetName(), "internal", innersource)
				if err != nil {
					return
				}
				body.Visibility = 2
			}

			_ = ghmgmtdb.PRProjectsInsert(body, username.(string))

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

func handleError(err error) {
	if err != nil {
		fmt.Printf("ERROR: %+v", err)
	}
}

func httpResponseError(w http.ResponseWriter, code int, errorMessage string) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(code)
	response, err := json.Marshal(errorMessage)
	handleError(err)
	w.Write(response)
}

func isRepoNameValid(value string) bool {
	regex := regexp.MustCompile(`^([a-zA-Z0-9\-\_]|\.{3,}|\.{1,}[a-zA-Z0-9\-\_])([a-zA-Z0-9\-\_\.]+)?`)
	matched := regex.FindAllString(value, 1)

	if matched == nil {
		return false
	}

	return matched[0] == value
}

type TypErrorJsonReturn struct {
	Error string `json:"error"`
}
