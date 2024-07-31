package contributionarea

import (
	"encoding/json"
	"errors"
	"main/model"
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
	filter := params["filter"][0]
	offset := params["offset"][0]
	search := params["search"][0]
	orderby := params["orderby"][0]
	ordertype := params["ordertype"][0]
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