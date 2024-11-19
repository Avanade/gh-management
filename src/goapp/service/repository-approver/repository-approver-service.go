package repositoryApprover

import (
	"main/model"
	"main/repository"
)

type repositoryApproverService struct {
	Repository *repository.Repository
}

func NewRepositoryApproverService(repository *repository.Repository) RepositoryApproverService {
	return &repositoryApproverService{repository}
}

func (s *repositoryApproverService) Create(approver *model.RepositoryApprover) error {
	return s.Repository.RepositoryApprover.Insert(approver)
}

func (s *repositoryApproverService) Get(approvalTypeId int) ([]model.RepositoryApprover, error) {
	return s.Repository.RepositoryApprover.SelectByApprovalTypeId(approvalTypeId)
}
