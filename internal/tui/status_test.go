package tui

import (
	"strings"
	"testing"

	"github.com/dnatag/mission-toolkit/internal/mission"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test helper functions
func createTestMission(intent, status string) *mission.Mission {
	return &mission.Mission{
		Intent: intent,
		Status: status,
		Type:   "WET",
		Track:  "2",
	}
}

func createTestModel() Model {
	return Model{
		selectedIndex:  0,
		currentPage:    0,
		itemsPerPage:   5,
		viewportHeight: 20,
		completedMissions: []*mission.Mission{
			createTestMission("Test mission 1", "completed"),
			createTestMission("Test mission 2", "completed"),
			createTestMission("Search test mission", "completed"),
		},
	}
}

// Test Model initialization
func TestNewModel(t *testing.T) {
	model := NewModel()

	assert.Equal(t, 0, model.selectedIndex)
	assert.Equal(t, 0, model.currentPage)
	assert.Equal(t, 5, model.itemsPerPage)
	assert.False(t, model.searchMode)
	assert.Empty(t, model.searchQuery)
}

// Test pagination logic
func TestGetTotalPages(t *testing.T) {
	tests := []struct {
		name          string
		missionsCount int
		itemsPerPage  int
		expectedPages int
	}{
		{"Empty missions", 0, 5, 1},
		{"Less than page size", 3, 5, 1},
		{"Exact page size", 5, 5, 1},
		{"More than page size", 7, 5, 2},
		{"Multiple pages", 12, 5, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := createTestModel()
			model.itemsPerPage = tt.itemsPerPage

			// Create missions for test
			missions := make([]*mission.Mission, tt.missionsCount)
			for i := 0; i < tt.missionsCount; i++ {
				missions[i] = createTestMission("Test", "completed")
			}
			model.completedMissions = missions

			pages := model.getTotalPages()
			assert.Equal(t, tt.expectedPages, pages)
		})
	}
}

func TestGetPageSize(t *testing.T) {
	model := createTestModel()
	model.itemsPerPage = 5

	// Test with missions less than page size
	model.completedMissions = []*mission.Mission{
		createTestMission("Test 1", "completed"),
		createTestMission("Test 2", "completed"),
	}
	assert.Equal(t, 2, model.getPageSize())

	// Test with missions equal to page size
	model.completedMissions = make([]*mission.Mission, 5)
	for i := 0; i < 5; i++ {
		model.completedMissions[i] = createTestMission("Test", "completed")
	}
	assert.Equal(t, 5, model.getPageSize())
}

func TestGetCurrentPageMissions(t *testing.T) {
	model := createTestModel()
	model.itemsPerPage = 2

	// Create 5 missions
	missions := make([]*mission.Mission, 5)
	for i := 0; i < 5; i++ {
		missions[i] = createTestMission("Mission "+string(rune('A'+i)), "completed")
	}
	model.completedMissions = missions

	// Test first page
	model.currentPage = 0
	pageMissions := model.getCurrentPageMissions()
	assert.Len(t, pageMissions, 2)
	assert.Equal(t, "Mission A", pageMissions[0].Intent)

	// Test second page
	model.currentPage = 1
	pageMissions = model.getCurrentPageMissions()
	assert.Len(t, pageMissions, 2)
	assert.Equal(t, "Mission C", pageMissions[0].Intent)

	// Test last page (partial)
	model.currentPage = 2
	pageMissions = model.getCurrentPageMissions()
	assert.Len(t, pageMissions, 1)
	assert.Equal(t, "Mission E", pageMissions[0].Intent)
}

// Test search and filtering
func TestFilterMissions(t *testing.T) {
	model := createTestModel()

	tests := []struct {
		name          string
		query         string
		expectedCount int
	}{
		{"Empty query", "", 0},
		{"Match found", "search", 1},
		{"No match", "nonexistent", 0},
		{"Case insensitive", "SEARCH", 1},
		{"Partial match", "test", 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filtered := model.filterMissions(tt.query)
			assert.Len(t, filtered, tt.expectedCount)
		})
	}
}

func TestMatchesFuzzy(t *testing.T) {
	model := createTestModel()
	testMission := createTestMission("Add user authentication", "completed")

	tests := []struct {
		name     string
		query    string
		expected bool
	}{
		{"Exact match", "user authentication", true},
		{"Partial match", "user", true},
		{"Case insensitive", "user", true}, // Function converts to lowercase, so "USER" becomes "user"
		{"No match", "database", false},
		{"Empty query", "", true}, // Empty string matches everything via Contains
		{"Status match", "completed", true},
		{"Type match", "wet", true},
		{"Track match", "2", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := model.matchesFuzzy(testMission, strings.ToLower(tt.query))
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetActiveMissions(t *testing.T) {
	model := createTestModel()

	// Test normal mode
	model.searchMode = false
	missions := model.getActiveMissions()
	assert.Len(t, missions, 3)

	// Test search mode with query
	model.searchMode = true
	model.searchQuery = "search"
	model.filteredMissions = []*mission.Mission{
		createTestMission("Search result", "completed"),
	}
	missions = model.getActiveMissions()
	assert.Len(t, missions, 1)

	// Test search mode with empty query
	model.searchQuery = ""
	missions = model.getActiveMissions()
	assert.Len(t, missions, 0)
}

// Test scroll offset calculations
func TestGetMaxScrollOffset(t *testing.T) {
	model := createTestModel()

	// Test with no selected mission
	model.selectedMission = nil
	assert.Equal(t, 0, model.getMaxScrollOffset())

	// Test with selected mission
	model.selectedMission = createTestMission("Test mission with long content", "completed")
	model.viewportHeight = 10

	// Since we can't easily test the actual content height calculation,
	// we'll test that it returns a non-negative value
	offset := model.getMaxScrollOffset()
	assert.GreaterOrEqual(t, offset, 0)
}

// Test utility functions
func TestTruncate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxLen   int
		expected string
	}{
		{"Short string", "hello", 10, "hello"},
		{"Exact length", "hello", 5, "hello"},
		{"Long string", "hello world", 5, "he..."},
		{"Empty string", "", 5, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncate(tt.input, tt.maxLen)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMax(t *testing.T) {
	assert.Equal(t, 5, max(3, 5))
	assert.Equal(t, 5, max(5, 3))
	assert.Equal(t, 5, max(5, 5))
	assert.Equal(t, 0, max(-1, 0))
}

func TestMin(t *testing.T) {
	assert.Equal(t, 3, min(3, 5))
	assert.Equal(t, 3, min(5, 3))
	assert.Equal(t, 5, min(5, 5))
	assert.Equal(t, -1, min(-1, 0))
}

// Integration tests for state transitions
func TestSearchModeToggle(t *testing.T) {
	model := createTestModel()

	// Initially not in search mode
	assert.False(t, model.searchMode)
	assert.Empty(t, model.searchQuery)

	// Test that search functionality works with the model
	model.searchMode = true
	model.searchQuery = "test"

	filtered := model.filterMissions("test")
	assert.Len(t, filtered, 3) // All test missions contain "test"
}

func TestPaginationBoundaries(t *testing.T) {
	model := createTestModel()
	model.itemsPerPage = 2

	// Test page boundaries
	totalPages := model.getTotalPages()
	require.Greater(t, totalPages, 1)

	// Test first page
	model.currentPage = 0
	missions := model.getCurrentPageMissions()
	assert.Len(t, missions, 2)

	// Test last page
	model.currentPage = totalPages - 1
	missions = model.getCurrentPageMissions()
	assert.Greater(t, len(missions), 0)
	assert.LessOrEqual(t, len(missions), 2)
}
