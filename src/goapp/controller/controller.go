package controller

import (
	cActivity "main/controller/activity"
	cActivityType "main/controller/activitytype"
	cApprovalType "main/controller/approval-type"
	cContributionArea "main/controller/contributionarea"
	cExternalLink "main/controller/externallink"
	cOssContributionSponsor "main/controller/osscontributionsponsor"
	"main/service"
)

type Controller struct {
	Activity               cActivity.ActivityController
	ActivityType           cActivityType.ActivityTypeController
	ApprovalType           cApprovalType.ApprovalTypeController
	ContributionArea       cContributionArea.ContributionAreaController
	ExternalLink           cExternalLink.ExternalLinkController
	OssContributionSponsor cOssContributionSponsor.OSSContributionSponsorController
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
