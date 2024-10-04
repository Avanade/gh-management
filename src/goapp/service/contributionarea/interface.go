package contributionarea

import "main/model"

type Total int64

type ContributionAreaService interface {
	Get(offset, filter, orderby, ordertype, search string) (contributionAreas []model.ContributionArea, total int64, err error)
	GetByID(id string) (*model.ContributionArea, error)
	Create(contributionArea *model.ContributionArea) (*model.ContributionArea, error)
	Update(id string, contributionArea *model.ContributionArea) (*model.ContributionArea, error)
	Validate(contributionArea *model.ContributionArea) error
}
