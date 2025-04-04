package controller

import (
	"main/config"
	cActivity "main/controller/activity"
	cActivityType "main/controller/activitytype"
	cAdoOrganization "main/controller/ado-organization"
	cApprovalType "main/controller/approval-type"
	cArticle "main/controller/article"
	cCategory "main/controller/category"
	cContributionArea "main/controller/contributionarea"
	cExternalLink "main/controller/externallink"
	cOssContributionSponsor "main/controller/osscontributionsponsor"
	cOtherRequest "main/controller/other-request"
	cRepositoryApprover "main/controller/repository-approver"
	cTopic "main/controller/topic"
	"main/service"
)

type Controller struct {
	Activity               cActivity.ActivityController
	ActivityType           cActivityType.ActivityTypeController
	AdoOrganization        cAdoOrganization.AdoOrganizationController
	ApprovalType           cApprovalType.ApprovalTypeController
	ContributionArea       cContributionArea.ContributionAreaController
	ExternalLink           cExternalLink.ExternalLinkController
	OssContributionSponsor cOssContributionSponsor.OSSContributionSponsorController
	OtherRequest           cOtherRequest.OtherRequestController
	RepositoryApprover     cRepositoryApprover.RepositoryApproverController
	Topic                  cTopic.TopicController
	Category               cCategory.CategoryController
	Article                cArticle.ArticleController
}

type ControllerOptionFunc func(*Controller)

func NewController(controllerOpts ...ControllerOptionFunc) *Controller {
	controller := &Controller{}

	for _, opt := range controllerOpts {
		opt(controller)
	}

	return controller
}

func NewActivityController(serv *service.Service) ControllerOptionFunc {
	return func(c *Controller) {
		c.Activity = cActivity.NewActivityController(serv)
	}
}

func NewActivityTypeController(serv *service.Service) ControllerOptionFunc {
	return func(c *Controller) {
		c.ActivityType = cActivityType.NewActivityTypeController(serv)
	}
}

func NewAdoOrganizationController(serv *service.Service, conf config.ConfigManager) ControllerOptionFunc {
	return func(c *Controller) {
		c.AdoOrganization = cAdoOrganization.NewAdoOrganizationController(serv, conf)
	}
}

func NewApprovalTypeController(serv *service.Service) ControllerOptionFunc {
	return func(c *Controller) {
		c.ApprovalType = cApprovalType.NewApprovalTypeController(serv)
	}
}

func NewContributionAreaController(serv *service.Service) ControllerOptionFunc {
	return func(c *Controller) {
		c.ContributionArea = cContributionArea.NewContributionAreaController(serv)
	}
}

func NewExternalLinkController(serv *service.Service) ControllerOptionFunc {
	return func(c *Controller) {
		c.ExternalLink = cExternalLink.NewExternalLinkController(serv)
	}
}

func NewOssContributionSponsorController(serv *service.Service) ControllerOptionFunc {
	return func(c *Controller) {
		c.OssContributionSponsor = cOssContributionSponsor.NewOssContributionSponsorController(serv)
	}
}

func NewOtherRequestController() ControllerOptionFunc {
	return func(c *Controller) {
		c.OtherRequest = cOtherRequest.NewOtherRequestController()
	}
}

func NewRepositoryApproverController(serv *service.Service, conf config.ConfigManager) ControllerOptionFunc {
	return func(c *Controller) {
		c.RepositoryApprover = cRepositoryApprover.NewRepositoryApproverController(serv, conf)
	}
}

func NewTopicController(serv *service.Service) ControllerOptionFunc {
	return func(c *Controller) {
		c.Topic = cTopic.NewTopicController(serv)
	}
}

func NewCategoryController(serv *service.Service) ControllerOptionFunc {
	return func(c *Controller) {
		c.Category = cCategory.NewCategoryController(serv)
	}
}

func NewArticleController(serv *service.Service) ControllerOptionFunc {
	return func(c *Controller) {
		c.Article = cArticle.NewArticleController(serv)
	}
}
