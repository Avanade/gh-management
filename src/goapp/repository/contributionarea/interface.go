package contributionarea

import (
	"main/model"
)

type ContributionAreaRepository interface {
	GetAll() ([]model.ContributionArea, error)
	GetByID(id int64) (*model.ContributionArea, error)
	GetByOption(offset, filter int, orderby, ordertype, search string) ([]model.ContributionArea, error)
	GetTotal() (int, error)
	Create(contributionArea *model.ContributionArea) (*model.ContributionArea, error)
	Update(id int64, contributionArea *model.ContributionArea) (*model.ContributionArea, error)
}
