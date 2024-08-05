package activity

import (
	"errors"
	"main/model"
	serviceActivity "main/service/activity"
	"net/http"

	"github.com/goccy/go-json"
	"github.com/gorilla/mux"
)

type activityController struct {
	activityService serviceActivity.ActivityService
}

// CreateActivity implements ActivityController.
func (c *activityController) CreateActivity(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var activity model.Activity
	err := json.NewDecoder(r.Body).Decode(&activity)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.New("error unmarshalling data"))
		return
	}
	err = c.activityService.Validate(&activity)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.New(err.Error()))
		return
	}

	result, err := c.activityService.Create(&activity)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.New("error saving the activity"))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// GetActivities implements ActivityController.
func (c *activityController) GetActivities(w http.ResponseWriter, r *http.Request) {
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
	w.Header().Set("Content-Type", "application/json")
	activities, total, err := c.activityService.Get(offset, filter, orderby, ordertype, search)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(GetResponseDto{
		Data:  activities,
		Total: total,
	})
}

// GetActivityById implements ActivityController.
func (c *activityController) GetActivityById(w http.ResponseWriter, r *http.Request) {
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
	activity, err := c.activityService.GetById(string(params["id"]))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(activity)
}

func NewActivityController(activityService serviceActivity.ActivityService) ActivityController {
	return &activityController{
		activityService: activityService,
	}
}
