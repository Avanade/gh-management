package osscontributionsponsor

import (
	"encoding/json"
	"errors"
	"main/model"
	"main/service"
	"net/http"

	"github.com/gorilla/mux"
)

type ossContributionSponsorController struct {
	*service.Service
}

func NewOssContributionSponsorController(serv *service.Service) OSSContributionSponsorController {
	return &ossContributionSponsorController{
		Service: serv,
	}
}

func (c *ossContributionSponsorController) CreateOssContributionSponsor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var ossContributionSponsor model.OSSContributionSponsor
	err := json.NewDecoder(r.Body).Decode(&ossContributionSponsor)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.New("error unmarshalling data"))
		return
	}
	err = c.Service.OssContributionSponsor.Validate(&ossContributionSponsor)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	result, err := c.Service.OssContributionSponsor.Create(&ossContributionSponsor)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.New("error saving the oss contribution sponsor"))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (c *ossContributionSponsorController) GetOssContributionSponsors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ossContributionSponsor, err := c.Service.OssContributionSponsor.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ossContributionSponsor)
}

func (c *ossContributionSponsorController) GetEnabledOssContributionSponsors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ossContributionSponsor, err := c.Service.OssContributionSponsor.GetAllEnabled()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ossContributionSponsor)
}

func (c *ossContributionSponsorController) UpdateOssContributionSponsor(w http.ResponseWriter, r *http.Request) {
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
	var ossContributionSponsor model.OSSContributionSponsor
	err := json.NewDecoder(r.Body).Decode(&ossContributionSponsor)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.New("error unmarshalling data"))
		return
	}
	err = c.Service.OssContributionSponsor.Validate(&ossContributionSponsor)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.New(err.Error()))
		return
	}

	id := params["id"]
	newOssContributionSponsor, err := c.Service.OssContributionSponsor.Update(id, &ossContributionSponsor)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.New("error saving the oss contribution sponsor"))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newOssContributionSponsor)
}
