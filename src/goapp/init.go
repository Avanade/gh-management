package main

import (
	"main/config"
	"main/repository"
	"main/router"

	controllerExternalLink "main/controller/externallink"

	repositoryExternalLink "main/repository/externallink"

	serviceExternalLink "main/service/externallink"
)

var (
	configManager          config.ConfigManager                          = config.NewEnvConfigManager()
	database               repository.Database                           = repository.NewDatabase(configManager)
	externalLinkRepository repositoryExternalLink.ExternalLinkRepository = repositoryExternalLink.NewExternalLinkRepository(database)
	externalLinkService    serviceExternalLink.ExternalLinkService       = serviceExternalLink.NewExternalLinkService(externalLinkRepository)
	externalLinkController controllerExternalLink.ExternalLinkController = controllerExternalLink.NewExternalLinkController(externalLinkService)

	httpRouter router.Router = router.NewMuxRouter()
)
