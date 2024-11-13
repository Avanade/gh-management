package approver

import (
	"main/model"
	"main/repository"
)

type approverService struct {
	Repository *repository.Repository
}

func NewApproverService(repository *repository.Repository) ApproverService {
	return &approverService{repository}
}

func (s *approverService) Get(approvalTypeId int) ([]model.RepositoryApprover, error) {
	return s.Repository.Approver.SelectByApprovalTypeId(approvalTypeId)
}
