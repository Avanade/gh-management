package model

type OSSContributionSponsor struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	IsArchived bool   `json:"isArchived"`
}
