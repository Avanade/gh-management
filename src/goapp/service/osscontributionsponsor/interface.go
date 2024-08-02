package osscontributionsponsor

import "main/model"

type OssContributionSponsorService interface {
	Create(ossContributionSponsor *model.OSSContributionSponsor) (*model.OSSContributionSponsor, error)
	Update(id string, ossContributionSponsor *model.OSSContributionSponsor) (*model.OSSContributionSponsor, error)
	GetAll() ([]model.OSSContributionSponsor, error)
	GetAllEnabled() ([]model.OSSContributionSponsor, error)
	Validate(ossContributionSponsor *model.OSSContributionSponsor) error
}
