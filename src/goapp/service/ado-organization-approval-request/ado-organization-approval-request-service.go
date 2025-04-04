package adoOrganizationApprovalRequest

import (
	"main/model"
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

func (s *adoOrganizationApprovalRequestService) SelectByAdoOrganizationId(adoOrganizationId int64) ([]model.ApprovalRequest, error) {
	return s.Repository.AdoOrganizationApprovalRequest.SelectByAdoOrganizationId(adoOrganizationId)
}
