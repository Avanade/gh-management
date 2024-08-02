package model

import "time"

// External Link Model
type ExternalLink struct {
	ID          int64     `json:"id"`
	DisplayName string    `json:"displayName"`
	IconSVGPath string    `json:"iconSVGPath"`
	Hyperlink   string    `json:"hyperlink"`
	IsEnabled   bool      `json:"isEnabled"`
	Created     time.Time `json:"created"`
	CreatedBy   string    `json:"createdBy"`
	Modified    time.Time `json:"modified"`
	ModifiedBy  string    `json:"modifiedBy"`
}
