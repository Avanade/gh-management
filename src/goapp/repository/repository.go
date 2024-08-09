package repository

import (
	"main/infrastructure/database"
	rActivity "main/repository/activity"
	rActivityContributionArea "main/repository/activitycontributionarea"
	rActivityHelp "main/repository/activityhelp"
	rActivityType "main/repository/activitytype"
	rContributionArea "main/repository/contributionarea"
	rExternalLink "main/repository/externallink"
	rOssContributionSponsor "main/repository/osscontributionsponsor"
)

type Repository struct {
	Activity                 rActivity.ActivityRepository
	ActivityContributionArea rActivityContributionArea.ActivityContributionAreaRepository
	ActivityHelp             rActivityHelp.ActivityHelpRepository
	ActivityType             rActivityType.ActivityTypeRepository
	ContributionArea         rContributionArea.ContributionAreaRepository
	ExternalLink             rExternalLink.ExternalLinkRepository
	OssContributionSponsor   rOssContributionSponsor.OssContributionSponsorRepository
}

type RepositoryOptionFunc func(*Repository)

func NewRepository(repoOpts ...RepositoryOptionFunc) *Repository {
	repository := &Repository{}

	for _, opt := range repoOpts {
		opt(repository)
	}

	return repository
}

func NewActivity(db database.Database) RepositoryOptionFunc {
	return func(r *Repository) {
		r.Activity = rActivity.NewActivityRepository(db)
	}
}

func NewActivityContributionArea(db database.Database) RepositoryOptionFunc {
	return func(r *Repository) {
		r.ActivityContributionArea = rActivityContributionArea.NewActivityContributionAreaRepository(db)
	}
}

func NewActivityType(db database.Database) RepositoryOptionFunc {
	return func(r *Repository) {
		r.ActivityType = rActivityType.NewActivityTypeRepository(db)
	}
}

func NewContributionArea(db database.Database) RepositoryOptionFunc {
	return func(r *Repository) {
		r.ContributionArea = rContributionArea.NewContributionAreaRepository(db)
	}
}

func NewExternalLink(db database.Database) RepositoryOptionFunc {
	return func(r *Repository) {
		r.ExternalLink = rExternalLink.NewExternalLinkRepository(db)
	}
}

func NewOssContributionSponsor(db database.Database) RepositoryOptionFunc {
	return func(r *Repository) {
		r.OssContributionSponsor = rOssContributionSponsor.NewOSSContributionSponsorRepository(db)
	}
}

func NewActivityHelp(db database.Database) RepositoryOptionFunc {
	return func(r *Repository) {
		r.ActivityHelp = rActivityHelp.NewActivityHelpRepository(db)
	}
}
