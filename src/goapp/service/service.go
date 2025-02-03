package service

import (
	"main/config"
	"main/repository"
	sActivity "main/service/activity"
	sActivityHelp "main/service/activityhelp"
	sActivityType "main/service/activitytype"
	sApprovalType "main/service/approval-type"
	sContributionArea "main/service/contributionarea"
	sEmail "main/service/email"
	sExternalLink "main/service/externallink"
	sOssContributionSponsor "main/service/osscontributionsponsor"
	sRepositoryApprover "main/service/repository-approver"
	sTopic "main/service/topic"
	sUser "main/service/user"
)

type Service struct {
	Activity               sActivity.ActivityService
	ActivityHelp           sActivityHelp.ActivityHelpService
	ActivityType           sActivityType.ActivityTypeService
	ApprovalType           sApprovalType.ApprovalTypeService
	ContributionArea       sContributionArea.ContributionAreaService
	Email                  sEmail.EmailService
	ExternalLink           sExternalLink.ExternalLinkService
	OssContributionSponsor sOssContributionSponsor.OssContributionSponsorService
	RepositoryApprover     sRepositoryApprover.RepositoryApproverService
	User                   sUser.UserService
	Topic                  sTopic.TopicService
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

func NewApprovalTypeService(repo *repository.Repository) ServiceOptionFunc {
	return func(s *Service) {
		s.ApprovalType = sApprovalType.NewApprovalTypeService(repo)
	}
}

func NewApproverService(repo *repository.Repository) ServiceOptionFunc {
	return func(s *Service) {
		s.RepositoryApprover = sRepositoryApprover.NewRepositoryApproverService(repo)
	}
}

func NewContributionAreaService(repo *repository.Repository) ServiceOptionFunc {
	return func(s *Service) {
		s.ContributionArea = sContributionArea.NewContributionAreaService(repo)
	}
}

func NewEmailService(conf config.ConfigManager) ServiceOptionFunc {
	return func(s *Service) {
		s.Email = sEmail.NewSdkEmailService(conf)
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

func NewUserService(repo *repository.Repository) ServiceOptionFunc {
	return func(s *Service) {
		s.User = sUser.NewUserService(repo)
	}
}

func NewTopicService(repo *repository.Repository) ServiceOptionFunc {
	return func(s *Service) {
		s.Topic = sTopic.NewTopicService(repo)
	}
}
