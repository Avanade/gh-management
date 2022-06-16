package routes

import (
	"encoding/json"
	"main/pkg/envvar"
	ghmgmt "main/pkg/ghmgmtdb"
	gh "main/pkg/github"
	session "main/pkg/session"
	"net/http"
	"sort"
	"strings"
	"sync"

	"github.com/gorilla/mux"
)

func GetAvanadeProjects(w http.ResponseWriter, r *http.Request) {
	var allRepos []gh.Repo

	o := envvar.GetEnvVar("GH_ORGANIZATIONS", "Avanade")
	organizations := strings.Split(o, " ")

	for _, org := range organizations {
		repos, err := gh.GetRepositoriesFromOrganization(org)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if repos != nil {
			allRepos = append(allRepos, repos...)
		}
	}

	sort.Slice(allRepos[:], func(i, j int) bool {
		return strings.ToLower(allRepos[i].Name) < strings.ToLower(allRepos[j].Name)
	})

	var wg = &sync.WaitGroup{}

	for i, project := range allRepos {
		wg.Add(1)
		go func(i int, p gh.Repo) {
			rec := ghmgmt.GetProjectByName(p.Name)
			if len(rec) == 0 {
				p.IsArchived = false
			} else {
				allRepos[i].IsArchived = rec[0]["IsArchived"].(bool)
			}
			wg.Done()
		}(i, project)
	}

	wg.Wait()

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(allRepos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func SetVisibility(w http.ResponseWriter, r *http.Request) {
	// Check if user is an admin
	isAdmin, err := session.IsUserAdmin(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !isAdmin {
		http.Error(w, "Not enough privilege to do the action.", http.StatusForbidden)
		return
	}

	req := mux.Vars(r)
	project := req["project"]
	archive := req["archive"]
	private := req["private"]

	err = ghmgmt.UpdateIsArchiveIsPrivate(project, archive == "1", private == "1")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//If project is currently public, set visibility to private
	if private == "0" {
		visibility := "private"
		if private == "0" {
			visibility = "public"
		}
		err := gh.SetProjectVisibility(project, visibility)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
