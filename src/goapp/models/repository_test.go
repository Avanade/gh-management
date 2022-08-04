package models

import (
	"context"
	"testing"
)

func TestNewStateMachine(t *testing.T) {
	t.Run("should return a state machine in innersource internal, unaffected by positive approval state", func(t *testing.T) {
		sm := GetRepositoryStateMachine(stateInnerSourceInternal, ReviewIsComplete, IsNotAdmin)

		got, err := sm.State(context.Background())
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		want := stateInnerSourceInternal
		if got != want {
			t.Errorf("Expected %s, got %s", want, got)
		}
	})
	t.Run("should return a state machine in innersource internal, unaffected by negative approval state", func(t *testing.T) {
		sm := GetRepositoryStateMachine(stateInnerSourceInternal, ReviewIsIncomplete, IsNotAdmin)

		got, err := sm.State(context.Background())
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		want := stateInnerSourceInternal
		if got != want {
			t.Errorf("Expected %s, got %s", want, got)
		}
	})
}

func TestTransitionToOpenSource(t *testing.T) {
	t.Run("cannot move without approval", func(t *testing.T) {
		sm := GetRepositoryStateMachine(stateInnerSourceInternal, ReviewIsIncomplete, IsNotAdmin)

		got, err := sm.CanFire(triggerInnerSourceToOpenSource)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		want := false
		if got != want {
			t.Errorf("Expected %v, got %v", want, got)
		}
	})
	t.Run("can move with approval", func(t *testing.T) {
		sm := GetRepositoryStateMachine(stateInnerSourceInternal, ReviewIsComplete, IsNotAdmin)

		got, err := sm.CanFire(triggerInnerSourceToOpenSource)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		want := true
		if got != want {
			t.Errorf("Expected %v, got %v", want, got)
		}
	})
}

// Test deletion is limited to admins
func TestDeletionRules(t *testing.T) {
	t.Run("only admin can delete - internal", func(t *testing.T) {
		sm := GetRepositoryStateMachine(stateInnerSourceInternal, ReviewIsIncomplete, IsNotAdmin)

		got, _ := sm.CanFire(triggerInnerSourceRetireToDelete)
		want := false
		if got != want {
			t.Errorf("Expected %v, got %v", want, got)
		}
	})
	t.Run("only admin can delete - private", func(t *testing.T) {
		sm := GetRepositoryStateMachine(stateInnerSourcePrivate, ReviewIsIncomplete, IsNotAdmin)

		got, _ := sm.CanFire(triggerInnerSourceRetireToDelete)
		want := false
		if got != want {
			t.Errorf("Expected %v, got %v", want, got)
		}
	})
	t.Run("only admin can delete - archived", func(t *testing.T) {
		sm := GetRepositoryStateMachine(stateInnerSourceArchived, ReviewIsIncomplete, IsNotAdmin)

		got, _ := sm.CanFire(triggerInnerSourceRetireToDelete)
		want := false
		if got != want {
			t.Errorf("Expected %v, got %v", want, got)
		}
	})
	t.Run("only admin can delete - oss", func(t *testing.T) {
		sm := GetRepositoryStateMachine(stateOpenSourcePublic, ReviewIsIncomplete, IsNotAdmin)

		got, _ := sm.CanFire(triggerInnerSourceRetireToDelete)
		want := false
		if got != want {
			t.Errorf("Expected %v, got %v", want, got)
		}
	})
	t.Run("only admin can delete - oss archive", func(t *testing.T) {
		sm := GetRepositoryStateMachine(stateOpenSourceArchived, ReviewIsIncomplete, IsNotAdmin)

		got, _ := sm.CanFire(triggerInnerSourceRetireToDelete)
		want := false
		if got != want {
			t.Errorf("Expected %v, got %v", want, got)
		}
	})
	t.Run("admin can delete - internal", func(t *testing.T) {
		sm := GetRepositoryStateMachine(stateInnerSourceInternal, ReviewIsIncomplete, IsAdmin)
		got, _ := sm.CanFire(triggerInnerSourceRetireToDelete)
		want := true
		if got != want {
			t.Errorf("Expected %v, got %v", want, got)
		}
	})
	t.Run("admin can delete - private", func(t *testing.T) {
		sm := GetRepositoryStateMachine(stateInnerSourcePrivate, ReviewIsIncomplete, IsAdmin)

		got, _ := sm.CanFire(triggerInnerSourceRetireToDelete)
		want := true
		if got != want {
			t.Errorf("Expected %v, got %v", want, got)
		}
	})
	t.Run("admin can delete - archive", func(t *testing.T) {
		sm := GetRepositoryStateMachine(stateInnerSourceArchived, ReviewIsIncomplete, IsAdmin)

		got, _ := sm.CanFire(triggerInnerSourceRetireToDelete)
		want := true
		if got != want {
			t.Errorf("Expected %v, got %v", want, got)
		}
	})
	t.Run("admin can delete - oss", func(t *testing.T) {
		sm := GetRepositoryStateMachine(stateOpenSourcePublic, ReviewIsIncomplete, IsAdmin)

		got, _ := sm.CanFire(triggerOpenSourceRetireToDelete)
		want := true
		if got != want {
			t.Errorf("Expected %v, got %v", want, got)
		}
	})
	t.Run("admin can delete - oss archive", func(t *testing.T) {
		sm := GetRepositoryStateMachine(stateOpenSourceArchived, ReviewIsIncomplete, IsAdmin)

		got, _ := sm.CanFire(triggerOpenSourceRetireToDelete)
		want := true
		if got != want {
			t.Errorf("Expected %v, got %v", want, got)
		}
	})
}
