package communityApprover

import (
	"main/model"
)

type CommunityApproverService interface {
	GetByCategory(category string) ([]model.CommunityApprover, error)
}
