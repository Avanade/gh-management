package main

import (
	"main/config"
	"main/infrastructure/database"
	r "main/repository"
	"main/router"
	s "main/service"

	controllerActivity "main/controller/activity"
	controllerActivityType "main/controller/activitytype"
	controllerContributionArea "main/controller/contributionarea"
	controllerExternalLink "main/controller/externallink"
	controllerOssContributionSponsor "main/controller/osscontributionsponsor"
)

var (
	conf config.ConfigManager = config.NewEnvConfigManager()
	db   database.Database    = database.NewDatabase(conf)

	repo = r.NewRepository(
		r.NewActivity(db),
		r.NewActivityContributionArea(db),
		r.NewActivityHelp(db),
		r.NewActivityType(db),
		r.NewContributionArea(db),
		r.NewExternalLink(db),
		r.NewOssContributionSponsor(db))

	serv = s.NewService(
		s.NewActivityService(repo),
		s.NewActivityTypeService(repo),
		s.NewContributionAreaService(repo),
		s.NewEmailService(conf),
		s.NewExternalLinkService(repo),
		s.NewOssContributionSponsorService(repo))

	externalLinkController           controllerExternalLink.ExternalLinkController                     = controllerExternalLink.NewExternalLinkController(serv)
	contributionAreaController       controllerContributionArea.ContributionAreaController             = controllerContributionArea.NewContributionAreaController(serv)
	ossContributionSponsorController controllerOssContributionSponsor.OSSContributionSponsorController = controllerOssContributionSponsor.NewOssContributionSponsorController(serv)
	activityTypeController           controllerActivityType.ActivityTypeController                     = controllerActivityType.NewActivityTypeController(serv)
	activityController               controllerActivity.ActivityController                             = controllerActivity.NewActivityController(serv)

	httpRouter router.Router = router.NewMuxRouter()
)
