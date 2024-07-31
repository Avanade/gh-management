package main

import (
	"main/config"
	"main/repository"
	"main/router"

	repositoruContributionArea "main/repository/contributionarea"
	repositoryExternalLink "main/repository/externallink"

	serviceContributionArea "main/service/contributionarea"
	serviceExternalLink "main/service/externallink"

	controllerContributionArea "main/controller/contributionarea"
	controllerExternalLink "main/controller/externallink"
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

	httpRouter router.Router = router.NewMuxRouter()
)
