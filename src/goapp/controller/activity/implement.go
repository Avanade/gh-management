package activity

import (
	"errors"
	"fmt"
	"main/model"
	"main/pkg/session"
	"main/service"
	"net/http"

	"github.com/goccy/go-json"
	"github.com/gorilla/mux"
)

type activityController struct {
	*service.Service
}

// CreateActivity implements ActivityController.
func (c *activityController) CreateActivity(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var requestBody CreateActivityRequest
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.New("error unmarshalling data"))
		return
	}

	activity := model.Activity{
		Name:           requestBody.Name,
		Date:           requestBody.Date,
		Url:            requestBody.Url,
		CommunityId:    requestBody.CommunityId,
		ActivityTypeId: requestBody.Type.ID,
	}

	activity.ActivityType = model.ActivityType{
		ID:   requestBody.Type.ID,
		Name: requestBody.Type.Name,
	}

	for _, contributionArea := range requestBody.ContributionAreas {
		activity.ActivityContributionAreas = append(activity.ActivityContributionAreas, model.ActivityContributionArea{
			ActivityId:         activity.ID,
			ContributionAreaId: contributionArea.ID,
			IsPrimary:          contributionArea.IsPrimary,
			ContributionArea: model.ContributionArea{
				ID:   contributionArea.ID,
				Name: contributionArea.Name,
			},
		})
	}

	// temporary
	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := fmt.Sprint(profile["preferred_username"])
	activity.CreatedBy = username

	err = c.Service.Activity.Validate(&activity)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.New(err.Error()))
		return
	}

	result, err := c.Service.Activity.Create(&activity)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.New("error saving the activity"))
		return
	}

	if requestBody.Help != nil {
		help := model.ActivityHelp{
			ActivityId: result.ID,
			HelpTypeId: requestBody.Help.ID,
			Details:    requestBody.Help.Details,
		}
		_, err := c.Service.ActivityHelp.Create(&help)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errors.New("error saving the help"))
			return
		}

		// send email
		// emailSender, err := c.Service.Email.Connect()
		// if err != nil {
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	json.NewEncoder(w).Encode(err)
		// 	return
		// }
		// to := []string{"TO"}
		// cc := []string{"CC"}
		// subject := "Subject"
		// body := "Body"
		// contentType := c.Email.Connect().Html
		// isSaveToSetItem := false
		// emailSender.SendEmail(to, cc, subject, body, contentType, isSaveToSetItem)
		// if err != nil {
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	json.NewEncoder(w).Encode(err)
		// 	return
		// }
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

	// temporary
	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := fmt.Sprint(profile["preferred_username"])
	createdBy := username

	activities, total, err := c.Service.Activity.Get(offset, filter, orderby, ordertype, search, createdBy)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(GetActivitiesResponse{
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
	activity, err := c.Service.Activity.GetById(string(params["id"]))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(activity)
}

func NewActivityController(serv *service.Service) ActivityController {
	return &activityController{
		Service: serv,
	}
}

func buildActivityHelpEmail(activity *model.Activity, help *model.ActivityHelp) string {
	return fmt.Sprintf("Activity: %s\n\nHelp: %s\n\nDetails: %s\n\n", activity.Name, help.HelpTypeId, help.Details)
}
