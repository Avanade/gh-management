package osscontributionsponsor

import (
	"main/model"
)

type OSSContributionSponsorRepository interface {
	Select() ([]model.OSSContributionSponsor, error)
	SelectByIsArchived(isArchived bool) ([]model.OSSContributionSponsor, error)
	Insert(ossContributionSponsor *model.OSSContributionSponsor) (*model.OSSContributionSponsor, error)
	Update(id int64, ossContributionSponsor *model.OSSContributionSponsor) (*model.OSSContributionSponsor, error)
}
