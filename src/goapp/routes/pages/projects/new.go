package routes

import (
	"encoding/json"
	models "main/models"
	ghmgmtdb "main/pkg/ghmgmtdb"
	githubAPI "main/pkg/github"
	session "main/pkg/session"
	template "main/pkg/template"
	"net/http"
	db "main/pkg/ghmgmtdb"
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
			ghmgmtdb.PRProjectsInsert(body, username.(string))
		}
	}
}