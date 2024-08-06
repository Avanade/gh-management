package activitytype

import (
	"encoding/json"
	serviceActivityType "main/service/activitytype"
	"net/http"
)

type activityTypeController struct {
	activityTypeService serviceActivityType.ActivityTypeService
}

func NewActivityTypeController(activityTypeService serviceActivityType.ActivityTypeService) ActivityTypeController {
	return &activityTypeController{activityTypeService}
}

func (c *activityTypeController) GetActivityTypes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	activityTypes, err := c.activityTypeService.Get()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(activityTypes)
}
