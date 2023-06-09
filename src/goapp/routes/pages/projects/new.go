package routes

import (
	"encoding/json"
	"log"
	"main/models"
	"net/http"
	"os"
	"regexp"
	"strings"

	db "main/pkg/ghmgmtdb"
	ghAPI "main/pkg/github"
	"main/pkg/session"
	"main/pkg/template"

	"github.com/gorilla/mux"
)

type TypErrorJsonReturnDto struct {
	Error string `json:"error"`
}

func ProjectsNewHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		req := mux.Vars(r)
		id := req["id"]

		sessionaz, _ := session.Store.Get(r, "auth-session")
		sessiongh, _ := session.GetGitHubUserData(w, r)
		iprofile := sessionaz.Values["profile"]
		profile := iprofile.(map[string]interface{})
		username := profile["preferred_username"]
		isInnersourceMember, isOpensourceMember, _ := ghAPI.OrganizationsIsMember(os.Getenv("GH_TOKEN"), sessiongh.Username)

		users := db.GetUsersWithGithub()
		data := map[string]interface{}{
			"Id":                  id,
			"users":               users,
			"email":               username,
			"isInnersourceMember": isInnersourceMember,
			"isOpensourceMember":  isOpensourceMember,
			"innersourceOrg":      os.Getenv("GH_ORG_INNERSOURCE"),
			"opensourceOrg":       os.Getenv("GH_ORG_OPENSOURCE"),
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
		go func() { checkDB <- db.Projects_IsExisting(body) }()
		go func() { b, _ := ghAPI.Repo_IsExisting(dashedProjName); checkGH <- b }()

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
			isOrgAllowInternalRepo, err := ghAPI.IsOrgAllowInternalRepo()
			if err != nil {
				HttpResponseError(w, http.StatusBadRequest, "There is a problem checking if the organization is enterprise or not.")
				return
			}

			repo, errRepo := ghAPI.CreatePrivateGitHubRepository(body, username.(string))
			if errRepo != nil {
				log.Println(errRepo.Error())
				HttpResponseError(w, http.StatusInternalServerError, "There is a problem creating the GitHub repository.")
				return
			}
			body.GithubId = repo.GetID()
			body.TFSProjectReference = repo.GetHTMLURL()
			body.Visibility = 1

			if isOrgAllowInternalRepo {
				innersource := os.Getenv("GH_ORG_INNERSOURCE")
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
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		db.PRProjectsUpdate(body, username.(string))

		w.WriteHeader(http.StatusOK)
	}
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

func IsRepoNameValid(value string) bool {
	regex := regexp.MustCompile(`^([a-zA-Z0-9\-\_]|\.{3,}|\.{1,}[a-zA-Z0-9\-\_])([a-zA-Z0-9\-\_\.]+)?`)
	matched := regex.FindAllString(value, 1)

	if matched == nil {
		return false
	}

	return matched[0] == value
}
