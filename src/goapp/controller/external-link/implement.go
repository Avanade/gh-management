package controller

import (
	"encoding/json"
	"errors"
	"main/model"
	service "main/service/external-link"
	"net/http"

	"github.com/gorilla/mux"
)

type externalLinkController struct {
	externalLinkService service.ExternalLinkService
}

// GetEnabledExternalLinks implements ExternalLinkController.
func (e *externalLinkController) GetEnabledExternalLinks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	externalLinks, err := e.externalLinkService.GetAllEnabled()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(externalLinks)
}

// GetExternalLinkById implements ExternalLinkController.
func (e *externalLinkController) GetExternalLinkById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if len(params) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.New("no parameters found"))
		return
	}
	if params["id"] == "" {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.New("no parameters found"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	externalLinks, err := e.externalLinkService.GetByID(string(params["id"]))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(externalLinks)
}

// GetExternalLinks implements ExternalLinkController.
func (e *externalLinkController) GetExternalLinks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	externalLinks, err := e.externalLinkService.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(externalLinks)
}

// CreateExternalLink implements ExternalLinkController.
func (e *externalLinkController) CreateExternalLink(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var externalLink model.ExternalLink
	err := json.NewDecoder(r.Body).Decode(&externalLink)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.New("error unmarshalling data"))
		return
	}
	err = e.externalLinkService.Validate(&externalLink)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.New(err.Error()))
		return
	}

	result, err := e.externalLinkService.Create(&externalLink)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.New("error saving the external link"))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// RemoveExternalLinkById implements ExternalLinkController.
func (e *externalLinkController) RemoveExternalLinkById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if len(params) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.New("no parameters found"))
		return
	}
	if params["id"] == "" {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.New("no parameters found"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	id := params["id"]
	err := e.externalLinkService.Delete(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.New("error saving the external link"))
		return
	}
	w.WriteHeader(http.StatusOK)
}

// UpdateExternalLinkById implements ExternalLinkController.
func (e *externalLinkController) UpdateExternalLinkById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if len(params) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.New("no parameters found"))
		return
	}
	if params["id"] == "" {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.New("no parameters found"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	var externalLink model.ExternalLink
	err := json.NewDecoder(r.Body).Decode(&externalLink)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.New("error unmarshalling data"))
		return
	}
	err = e.externalLinkService.Validate(&externalLink)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.New(err.Error()))
		return
	}

	id := params["id"]
	newExternalLink, err := e.externalLinkService.Update(id, &externalLink)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.New("error saving the external link"))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newExternalLink)
}

func NewExternalLinkController(externalLinkService service.ExternalLinkService) ExternalLinkController {
	return &externalLinkController{
		externalLinkService: externalLinkService,
	}
}
