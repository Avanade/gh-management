package communityApprover

import (
	"main/model"
	"main/repository"
)

type communityApproverService struct {
	repository *repository.Repository
}

func NewCommunityApproverService(repository *repository.Repository) CommunityApproverService {
	return &communityApproverService{repository}
}

func (s *communityApproverService) GetByCategory(category string) ([]model.CommunityApprover, error) {
	return s.repository.CommunityApprover.GetByCategory(category)
}
