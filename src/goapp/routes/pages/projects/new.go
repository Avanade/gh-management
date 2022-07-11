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
