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

func (s *approverService) GetApproversByApprovalTypeId(approvalTypeId int) ([]model.Approver, error) {
	return s.Repository.Approver.GetApproversByApprovalTypeId(approvalTypeId)
}
