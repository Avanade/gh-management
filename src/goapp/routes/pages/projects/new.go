package routes

import (
	"encoding/json"
	"fmt"
	models "main/models"
	ghmgmtdb "main/pkg/ghmgmtdb"
	githubAPI "main/pkg/github"
	session "main/pkg/session"
	template "main/pkg/template"
	"net/http"
)

func ProjectsNewHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		template.UseTemplate(&w, r, "projects/new", nil)
	case "POST":
		sessionaz, _ := session.Store.Get(r, "auth-session")
		iprofile := sessionaz.Values["profile"]
		profile := iprofile.(map[string]interface{})
		fmt.Println(profile["preferred_username"])
		username := profile["preferred_username"]
		r.ParseForm()

		var body models.TypNewProjectReqBody
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if ghmgmtdb.Projects_IsExisting(body) {
			http.Error(w, "Existing Project Name", http.StatusBadRequest)
		} else {
			ghmgmtdb.PRProjectsInsert(body, username.(string))
		}

		_, err = githubAPI.CreatePrivateGitHubRepository(body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	}
}
