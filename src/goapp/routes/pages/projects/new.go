package routes

import (
	"encoding/json"
	models "main/models"
	githubAPI "main/pkg/github"
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
		var body models.TypNewProjectReqBody
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = githubAPI.CreatePrivateGitHubRepository(body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}
