package main

import (
	"main/config"
	controllerExternalLink "main/controller/external-link"
	router "main/http"
	"main/repository"
	repositoryExternalLink "main/repository/external-link"
	serviceExternalLink "main/service/external-link"
)

var (
	configManager          config.ConfigManager                          = config.NewEnvConfigManager()
	database               repository.Database                           = repository.NewDatabase(configManager)
	externalLinkRepository repositoryExternalLink.ExternalLinkRepository = repositoryExternalLink.NewExternalLinkRepository(database)
	externalLinkService    serviceExternalLink.ExternalLinkService       = serviceExternalLink.NewExternalLinkService(externalLinkRepository)
	externalLinkController controllerExternalLink.ExternalLinkController = controllerExternalLink.NewExternalLinkController(externalLinkService)

	httpRouter router.Router = router.NewMuxRouter()
)
