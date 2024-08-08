package main

import (
	"main/config"
	"main/infrastructure/database"
	r "main/repository"
	"main/router"

	serviceActivity "main/service/activity"
	serviceActivityType "main/service/activitytype"
	serviceContributionArea "main/service/contributionarea"
	serviceExternalLink "main/service/externallink"
	serviceOssContributionSponsor "main/service/osscontributionsponsor"

	controllerActivity "main/controller/activity"
	controllerActivityType "main/controller/activitytype"
	controllerContributionArea "main/controller/contributionarea"
	controllerExternalLink "main/controller/externallink"
	controllerOssContributionSponsor "main/controller/osscontributionsponsor"
)

var (
	configManager config.ConfigManager = config.NewEnvConfigManager()
	db            database.Database    = database.NewDatabase(configManager)

	repo = r.NewRepository(
		r.NewActivity(db),
		r.NewActivityType(db),
		r.NewContributionArea(db),
		r.NewExternalLink(db),
		r.NewOssContributionSponsor(db))

	externalLinkService    serviceExternalLink.ExternalLinkService       = serviceExternalLink.NewExternalLinkService(repo)
	externalLinkController controllerExternalLink.ExternalLinkController = controllerExternalLink.NewExternalLinkController(externalLinkService)

	contributionAreaService    serviceContributionArea.ContributionAreaService       = serviceContributionArea.NewContributionAreaService(repo)
	contributionAreaController controllerContributionArea.ContributionAreaController = controllerContributionArea.NewContributionAreaController(contributionAreaService)

	ossContributionSponsorService    serviceOssContributionSponsor.OssContributionSponsorService       = serviceOssContributionSponsor.NewOssContributionSponsorService(repo)
	ossContributionSponsorController controllerOssContributionSponsor.OSSContributionSponsorController = controllerOssContributionSponsor.NewOssContributionSponsorController(ossContributionSponsorService)

	activityTypeService    serviceActivityType.ActivityTypeService       = serviceActivityType.NewActivityTypeService(repo)
	activityTypeController controllerActivityType.ActivityTypeController = controllerActivityType.NewActivityTypeController(activityTypeService)

	activityService    serviceActivity.ActivityService       = serviceActivity.NewActivityService(repo)
	activityController controllerActivity.ActivityController = controllerActivity.NewActivityController(activityService)

	httpRouter router.Router = router.NewMuxRouter()
)
