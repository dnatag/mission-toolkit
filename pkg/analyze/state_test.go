package analyze

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadState(t *testing.T) {
	fs := afero.NewMemMapFs()

	t.Run("loads valid plan state", func(t *testing.T) {
		content := `{
  "original_intent": "test intent",
  "refined_intent": "refined test intent",
  "mission_type": "WET",
  "track": 2
}`
		err := afero.WriteFile(fs, ".mission/plan.json", []byte(content), 0644)
		require.NoError(t, err)

		state, err := LoadState(fs, ".mission/plan.json")
		require.NoError(t, err)
		assert.Equal(t, "test intent", state.OriginalIntent)
		assert.Equal(t, "refined test intent", state.RefinedIntent)
		assert.Equal(t, "WET", state.MissionType)
		assert.Equal(t, 2, state.Track)
	})

	t.Run("returns error for missing file", func(t *testing.T) {
		_, err := LoadState(fs, ".mission/missing.json")
		assert.Error(t, err)
	})

	t.Run("returns error for invalid JSON", func(t *testing.T) {
		err := afero.WriteFile(fs, ".mission/invalid.json", []byte("not json"), 0644)
		require.NoError(t, err)

		_, err = LoadState(fs, ".mission/invalid.json")
		assert.Error(t, err)
	})
}

func TestSaveState(t *testing.T) {
	fs := afero.NewMemMapFs()

	t.Run("saves plan state successfully", func(t *testing.T) {
		state := &PlanState{
			OriginalIntent: "test intent",
			RefinedIntent:  "refined",
			MissionType:    "WET",
			Track:          2,
			Scope:          []string{"file1.go", "file2.go"},
		}

		err := SaveState(fs, state, ".mission/plan.json")
		require.NoError(t, err)

		loaded, err := LoadState(fs, ".mission/plan.json")
		require.NoError(t, err)
		assert.Equal(t, state.OriginalIntent, loaded.OriginalIntent)
		assert.Equal(t, state.RefinedIntent, loaded.RefinedIntent)
		assert.Equal(t, state.MissionType, loaded.MissionType)
		assert.Equal(t, state.Track, loaded.Track)
		assert.Equal(t, state.Scope, loaded.Scope)
	})

	t.Run("creates directory if missing", func(t *testing.T) {
		state := &PlanState{
			OriginalIntent: "test",
		}

		err := SaveState(fs, state, ".mission/nested/plan.json")
		require.NoError(t, err)

		exists, err := afero.DirExists(fs, ".mission/nested")
		require.NoError(t, err)
		assert.True(t, exists)
	})
}
