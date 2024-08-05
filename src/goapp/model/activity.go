package model

import "time"

type Activity struct {
	ID             int64     `json:"id"`
	CommunityId    int64     `json:"communityId"`
	Date           time.Time `json:"date"`
	Name           string    `json:"name"`
	ActivityTypeId int64     `json:"activityTypeId"`
	Url            string    `json:"url"`
	Created        time.Time `json:"created"`
	CreatedBy      string    `json:"createdBy"`
	Modified       time.Time `json:"modified"`
	ModifiedBy     string    `json:"modifiedBy"`
}
