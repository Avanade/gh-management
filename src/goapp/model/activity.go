package model

import "time"

type Activity struct {
	ID             int64     `json:"id"`
	Date           time.Time `json:"date"`
	Name           string    `json:"name"`
	Url            string    `json:"url"`
	Created        time.Time `json:"created"`
	CreatedBy      string    `json:"createdBy"`
	Modified       time.Time `json:"modified"`
	ModifiedBy     string    `json:"modifiedBy"`
	CommunityId    int64     `json:"communityId"`
	ActivityTypeId int64     `json:"activityTypeId"`

	Community                 Community                  `json:"community"`
	ActivityType              ActivityType               `json:"type"`
	ActivityContributionAreas []ActivityContributionArea `json:"activityContributionAreas"`
}
