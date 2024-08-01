package contributionarea

import (
	"main/model"
)

type ContributionAreaRepository interface {
	Select() ([]model.ContributionArea, error)
	SelectById(id int64) (*model.ContributionArea, error)
	SelectByOption(offset, filter int, orderby, ordertype, search string) ([]model.ContributionArea, error)
	Total() (int, error)
	Insert(contributionArea *model.ContributionArea) (*model.ContributionArea, error)
	Update(id int64, contributionArea *model.ContributionArea) (*model.ContributionArea, error)
}
