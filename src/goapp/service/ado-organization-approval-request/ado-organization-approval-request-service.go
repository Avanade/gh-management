package adoOrganizationApprovalRequest

import (
	"main/repository"
)

type adoOrganizationApprovalRequestService struct {
	Repository *repository.Repository
}

func NewAdoOrganizationApprovalRequestService(repo *repository.Repository) AdoOrganizationApprovalRequestService {
	return &adoOrganizationApprovalRequestService{repo}
}

func (s *adoOrganizationApprovalRequestService) Insert(adoOrganizationId int, approvalRequestId int64) error {
	return s.Repository.AdoOrganizationApprovalRequest.Insert(adoOrganizationId, approvalRequestId)
}
