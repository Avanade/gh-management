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
		r.NewActivity(&db),
		r.NewActivityContributionArea(&db),
		r.NewActivityHelp(&db),
		r.NewActivityType(&db),
		r.NewApprovalType(&db),
		r.NewApprover(&db),
		r.NewContributionArea(&db),
		r.NewExternalLink(&db),
		r.NewOssContributionSponsor(&db),
		r.NewUser(&db),
		r.NewPopularTopic(&db),
		r.NewCategory(&db),
		r.NewArticle(&db))

	serv = s.NewService(
		s.NewActivityService(repo),
		s.NewActivityTypeService(repo),
		s.NewActivityHelpService(repo),
		s.NewApprovalTypeService(repo),
		s.NewApproverService(repo),
		s.NewContributionAreaService(repo),
		s.NewEmailService(conf),
		s.NewExternalLinkService(repo),
		s.NewOssContributionSponsorService(repo),
		s.NewUserService(repo),
		s.NewTopicService(repo),
		s.NewCategoryService(repo),
		s.NewArticleService(repo))

	cont = c.NewController(
		c.NewActivityController(serv),
		c.NewActivityTypeController(serv),
		c.NewApprovalTypeController(serv),
		c.NewContributionAreaController(serv),
		c.NewExternalLinkController(serv),
		c.NewOssContributionSponsorController(serv),
		c.NewRepositoryApproverController(serv, conf),
		c.NewTopicController(serv),
		c.NewCategoryController(serv),
		c.NewArticleController(serv))

	httpRouter router.Router = router.NewMuxRouter()
)
