package model

type ActivityContributionArea struct {
	ID                 int64 `json:"id"`
	ActivityId         int64 `json:"activityId"`
	ContributionAreaId int64 `json:"contributionAreaId"`
	IsPrimary          bool  `json:"isPrimary"`

	ContributionArea ContributionArea `json:"contributionArea"`
}
