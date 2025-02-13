package model

import "time"

type Article struct {
	Id         int64     `json:"id"`
	Name       string    `json:"name"`
	Url        string    `json:"url"`
	Body       string    `json:"body"`
	CategoryId int64     `json:"categoryId"`
	Created    time.Time `json:"created"`
	CreatedBy  string    `json:"created_by"`
	Modified   time.Time `json:"modified"`
	ModifiedBy string    `json:"modified_by"`

	Category Category `json:"category"`
}
