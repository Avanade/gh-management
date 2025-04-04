package adoOrganization

import (
	"main/model"
)

type AdoOrganizationService interface {
	GetByUser(user string) ([]model.AdoOrganizationRequest, error)
	Insert(adoOrgRequest *model.AdoOrganizationRequest) (int, error)
}
