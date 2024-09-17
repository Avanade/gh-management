package activity

import (
	"main/model"
	"time"
)

type CreateActivityRequest struct {
	Name        string    `json:"name"`
	Date        time.Time `json:"date"`
	Url         string    `json:"url"`
	CommunityId int64     `json:"communityId"`
	Type        struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"type"`
	ContributionAreas []struct {
		ID        int64  `json:"id"`
		Name      string `json:"name"`
		IsPrimary bool   `json:"isPrimary"`
	} `json:"contributionAreas"`
	Help *struct {
		ID      int64  `json:"id"`
		Name    string `json:"name"`
		Details string `json:"details"`
	} `json:"help"`
}

type GetActivitiesResponse struct {
	Data  []model.Activity `json:"data"`
	Total int64            `json:"total"`
}
