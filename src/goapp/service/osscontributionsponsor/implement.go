package osscontributionsponsor

import (
	"errors"
	"main/model"
	"main/repository"
	"strconv"
)

type ossContributionSponsorService struct {
	Repository *repository.Repository
}

func NewOssContributionSponsorService(repo *repository.Repository) OssContributionSponsorService {
	return &ossContributionSponsorService{
		Repository: repo,
	}
}

func (s *ossContributionSponsorService) Create(ossContributionSponsor *model.OSSContributionSponsor) (*model.OSSContributionSponsor, error) {
	return s.Repository.OssContributionSponsor.Insert(ossContributionSponsor)
}

func (s *ossContributionSponsorService) GetAll() ([]model.OSSContributionSponsor, error) {
	return s.Repository.OssContributionSponsor.Select()
}

func (s *ossContributionSponsorService) GetAllEnabled() ([]model.OSSContributionSponsor, error) {
	return s.Repository.OssContributionSponsor.SelectByIsArchived(false)
}

func (s *ossContributionSponsorService) Update(id string, ossContributionSponsor *model.OSSContributionSponsor) (*model.OSSContributionSponsor, error) {
	parseId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}
	return s.Repository.OssContributionSponsor.Update(parseId, ossContributionSponsor)
}

func (s *ossContributionSponsorService) Validate(ossContributionSponsor *model.OSSContributionSponsor) error {
	if ossContributionSponsor.Name == "" {
		return errors.New("name is required")
	}
	return nil
}
