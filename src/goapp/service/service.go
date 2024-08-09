package service

import (
	"main/config"
	"main/repository"
	sActivity "main/service/activity"
	sActivityHelp "main/service/activityhelp"
	sActivityType "main/service/activitytype"
	sContributionArea "main/service/contributionarea"
	sEmail "main/service/email"
	sExternalLink "main/service/externallink"
	sOssContributionSponsor "main/service/osscontributionsponsor"
)

type Service struct {
	Activity               sActivity.ActivityService
	ActivityHelp           sActivityHelp.ActivityHelpService
	ActivityType           sActivityType.ActivityTypeService
	ContributionArea       sContributionArea.ContributionAreaService
	ExternalLink           sExternalLink.ExternalLinkService
	OssContributionSponsor sOssContributionSponsor.OssContributionSponsorService
	Email                  sEmail.EmailService
}

type ServiceOptionFunc func(*Service)

func NewService(serviceOpts ...ServiceOptionFunc) *Service {
	service := &Service{}

	for _, opt := range serviceOpts {
		opt(service)
	}

	return service
}

func NewActivityService(repo *repository.Repository) ServiceOptionFunc {
	return func(s *Service) {
		s.Activity = sActivity.NewActivityService(repo)
	}
}

func NewActivityHelpService(repo *repository.Repository) ServiceOptionFunc {
	return func(s *Service) {
		s.ActivityHelp = sActivityHelp.NewActivityHelpService(repo)
	}
}

func NewActivityTypeService(repo *repository.Repository) ServiceOptionFunc {
	return func(s *Service) {
		s.ActivityType = sActivityType.NewActivityTypeService(repo)
	}
}

func NewContributionAreaService(repo *repository.Repository) ServiceOptionFunc {
	return func(s *Service) {
		s.ContributionArea = sContributionArea.NewContributionAreaService(repo)
	}
}

func NewEmailService(conf config.ConfigManager) ServiceOptionFunc {
	return func(s *Service) {
		s.Email = sEmail.NewHttpEmailService(conf)
	}
}

func NewExternalLinkService(repo *repository.Repository) ServiceOptionFunc {
	return func(s *Service) {
		s.ExternalLink = sExternalLink.NewExternalLinkService(repo)
	}
}

func NewOssContributionSponsorService(repo *repository.Repository) ServiceOptionFunc {
	return func(s *Service) {
		s.OssContributionSponsor = sOssContributionSponsor.NewOssContributionSponsorService(repo)
	}
}
