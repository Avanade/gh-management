package externallink

import (
	"encoding/json"
	"errors"
	"fmt"
	"main/model"
	"main/pkg/session"
	"main/service"
	"net/http"

	"github.com/gorilla/mux"
)

type externalLinkController struct {
	*service.Service
}

func NewExternalLinkController(serv *service.Service) ExternalLinkController {
	return &externalLinkController{
		Service: serv,
	}
}

func (c *externalLinkController) GetEnabledExternalLinks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	externalLinks, err := c.Service.ExternalLink.GetAllEnabled()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(externalLinks)
}

func (c *externalLinkController) GetExternalLinkById(w http.ResponseWriter, r *http.Request) {
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
	externalLinks, err := c.Service.ExternalLink.GetByID(string(params["id"]))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(externalLinks)
}

func (c *externalLinkController) GetExternalLinks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	externalLinks, err := c.Service.ExternalLink.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(externalLinks)
}

func (c *externalLinkController) CreateExternalLink(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var externalLink model.ExternalLink
	err := json.NewDecoder(r.Body).Decode(&externalLink)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.New("error unmarshalling data"))
		return
	}
	err = c.Service.ExternalLink.Validate(&externalLink)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.New(err.Error()))
		return
	}

	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := fmt.Sprint(profile["preferred_username"])
	externalLink.CreatedBy = username

	result, err := c.Service.ExternalLink.Create(&externalLink)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.New("error saving the external link"))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (c *externalLinkController) RemoveExternalLinkById(w http.ResponseWriter, r *http.Request) {
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
	err := c.Service.ExternalLink.Delete(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.New("error saving the external link"))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (c *externalLinkController) UpdateExternalLinkById(w http.ResponseWriter, r *http.Request) {
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
	err = c.Service.ExternalLink.Validate(&externalLink)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.New(err.Error()))
		return
	}

	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := fmt.Sprint(profile["preferred_username"])
	externalLink.ModifiedBy = username

	id := params["id"]
	newExternalLink, err := c.Service.ExternalLink.Update(id, &externalLink)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.New("error saving the external link"))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newExternalLink)
}
