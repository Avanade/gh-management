package routes

import (
	"encoding/json"
	"fmt"
	models "main/models"
	db "main/pkg/ghmgmtdb"
	session "main/pkg/session"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ActivitiesDto struct {
	Data  interface{} `json: "data"`
	Total int         `json: "total"`
}

type ActivityDto struct {
	Name        string  `json: "name"`
	Url         string  `json: "url"`
	Date        string  `json: "date"`
	Type        ItemDto `json: "type"`
	CommunityId int     `json: "communityid"`
	CreatedBy   string
	ModifiedBy  string

	PrimaryContributionArea     ItemDto   `json: "primarycontributionarea"`
	AdditionalContributionAreas []ItemDto `json: "additionalcontributionareas"`
}

type ItemDto struct {
	Id   int    `json: "id"`
	Name string `json: "name"`
}

type CommunityActivitiesContributionAreasDto struct {
	CommunityActivityId int
	ContributionAreaId  int
	IsPrimary           bool
	CreatedBy           string
	ModifiedBy          string
}

func GetActivities(w http.ResponseWriter, r *http.Request) {
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
			Total: db.CommunitiesActivities_TotalCount_ByCreatedBy(username),
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
}

func CreateActivity(w http.ResponseWriter, r *http.Request) {
	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := fmt.Sprint(profile["preferred_username"])

	var body ActivityDto
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// CHECK ACTIVITY TYPE IF EXIST / INSERT IF NOT EXIST
	if body.Type.Id == 0 {
		id, err := db.ActivityTypes_Insert(body.Type.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		body.Type.Id = id
	}

	// COMMUNITY ACTIVITY
	communityActivityId, err := db.CommunitiesActivities_Insert(models.Activity{
		Name:        body.Name,
		Url:         body.Url,
		Date:        body.Date,
		TypeId:      body.Type.Id,
		CommunityId: body.CommunityId,
		CreatedBy:   username,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// PRIMARY CONTRIBUTION AREA
	err = InsertCommunityActivitiesContributionArea(body.PrimaryContributionArea, models.CommunityActivitiesContributionAreas{
		CommunityActivityId: communityActivityId,
		ContributionAreaId:  body.PrimaryContributionArea.Id,
		IsPrimary:           true,
		CreatedBy:           username,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// ADDITIONAL CONTRIBUTION AREA
	for _, contributionArea := range body.AdditionalContributionAreas {
		err = InsertCommunityActivitiesContributionArea(contributionArea, models.CommunityActivitiesContributionAreas{
			CommunityActivityId: communityActivityId,
			ContributionAreaId:  contributionArea.Id,
			IsPrimary:           false,
			CreatedBy:           username,
		})
		if err != nil {
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
		w.WriteHeader(http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func InsertCommunityActivitiesContributionArea(ca ItemDto, caca models.CommunityActivitiesContributionAreas) error {
	if ca.Id == 0 {
		id, err := db.ContributionAreas_Insert(ca.Name, caca.CreatedBy)
		if err != nil {
			return err
		}

		caca.ContributionAreaId = id
	}

	_, err := db.CommunityActivitiesContributionAreas_Insert(models.CommunityActivitiesContributionAreas{
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
