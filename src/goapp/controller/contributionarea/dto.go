package contributionarea

import "main/model"

type GetResponseDto struct {
	Data  []model.ContributionArea `json:"data"`
	Total int64                    `json:"total"`
}
