package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	models "main/models"
	db "main/pkg/ghmgmtdb"
	session "main/pkg/session"
)

type ActivityDto struct {
	Name        string `json: "name"`
	Url         string `json: "url"`
	Date        string `json: "date"`
	TypeId      int    `json: "typeid"`
	CommunityId int    `json: "communityid"`
	CreatedBy   string
	ModifiedBy  string

	PrimaryContributionArea     []int `json: "primarycontributionarea"`
	AdditionalContributionAreas []int `json: "additionalcontributionareas"`
}

type CommunityActivitiesContributionAreasDto struct {
	CommunityActivityId int
	ContributionAreaId  int
	IsPrimary           bool
	CreatedBy           string
	ModifiedBy          string
}

func CreateActivity(w http.ResponseWriter, r *http.Request) {
	// COMMUNITY :  SECOND
	// DATE :  2022-06-15
	// ACTIVITY :  FIRST ACTIVITY
	// TYPE First Activity
	// URL :  sample.com
	// PRIMARY CONTRIBUTION AREA :  FirstContributionArea
	// ADDITIONAL CONTRIBUTION AREAS
	// 	FirstContributionArea
	// 	SecondContributionArea

	// MOVE DATA FROM FORM TO MODEL
	sessionaz, _ := session.Store.Get(r, "auth-session")
	iprofile := sessionaz.Values["profile"]
	profile := iprofile.(map[string]interface{})
	username := fmt.Sprint(profile["preferred_username"])

	// communityId, _ := strconv.Atoi(r.FormValue("community_id"))
	// typeId, _ := strconv.Atoi(r.FormValue("type_id"))

	var body ActivityDto
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	body.CreatedBy = fmt.Sprint(profile["preferred_username"])

	// COMMUNITY ACTIVITY
	communityActivityId, err := db.CommunitiesActivities_Insert(models.Activity{
		Name:        body.Name,
		Url:         body.Url,
		Date:        body.Date,
		TypeId:      body.TypeId,
		CommunityId: body.CommunityId,
		CreatedBy:   username,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// PRIMARY CONTRIBUTION AREA
	for _, contributionAreaId := range body.PrimaryContributionArea {
		db.CommunityActivitiesContributionAreas_Insert(models.CommunityActivitiesContributionAreas{
			CommunityActivityId: communityActivityId,
			ContributionAreaId:  contributionAreaId,
			IsPrimary:           true,
			CreatedBy:           username,
		})
	}

	// ADDITIONAL CONTRIBUTION AREA
	for _, contributionAreaId := range body.AdditionalContributionAreas {
		db.CommunityActivitiesContributionAreas_Insert(models.CommunityActivitiesContributionAreas{
			CommunityActivityId: communityActivityId,
			ContributionAreaId:  contributionAreaId,
			IsPrimary:           false,
			CreatedBy:           username,
		})
	}

	fmt.Fprint(w, body)
}
