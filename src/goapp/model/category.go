package model

import "time"

type Category struct {
	Id         int64     `json:"id"`
	Name       string    `json:"name"`
	Created    time.Time `json:"created"`
	CreatedBy  string    `json:"created_by"`
	Modified   time.Time `json:"modified"`
	ModifiedBy string    `json:"modified_by"`

	Articles []Article `json:"article"`
}
