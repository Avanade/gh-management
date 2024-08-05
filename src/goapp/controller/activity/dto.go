package activity

import "main/model"

type GetResponseDto struct {
	Data  []model.Activity `json:"data"`
	Total int64            `json:"total"`
}
