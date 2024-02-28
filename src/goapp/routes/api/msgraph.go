package routes

import (
	"encoding/json"
	"main/pkg/appinsights_wrapper"
	db "main/pkg/ghmgmtdb"
	"main/pkg/msgraph"
	"net/http"

	"github.com/gorilla/mux"
)

func GetAllUserFromActiveDirectory(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	users, err := msgraph.GetAllUsers()
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(users)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func SearchUserFromActiveDirectory(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	vars := mux.Vars(r)
	search := vars["search"]

	users, err := msgraph.SearchUsers(search)

	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(users)
	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func IndexADGroups(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	groups, err := msgraph.GetADGroups()

	if err != nil {
		logger.LogException(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, group := range groups {
		hasGitHubAccess, err := msgraph.HasGitHubAccess(group.Id)
		if err != nil {
			logger.LogException(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if hasGitHubAccess {
			db.ADGroup_Insert(group.Id, group.Name)
		}
	}

	w.WriteHeader(http.StatusOK)
}
