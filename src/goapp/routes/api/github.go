package routes

import (
	"encoding/json"
	gh "main/pkg/github"
	"net/http"
)

func GetAvanadeProjects(w http.ResponseWriter, r *http.Request) {
	var allRepos []gh.Repo
	repos, err := gh.GetRepositoriesFromOrganization("ava-innersource")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if repos != nil {
		allRepos = append(allRepos, repos...)
	}

	repos, err = gh.GetRepositoriesFromOrganization("Avanade")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if repos != nil {
		allRepos = append(allRepos, repos...)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(allRepos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}
