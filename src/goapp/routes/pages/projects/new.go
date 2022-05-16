package routes

import (
	"encoding/json"
	models "main/models"
	githubAPI "main/pkg/github"
	template "main/pkg/template"
	"net/http"
)

func ProjectsNewHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		template.UseTemplate(&w, r, "projects/new", nil)
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
