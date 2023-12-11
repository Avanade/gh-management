package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"main/pkg/appinsights_wrapper"
	"main/pkg/email"
	"main/pkg/envvar"
	db "main/pkg/ghmgmtdb"
	"main/pkg/notification"
	"main/pkg/session"

	"github.com/gorilla/mux"
)

type ActivitiesDto struct {
	Data  interface{} `json:"data"`
	Total int         `json:"total"`
}

type ActivityDto struct {
	Name        string  `json:"name"`
	Url         string  `json:"url"`
	Date        string  `json:"date"`
	Type        ItemDto `json:"type"`
	CommunityId int     `json:"communityid"`
	CreatedBy   string
	ModifiedBy  string

	PrimaryContributionArea     ItemDto   `json:"primarycontributionarea"`
	AdditionalContributionAreas []ItemDto `json:"additionalcontributionareas"`

	Help HelpDto `json:"help"`
}

type HelpDto struct {
	ItemDto
	Details string `json:"details"`
}

type ItemDto struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func GetActivities(w http.ResponseWriter, r *http.Request) {
	client := appinsights_wrapper.NewClient()
	client.StartOperation("GET ACTIVITIES")

	client.TrackEvent("START GET ACTIVITIES")

	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := fmt.Sprint(profile["preferred_username"])

	var result ActivitiesDto

	params := r.URL.Query()

	if params.Has("offset") && params.Has("filter") {
		filter, _ := strconv.Atoi(params["filter"][0])
		offset, _ := strconv.Atoi(params["offset"][0])
		search := params["search"][0]
		orderby := params["orderby"][0]
		ordertype := params["ordertype"][0]
		result = ActivitiesDto{
			Data:  db.CommunitiesActivities_Select_ByOffsetAndFilterAndCreatedBy(offset, filter, orderby, ordertype, search, username),
			Total: db.CommunitiesActivities_TotalCount_ByCreatedBy(username, search),
		}
	} else {
		result = ActivitiesDto{
			Data:  db.CommunitiesActivities_Select(),
			Total: db.CommunitiesActivities_TotalCount(),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)

	client.TrackEvent("END GET ACTIVITIES")
	client.EndOperation()
}

func CreateActivity(w http.ResponseWriter, r *http.Request) {
	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := fmt.Sprint(profile["preferred_username"])

	var body ActivityDto
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// CHECK ACTIVITY TYPE IF EXIST / INSERT IF NOT EXIST
	if body.Type.Id == 0 {
		id, err := db.ActivityTypes_Insert(body.Type.Name)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		body.Type.Id = id
	}

	// COMMUNITY ACTIVITY
	communityActivityId, err := db.CommunitiesActivities_Insert(db.Activity{
		Name:        body.Name,
		Url:         body.Url,
		Date:        body.Date,
		TypeId:      body.Type.Id,
		CommunityId: body.CommunityId,
		CreatedBy:   username,
	})
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if body.Help.Id != 0 {
		if os.Getenv("NOTIFICATION_EMAIL_SUPPORT") == "" {
			log.Println("Activity does not send successfully because the notification support email environment is not set.")
			return
		}

		recipients := strings.Split(os.Getenv("NOTIFICATION_EMAIL_SUPPORT"), ",")

		scheme := envvar.GetEnvVar("SCHEME", "https")

		activityLink := fmt.Sprint(scheme, "://", r.Host, "/activities/view/", communityActivityId)

		messageBody := notification.ActivityAddedRequestForHelpMessageBody{
			Recipients:   recipients,
			ActivityLink: activityLink,
			UserName:     profile["name"].(string),
		}
		err = messageBody.Send()
		if err != nil {
			log.Println(err.Error())
		}

		errHelp := processHelp(communityActivityId, activityLink, username, profile["name"].(string), body.Help)
		if errHelp != nil {
			log.Println(err.Error())
			http.Error(w, errHelp.Error(), http.StatusBadRequest)
			return
		}
	}

	// PRIMARY CONTRIBUTION AREA
	err = insertCommunityActivitiesContributionArea(body.PrimaryContributionArea, db.CommunityActivitiesContributionAreas{
		CommunityActivityId: communityActivityId,
		ContributionAreaId:  body.PrimaryContributionArea.Id,
		IsPrimary:           true,
		CreatedBy:           username,
	})
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// ADDITIONAL CONTRIBUTION AREA
	for _, contributionArea := range body.AdditionalContributionAreas {
		err = insertCommunityActivitiesContributionArea(contributionArea, db.CommunityActivitiesContributionAreas{
			CommunityActivityId: communityActivityId,
			ContributionAreaId:  contributionArea.Id,
			IsPrimary:           false,
			CreatedBy:           username,
		})
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	fmt.Fprint(w, body)
}

func GetActivityById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	result, err := db.CommunitiesActivities_Select_ById(id)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func insertCommunityActivitiesContributionArea(ca ItemDto, caca db.CommunityActivitiesContributionAreas) error {
	if ca.Id == 0 {
		id, err := db.ContributionAreas_Insert(ca.Name, caca.CreatedBy)
		if err != nil {
			return err
		}

		caca.ContributionAreaId = id
	}

	_, err := db.CommunityActivitiesContributionAreas_Insert(db.CommunityActivitiesContributionAreas{
		CommunityActivityId: caca.CommunityActivityId,
		ContributionAreaId:  caca.ContributionAreaId,
		IsPrimary:           caca.IsPrimary,
		CreatedBy:           caca.CreatedBy,
	})
	if err != nil {
		return err
	}

	return nil
}

func processHelp(activityId int, activityLink, requestorEmail string, requestorName string, h HelpDto) error {
	// INSERT
	_, err := db.CommunityActivitiesHelpTypes_Insert(activityId, h.Id, h.Details)
	if err != nil {
		return err
	}

	body := fmt.Sprintf("<p>%s added an <a href=\"%s\">activity</a> and is requesting for help below are the details of the request.</p><p>DETAILS: %s</p>", requestorName, activityLink, h.Details)

	// SEND EMAIL
	m := email.Message{
		Subject: fmt.Sprintf("%s : Community Portal", h.Name),
		Body: email.Body{
			Content: body,
			Type:    email.HtmlMessageType,
		},
		ToRecipients: []email.Recipient{
			{
				Email: os.Getenv("EMAIL_SUPPORT"),
			},
		},
		CcRecipients: []email.Recipient{
			{
				Email: requestorEmail,
			},
		},
	}

	errEmail := email.SendEmail(m, false)
	if errEmail != nil {
		return errEmail
	}

	// NO ERROR
	return nil
}
