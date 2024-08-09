package activitytype

import (
	"encoding/json"
	"main/service"
	"net/http"
)

type activityTypeController struct {
	*service.Service
}

func NewActivityTypeController(serv *service.Service) ActivityTypeController {
	return &activityTypeController{
		Service: serv,
	}
}

func (c *activityTypeController) GetActivityTypes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	activityTypes, err := c.Service.ActivityType.Get()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(activityTypes)
}
