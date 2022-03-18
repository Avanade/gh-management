package main

import (
	"net/http"

	"webserver/models"
)

func handleGraphRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/text")
	sm := models.GetRepositoryStateMachine("InnerSourceInternal", true, false)

	graphText := sm.ToGraph()
	w.Write([]byte(graphText))
	return
}
