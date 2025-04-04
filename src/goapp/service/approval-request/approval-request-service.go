package approvalRequest

import (
	"main/repository"
)

type approvalRequestService struct {
	Repository *repository.Repository
}

func NewApprovalRequestService(repository *repository.Repository) ApprovalRequestService {
	return &approvalRequestService{repository}
}

func (s *approvalRequestService) Insert(approver, description, requestor string) (id int64, err error) {
	return s.Repository.ApprovalRequest.Insert(approver, description, requestor)
}

func (s *approvalRequestService) UpdateApprovalSystemGUID(requestId int64, approvalSystemGUID string) error {
	return s.Repository.ApprovalRequest.UpdateApprovalSystemGUID(requestId, approvalSystemGUID)
}
