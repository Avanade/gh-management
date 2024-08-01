package osscontributionsponsor

import (
	"errors"
	"main/model"
	repositoryOssContributionSponsor "main/repository/osscontributionsponsor"
	"strconv"
)

type ossContributionSponsorService struct {
	repositoryOssContributionSponsor repositoryOssContributionSponsor.OSSContributionSponsorRepository
}

// Create implements OssContributionSponsorService.
func (s *ossContributionSponsorService) Create(ossContributionSponsor *model.OSSContributionSponsor) (*model.OSSContributionSponsor, error) {
	return s.repositoryOssContributionSponsor.Insert(ossContributionSponsor)
}

// GetAll implements OssContributionSponsorService.
func (s *ossContributionSponsorService) GetAll() ([]model.OSSContributionSponsor, error) {
	return s.repositoryOssContributionSponsor.Select()
}

// GetByIsArchived implements OssContributionSponsorService.
func (s *ossContributionSponsorService) GetAllEnabled() ([]model.OSSContributionSponsor, error) {
	return s.repositoryOssContributionSponsor.SelectByIsArchived(false)
}

// Update implements OssContributionSponsorService.
func (s *ossContributionSponsorService) Update(id string, ossContributionSponsor *model.OSSContributionSponsor) (*model.OSSContributionSponsor, error) {
	parseId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}
	return s.repositoryOssContributionSponsor.Update(parseId, ossContributionSponsor)
}

// Validate implements OssContributionSponsorService.
func (s *ossContributionSponsorService) Validate(ossContributionSponsor *model.OSSContributionSponsor) error {
	if ossContributionSponsor.Name == "" {
		return errors.New("name is required")
	}
	return nil
}

func NewOssContributionSponsorService(repositoryOssContributionSponsor repositoryOssContributionSponsor.OSSContributionSponsorRepository) OssContributionSponsorService {
	return &ossContributionSponsorService{
		repositoryOssContributionSponsor: repositoryOssContributionSponsor,
	}
}
