package model

import "time"

type ContributionArea struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Created    time.Time `json:"created"`
	CreatedBy  string    `json:"createdBy"`
	Modified   time.Time `json:"modified"`
	ModifiedBy string    `json:"modifiedBy"`
}
