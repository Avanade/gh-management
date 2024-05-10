package routes

import (
	"encoding/json"
	"fmt"
	"main/pkg/appinsights_wrapper"
	db "main/pkg/ghmgmtdb"
	"main/pkg/msgraph"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/microsoft/ApplicationInsights-Go/appinsights/contracts"
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
	go func() {
		logger := appinsights_wrapper.NewClient()
		defer logger.EndOperation()

		logger.LogTrace("Pulling list of AD groups...", contracts.Information)
		groups, err := msgraph.GetADGroups()
		logger.LogTrace(fmt.Sprintf("%d AD groups found", len(groups)), contracts.Information)
		if err != nil {
			logger.LogException(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var wg sync.WaitGroup
		maxGoroutines := 100
		guard := make(chan struct{}, maxGoroutines)

		for _, repo := range groups {
			guard <- struct{}{}
			wg.Add(1)
			go func(g msgraph.ADGroup) {
				indexGroup(g, logger)
				<-guard
				wg.Done()
			}(repo)
		}
		wg.Wait()

	}()

	w.WriteHeader(http.StatusAccepted)
}

func indexGroup(g msgraph.ADGroup, logger *appinsights_wrapper.TelemetryClient) {
	hasGitHubAccess, err := msgraph.HasGitHubAccess(g.Id)
	if err != nil {
		logger.LogException(err)
		return
	}

	if hasGitHubAccess {
		db.ADGroup_Insert(g.Id, g.Name)
		logger.LogTrace(fmt.Sprintf("%s indexed.", g.Name), contracts.Information)
	}
}
