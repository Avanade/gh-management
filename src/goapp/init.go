package main

import (
	"main/config"
	c "main/controller"
	"main/infrastructure/database"
	r "main/repository"
	"main/router"
	s "main/service"
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
		s.NewActivityHelpService(repo),
		s.NewEmailService(conf),
		s.NewExternalLinkService(repo),
		s.NewOssContributionSponsorService(repo))

	cont = c.NewController(
		c.NewActivityController(serv),
		c.NewActivityTypeController(serv),
		c.NewContributionAreaController(serv),
		c.NewExternalLinkController(serv),
		c.NewOssContributionSponsorController(serv))

	httpRouter router.Router = router.NewMuxRouter()
)
