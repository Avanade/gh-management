package main

import (
	"main/config"
	"main/repository"
	"main/router"

	repositoryActivity "main/repository/activity"
	repositroyActivityType "main/repository/activitytype"
	repositoruContributionArea "main/repository/contributionarea"
	repositoryExternalLink "main/repository/externallink"
	repositoryOssContributionSponsor "main/repository/osscontributionsponsor"

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
	database      repository.Database  = repository.NewDatabase(configManager)

	externalLinkRepository repositoryExternalLink.ExternalLinkRepository = repositoryExternalLink.NewExternalLinkRepository(database)
	externalLinkService    serviceExternalLink.ExternalLinkService       = serviceExternalLink.NewExternalLinkService(externalLinkRepository)
	externalLinkController controllerExternalLink.ExternalLinkController = controllerExternalLink.NewExternalLinkController(externalLinkService)

	contributionAreaRepository repositoruContributionArea.ContributionAreaRepository = repositoruContributionArea.NewContributionAreaRepository(database)
	contributionAreaService    serviceContributionArea.ContributionAreaService       = serviceContributionArea.NewContributionAreaService(contributionAreaRepository)
	contributionAreaController controllerContributionArea.ContributionAreaController = controllerContributionArea.NewContributionAreaController(contributionAreaService)

	ossContributionSponsorRepository repositoryOssContributionSponsor.OSSContributionSponsorRepository = repositoryOssContributionSponsor.NewOSSContributionSponsorRepository(database)
	ossContributionSponsorService    serviceOssContributionSponsor.OssContributionSponsorService       = serviceOssContributionSponsor.NewOssContributionSponsorService(ossContributionSponsorRepository)
	ossContributionSponsorController controllerOssContributionSponsor.OSSContributionSponsorController = controllerOssContributionSponsor.NewOssContributionSponsorController(ossContributionSponsorService)

	activityTypeRepository repositroyActivityType.ActivityTypeRepository = repositroyActivityType.NewActivityTypeRepository(database)
	activityTypeService    serviceActivityType.ActivityTypeService       = serviceActivityType.NewActivityTypeService(activityTypeRepository)
	activityTypeController controllerActivityType.ActivityTypeController = controllerActivityType.NewActivityTypeController(activityTypeService)

	activityRepository repositoryActivity.ActivityRepository = repositoryActivity.NewActivityRepository(database)
	activityService    serviceActivity.ActivityService       = serviceActivity.NewActivityService(activityRepository)
	activityController controllerActivity.ActivityController = controllerActivity.NewActivityController(activityService)

	httpRouter router.Router = router.NewMuxRouter()
)
