package models

type CommunityActivitiesContributionAreas struct {
	CommunityActivityId int
	ContributionAreaId  int
	IsPrimary           bool
	CreatedBy           string
	ModifiedBy          string
}
