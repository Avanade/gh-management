package routes

import (
	"encoding/json"
	gh "main/pkg/github"
	"net/http"
)

func GetAvanadeProjects(w http.ResponseWriter, r *http.Request) {

	repos, err := gh.GetRepositoriesFromOrganization("Avanade")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(repos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}
