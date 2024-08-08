package activitycontributionarea

import "main/model"

type ActivityContributionAreaRepository interface {
	Insert(activityContributionArea *model.ActivityContributionArea) (*model.ActivityContributionArea, error)
	SelectByActivityId(activityId int64) ([]model.ActivityContributionArea, error)
}
