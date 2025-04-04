package adoOrganization

import (
	"main/model"
)

type AdoOrganizationRepository interface {
	Insert(*model.AdoOrganizationRequest) (int, error)
	SelectByUser(user string) ([]model.AdoOrganizationRequest, error)
}
