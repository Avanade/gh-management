package article

import (
	"main/model"
)

type CreateNewArticleRequest struct {
	Body     string         `json:"body"`
	Name     string         `json:"name"`
	Url      string         `json:"url"`
	Category model.Category `json:"category"`
}

type UpdateArticleRequest struct {
	Id       int64          `json:"id"`
	Body     string         `json:"body"`
	Name     string         `json:"name"`
	Url      string         `json:"url"`
	Category model.Category `json:"category"`
}
