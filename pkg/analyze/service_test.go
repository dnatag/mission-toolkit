package analyze

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewService(t *testing.T) {
	fs := afero.NewMemMapFs()
	svc := NewService(fs, "test-mission-123")

	assert.NotNil(t, svc)
	assert.Equal(t, "test-mission-123", svc.missionID)
}

func TestInitializePlan(t *testing.T) {
	fs := afero.NewMemMapFs()
	svc := NewService(fs, "test-mission-123")

	t.Run("creates new plan with intent", func(t *testing.T) {
		err := svc.InitializePlan("test intent")
		require.NoError(t, err)

		state, err := LoadState(fs, ".mission/plan.json")
		require.NoError(t, err)
		assert.Equal(t, "test intent", state.OriginalIntent)
	})
}

func TestGetPlanState(t *testing.T) {
	fs := afero.NewMemMapFs()
	svc := NewService(fs, "test-mission-123")

	t.Run("retrieves existing plan state", func(t *testing.T) {
		err := svc.InitializePlan("test intent")
		require.NoError(t, err)

		state, err := svc.GetPlanState()
		require.NoError(t, err)
		assert.Equal(t, "test intent", state.OriginalIntent)
	})

	t.Run("returns error for missing plan", func(t *testing.T) {
		newFs := afero.NewMemMapFs()
		newSvc := NewService(newFs, "test-mission-456")

		_, err := newSvc.GetPlanState()
		assert.Error(t, err)
	})
}

func TestUpdatePlanState(t *testing.T) {
	fs := afero.NewMemMapFs()
	svc := NewService(fs, "test-mission-123")

	t.Run("updates existing plan state", func(t *testing.T) {
		err := svc.InitializePlan("original intent")
		require.NoError(t, err)

		state, err := svc.GetPlanState()
		require.NoError(t, err)

		state.RefinedIntent = "refined intent"
		state.MissionType = "WET"
		state.Track = 2

		err = svc.UpdatePlanState(state)
		require.NoError(t, err)

		updated, err := svc.GetPlanState()
		require.NoError(t, err)
		assert.Equal(t, "original intent", updated.OriginalIntent)
		assert.Equal(t, "refined intent", updated.RefinedIntent)
		assert.Equal(t, "WET", updated.MissionType)
		assert.Equal(t, 2, updated.Track)
	})
}
