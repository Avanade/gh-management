package contributionarea

import (
	"encoding/json"
	"errors"
	"fmt"
	"main/model"
	"main/pkg/session"
	service "main/service/contributionarea"
	"net/http"

	"github.com/gorilla/mux"
)

type contributionAreaController struct {
	contributionAreaService service.ContributionAreaService
}

// CreateContributionAreas implements ContributionAreaController.
func (c *contributionAreaController) CreateContributionAreas(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var contributionArea model.ContributionArea
	err := json.NewDecoder(r.Body).Decode(&contributionArea)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.New("error unmarshalling data"))
		return
	}
	err = c.contributionAreaService.Validate(&contributionArea)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.New(err.Error()))
		return
	}

	// temporary
	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := fmt.Sprint(profile["preferred_username"])
	contributionArea.CreatedBy = username

	result, err := c.contributionAreaService.Create(&contributionArea)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.New("error saving the external link"))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// GetContributionAreaById implements ContributionAreaController.
func (c *contributionAreaController) GetContributionAreaById(w http.ResponseWriter, r *http.Request) {
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
	contributionArea, err := c.contributionAreaService.GetByID(string(params["id"]))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contributionArea)
}

// GetContributionAreas implements ContributionAreaController.
func (c *contributionAreaController) GetContributionAreas(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	filter := ""
	offset := ""
	search := ""
	orderby := ""
	ordertype := ""
	if params["filter"] != nil {
		filter = params["filter"][0]
	}
	if params["offset"] != nil {
		offset = params["offset"][0]
	}
	if params["search"] != nil {
		search = params["search"][0]
	}
	if params["orderby"] != nil {
		orderby = params["orderby"][0]
	}
	if params["ordertype"] != nil {
		ordertype = params["ordertype"][0]
	}
	contributionAreas, total, err := c.contributionAreaService.Get(offset, filter, orderby, ordertype, search)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(
		GetResponseDto{
			Data:  contributionAreas,
			Total: total,
		},
	)
}

// UpdateContributionArea implements ContributionAreaController.
func (c *contributionAreaController) UpdateContributionArea(w http.ResponseWriter, r *http.Request) {
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
	var contributionArea model.ContributionArea
	err := json.NewDecoder(r.Body).Decode(&contributionArea)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.New("error unmarshalling data"))
		return
	}
	err = c.contributionAreaService.Validate(&contributionArea)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.New(err.Error()))
		return
	}

	// temporary
	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := fmt.Sprint(profile["preferred_username"])
	contributionArea.ModifiedBy = username

	id := params["id"]
	newContributionArea, err := c.contributionAreaService.Update(id, &contributionArea)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.New("error saving the external link"))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newContributionArea)
}

func NewContributionAreaController(contributionAreaService service.ContributionAreaService) ContributionAreaController {
	return &contributionAreaController{
		contributionAreaService: contributionAreaService,
	}
}
