package osscontributionsponsor

import (
	"main/model"
)

type OSSContributionSponsorRepository interface {
	GetAll() ([]model.OSSContributionSponsor, error)
	GetByIsArchived(isArchived bool) ([]model.OSSContributionSponsor, error)
	Create(ossContributionSponsor *model.OSSContributionSponsor) (*model.OSSContributionSponsor, error)
	Update(id int64, ossContributionSponsor *model.OSSContributionSponsor) (*model.OSSContributionSponsor, error)
}
