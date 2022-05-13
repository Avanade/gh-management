package models

import (
	"context"

	"github.com/qmuntal/stateless"
)

type Repository struct {
	ID          string
	Name        string
	Visibility  string
	URL         string
	Description string
}

type RepositoryState struct {
	state          stateless.State
	reviewApproved bool
	adminRights    bool
}

const (
	ReviewIsComplete   bool = true
	ReviewIsIncomplete bool = false
	IsAdmin            bool = true
	IsNotAdmin         bool = false
)

const (
	triggerInnerSourceToPrivate           = "InnerSourceToPrivate"
	triggerInnerSourceToInternal          = "InnerSourceToInternal"
	triggerInnerSourceRetireToDelete      = "InnerSourceToDelete"
	triggerInnerSourceRetireToArchive     = "InnerSourceToArchive"
	triggerInnerSourceArchiveToUnarchived = "InnerSourceArchiveToUnarchived"
	triggerInnerSourceToOpenSource        = "InnerSourceToOpenSource"
	triggerOpenSourceRetireToDelete       = "OpenSourceToDelete"
	triggerOpenSourceRetireToArchive      = "OpenSourceToArchive"
	triggerOpenSourceArchiveToUnarchived  = "OpenSourceArchiveToUnarchived"
	triggerOpenSourceToInnerSource        = "OpenSourceToInnerSource"
)

const (
	stateInnerSourceInternal = "InnerSourceInternal"
	stateInnerSourcePrivate  = "InnerSourcePrivate"
	stateInnerSourceDeleted  = "InnerSourceDeleted"
	stateInnerSourceArchived = "InnerSourceArchived"
	stateOpenSourcePublic    = "OpenSourcePublic"
	stateOpenSourceArchived  = "OpenSourceArchived"
	stateOpenSourceDeleted   = "OpenSourceDeleted"
)

func GetRepositoryStateMachine(initialState stateless.State, reviewApproved bool, adminRights bool) *stateless.StateMachine {
	repositoryContext := RepositoryState{
		state:          initialState,
		reviewApproved: reviewApproved,
		adminRights:    adminRights,
	}
	repositoryMachine := stateless.NewStateMachineWithExternalStorage(func(_ context.Context) (stateless.State, error) {
		return repositoryContext.state, nil
	}, func(_ context.Context, state stateless.State) error {
		repositoryContext.state = state
		return nil
	}, stateless.FiringImmediate)

	// InnerSource - Active
	repositoryMachine.Configure(stateInnerSourceInternal).
		Permit(triggerInnerSourceToPrivate, stateInnerSourcePrivate).
		Permit(triggerInnerSourceToOpenSource, stateOpenSourcePublic, func(_ context.Context, _ ...interface{}) bool {
			return repositoryContext.reviewApproved
		}).
		Permit(triggerInnerSourceRetireToArchive, stateInnerSourceArchived).
		Permit(triggerInnerSourceRetireToDelete, stateInnerSourceDeleted, func(_ context.Context, _ ...interface{}) bool {
			return repositoryContext.adminRights
		})
	repositoryMachine.Configure(stateInnerSourcePrivate).
		Permit(triggerInnerSourceToInternal, stateInnerSourceInternal).
		Permit(triggerInnerSourceRetireToArchive, stateInnerSourceArchived).
		Permit(triggerInnerSourceRetireToDelete, stateInnerSourceDeleted, func(_ context.Context, _ ...interface{}) bool {
			return repositoryContext.adminRights
		})

	// InnerSource - Retired
	repositoryMachine.Configure(stateInnerSourceDeleted)
	repositoryMachine.Configure(stateInnerSourceArchived).
		Permit(triggerInnerSourceArchiveToUnarchived, stateInnerSourceInternal).
		Permit(triggerInnerSourceRetireToDelete, stateInnerSourceDeleted, func(_ context.Context, _ ...interface{}) bool {
			return repositoryContext.adminRights
		})

	// OpenSource - Active
	repositoryMachine.Configure(stateOpenSourcePublic).
		Permit(triggerOpenSourceRetireToArchive, stateOpenSourceArchived).
		Permit(triggerOpenSourceToInnerSource, stateInnerSourceInternal, func(_ context.Context, _ ...interface{}) bool {
			return repositoryContext.adminRights
		}).
		Permit(triggerOpenSourceRetireToDelete, stateOpenSourceDeleted, func(_ context.Context, _ ...interface{}) bool {
			return repositoryContext.adminRights
		})

	// OpenSource - Retired
	repositoryMachine.Configure(stateOpenSourceArchived).
		Permit(triggerOpenSourceArchiveToUnarchived, stateOpenSourcePublic).
		Permit(triggerOpenSourceRetireToDelete, stateOpenSourceDeleted, func(_ context.Context, _ ...interface{}) bool {
			return repositoryContext.adminRights
		})
	repositoryMachine.Configure(stateOpenSourceDeleted)

	return repositoryMachine
}
