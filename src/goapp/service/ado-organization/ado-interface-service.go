package adoOrganization

import (
	"main/model"
	"main/repository"
)

type adoOrganizationService struct {
	Repository *repository.Repository
}

func NewAdoOrganizationService(repository *repository.Repository) AdoOrganizationService {
	return &adoOrganizationService{repository}
}

func (s *adoOrganizationService) GetByUser(user string) ([]model.AdoOrganizationRequest, error) {
	data, err := s.Repository.AdoOrganization.SelectByUser(user)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *adoOrganizationService) Insert(adoOrgRequest *model.AdoOrganizationRequest) (int, error) {
	id, err := s.Repository.AdoOrganization.Insert(adoOrgRequest)
	if err != nil {
		return 0, err
	}
	return id, nil
}
