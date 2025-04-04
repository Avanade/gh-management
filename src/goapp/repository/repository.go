package repository

import (
	"main/infrastructure/database"
	rActivity "main/repository/activity"
	rActivityContributionArea "main/repository/activitycontributionarea"
	rActivityHelp "main/repository/activityhelp"
	rActivityType "main/repository/activitytype"
	rAdoOrganization "main/repository/ado-organization"
	rAdoOrganizationApprovalRequest "main/repository/ado-organization-approval-request"
	rApprovalRequest "main/repository/approval-request"
	rApprovalType "main/repository/approval-type"
	rArticle "main/repository/article"
	rCategory "main/repository/category"
	rCommunityApprover "main/repository/community-approver"
	rContributionArea "main/repository/contributionarea"
	rExternalLink "main/repository/externallink"
	rOssContributionSponsor "main/repository/osscontributionsponsor"
	rRepositoryApprover "main/repository/repository-approver"
	rTopic "main/repository/topic"
	rUser "main/repository/user"
)

type Repository struct {
	Activity                       rActivity.ActivityRepository
	ActivityContributionArea       rActivityContributionArea.ActivityContributionAreaRepository
	ActivityHelp                   rActivityHelp.ActivityHelpRepository
	ActivityType                   rActivityType.ActivityTypeRepository
	AdoOrganization                rAdoOrganization.AdoOrganizationRepository
	AdoOrganizationApprovalRequest rAdoOrganizationApprovalRequest.AdoOrganizationApprovalRequestRepository
	ApprovalRequest                rApprovalRequest.ApprovalRequestRepository
	ApprovalType                   rApprovalType.ApprovalTypeRepository
	CommunityApprover              rCommunityApprover.CommunityApproverRepository
	ContributionArea               rContributionArea.ContributionAreaRepository
	ExternalLink                   rExternalLink.ExternalLinkRepository
	OssContributionSponsor         rOssContributionSponsor.OssContributionSponsorRepository
	RepositoryApprover             rRepositoryApprover.RepositoryApproverRepository
	User                           rUser.UserRepository
	Topic                          rTopic.TopicRepository
	Category                       rCategory.CategoryRepository
	Article                        rArticle.ArticleRepository
}

type RepositoryOptionFunc func(*Repository)

func NewRepository(repoOpts ...RepositoryOptionFunc) *Repository {
	repository := &Repository{}

	for _, opt := range repoOpts {
		opt(repository)
	}

	return repository
}

func NewActivity(db *database.Database) RepositoryOptionFunc {
	return func(r *Repository) {
		r.Activity = rActivity.NewActivityRepository(db)
	}
}

func NewActivityContributionArea(db *database.Database) RepositoryOptionFunc {
	return func(r *Repository) {
		r.ActivityContributionArea = rActivityContributionArea.NewActivityContributionAreaRepository(db)
	}
}

func NewActivityType(db *database.Database) RepositoryOptionFunc {
	return func(r *Repository) {
		r.ActivityType = rActivityType.NewActivityTypeRepository(db)
	}
}

func NewAdoOrganization(db *database.Database) RepositoryOptionFunc {
	return func(r *Repository) {
		r.AdoOrganization = rAdoOrganization.NewAdoOrganizationRepository(db)
	}
}

func NewAdoOrganizationApprovalRequest(db *database.Database) RepositoryOptionFunc {
	return func(r *Repository) {
		r.AdoOrganizationApprovalRequest = rAdoOrganizationApprovalRequest.NewAdoOrganizationApprovalRequestRepository(db)
	}
}

func NewApprovalRequest(db *database.Database) RepositoryOptionFunc {
	return func(r *Repository) {
		r.ApprovalRequest = rApprovalRequest.NewApprovalRequestRepository(db)
	}
}

func NewApprovalType(db *database.Database) RepositoryOptionFunc {
	return func(r *Repository) {
		r.ApprovalType = rApprovalType.NewApprovalTypeRepository(db)
	}
}

func NewApprover(db *database.Database) RepositoryOptionFunc {
	return func(r *Repository) {
		r.RepositoryApprover = rRepositoryApprover.NewRepostioryApproverRepository(db)
	}
}

func NewCommunityApprover(db *database.Database) RepositoryOptionFunc {
	return func(r *Repository) {
		r.CommunityApprover = rCommunityApprover.NewCommunityApproverRepository(db)
	}
}

func NewContributionArea(db *database.Database) RepositoryOptionFunc {
	return func(r *Repository) {
		r.ContributionArea = rContributionArea.NewContributionAreaRepository(db)
	}
}

func NewExternalLink(db *database.Database) RepositoryOptionFunc {
	return func(r *Repository) {
		r.ExternalLink = rExternalLink.NewExternalLinkRepository(db)
	}
}

func NewOssContributionSponsor(db *database.Database) RepositoryOptionFunc {
	return func(r *Repository) {
		r.OssContributionSponsor = rOssContributionSponsor.NewOSSContributionSponsorRepository(db)
	}
}

func NewActivityHelp(db *database.Database) RepositoryOptionFunc {
	return func(r *Repository) {
		r.ActivityHelp = rActivityHelp.NewActivityHelpRepository(db)
	}
}

func NewUser(db *database.Database) RepositoryOptionFunc {
	return func(r *Repository) {
		r.User = rUser.NewUserRepository(db)
	}
}

func NewPopularTopic(db *database.Database) RepositoryOptionFunc {
	return func(r *Repository) {
		r.Topic = rTopic.NewTopicRepository(db)
	}
}

func NewCategory(db *database.Database) RepositoryOptionFunc {
	return func(r *Repository) {
		r.Category = rCategory.NewCategoryRepository(db)
	}
}

func NewArticle(db *database.Database) RepositoryOptionFunc {
	return func(r *Repository) {
		r.Article = rArticle.NewArticleRepository(db)
	}
}
