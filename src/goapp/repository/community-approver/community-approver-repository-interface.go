package communityApprover

import (
	"main/model"
)

type CommunityApproverRepository interface {
	GetByCategory(category string) ([]model.CommunityApprover, error)
}
