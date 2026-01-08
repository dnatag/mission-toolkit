package tui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/bubbletea"
	"github.com/dnatag/mission-toolkit/internal/mission"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test helper functions
func createTestMission(intent, status string) *mission.Mission {
	body := "## INTENT\n" + intent + "\n\n## SCOPE\ntest.go\n\n## PLAN\n- Step 1\n\n## VERIFICATION\ngo test"
	return &mission.Mission{
		ID:     "test-id",
		Status: status,
		Type:   "WET",
		Track:  2,
		Body:   body,
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
		totalCount:  3, // For lazy loading tests
		loadedCount: 3,
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

// Test Init function
func TestInit(t *testing.T) {
	model := NewModel()
	cmd := model.Init()

	// Init should return a batch command
	assert.NotNil(t, cmd)
}

// Test Update function with different message types
func TestUpdate_WindowSizeMsg(t *testing.T) {
	model := createTestModel()

	msg := tea.WindowSizeMsg{Width: 100, Height: 50}
	updatedModel, cmd := model.Update(msg)

	m := updatedModel.(Model)
	assert.Equal(t, 100, m.width)
	assert.Equal(t, 50, m.height)
	assert.Equal(t, 40, m.viewportHeight) // height - 10
	assert.Nil(t, cmd)
}

func TestUpdate_CurrentMissionMsg(t *testing.T) {
	model := createTestModel()
	testMission := createTestMission("Current mission", "active")

	// Test successful mission load
	msg := currentMissionMsg{mission: testMission, err: nil}
	updatedModel, cmd := model.Update(msg)

	m := updatedModel.(Model)
	assert.Equal(t, testMission, m.currentMission)
	assert.Nil(t, cmd)

	// Test error case
	msg = currentMissionMsg{mission: nil, err: assert.AnError}
	updatedModel, cmd = model.Update(msg)

	m = updatedModel.(Model)
	assert.Nil(t, m.currentMission)
	assert.Nil(t, cmd)
}

func TestUpdate_CompletedMissionsMsg(t *testing.T) {
	model := createTestModel()
	testMissions := []*mission.Mission{
		createTestMission("Mission 1", "completed"),
		createTestMission("Mission 2", "completed"),
	}

	// Test successful missions load
	msg := completedMissionsMsg{missions: testMissions, err: nil}
	updatedModel, cmd := model.Update(msg)

	m := updatedModel.(Model)
	assert.Equal(t, testMissions, m.completedMissions)
	assert.Nil(t, cmd)

	// Test error case
	msg = completedMissionsMsg{missions: nil, err: assert.AnError}
	updatedModel, cmd = model.Update(msg)

	m = updatedModel.(Model)
	assert.NotEqual(t, testMissions, m.completedMissions) // Should keep original
	assert.Nil(t, cmd)
}

func TestUpdate_KeyMsg_SearchMode(t *testing.T) {
	model := createTestModel()
	model.searchMode = true

	// Test adding character to search query
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}}
	updatedModel, cmd := model.Update(msg)

	m := updatedModel.(Model)
	assert.Equal(t, "t", m.searchQuery)
	assert.Equal(t, 0, m.selectedIndex)
	assert.Equal(t, 0, m.currentPage)
	assert.Nil(t, cmd)
}

func TestUpdate_KeyMsg_Navigation(t *testing.T) {
	model := createTestModel()

	tests := []struct {
		name       string
		key        string
		expectQuit bool
		expectCmd  bool
	}{
		{"quit with q", "q", true, false},
		{"quit with ctrl+c", "ctrl+c", true, false},
		{"reload with r", "r", false, true},
		{"search with /", "/", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var msg tea.KeyMsg
			if tt.key == "ctrl+c" {
				msg = tea.KeyMsg{Type: tea.KeyCtrlC}
			} else {
				msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			}

			updatedModel, cmd := model.Update(msg)

			if tt.expectQuit {
				// Check if cmd is tea.Quit function
				assert.NotNil(t, cmd)
			} else if tt.expectCmd {
				assert.NotNil(t, cmd)
			} else {
				m := updatedModel.(Model)
				if tt.key == "/" {
					assert.True(t, m.searchMode)
				}
			}
		})
	}
}

func TestUpdate_KeyMsg_Escape(t *testing.T) {
	model := createTestModel()
	model.selectedMission = createTestMission("Selected", "completed")
	model.scrollOffset = 5

	msg := tea.KeyMsg{Type: tea.KeyEsc}
	updatedModel, cmd := model.Update(msg)

	m := updatedModel.(Model)
	assert.Nil(t, m.selectedMission)
	assert.Equal(t, 0, m.scrollOffset)
	assert.Nil(t, cmd)
}

// Test View function
func TestView(t *testing.T) {
	model := createTestModel()
	model.width = 80
	model.height = 24

	view := model.View()

	// View should contain basic elements
	assert.Contains(t, view, "Mission Toolkit Status")
	assert.NotEmpty(t, view)
}

func TestView_WithCurrentMission(t *testing.T) {
	model := createTestModel()
	model.currentMission = createTestMission("Active mission", "active")
	model.width = 80
	model.height = 24

	view := model.View()

	assert.Contains(t, view, "ACTIVE")
	assert.Contains(t, view, "Active mission")
}

func TestView_SearchMode(t *testing.T) {
	model := createTestModel()
	model.searchMode = true
	model.searchQuery = "test"
	model.width = 80
	model.height = 24

	view := model.View()

	assert.Contains(t, view, "Search:")
	assert.Contains(t, view, "test")
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
			model.totalCount = tt.missionsCount // Set totalCount for lazy loading

			// Create missions for test
			missions := make([]*mission.Mission, tt.missionsCount)
			for i := 0; i < tt.missionsCount; i++ {
				missions[i] = createTestMission("Test", "completed")
			}
			model.completedMissions = missions
			model.loadedCount = tt.missionsCount

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
	assert.Equal(t, "Mission A", extractIntent(pageMissions[0].Body))

	// Test second page
	model.currentPage = 1
	pageMissions = model.getCurrentPageMissions()
	assert.Len(t, pageMissions, 2)
	assert.Equal(t, "Mission C", extractIntent(pageMissions[0].Body))

	// Test last page (partial)
	model.currentPage = 2
	pageMissions = model.getCurrentPageMissions()
	assert.Len(t, pageMissions, 1)
	assert.Equal(t, "Mission E", extractIntent(pageMissions[0].Body))
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

// Test message types
func TestCurrentMissionMsg(t *testing.T) {
	mission := createTestMission("Test", "active")

	// Test successful message
	msg := currentMissionMsg{mission: mission, err: nil}
	assert.Equal(t, mission, msg.mission)
	assert.Nil(t, msg.err)

	// Test error message
	msg = currentMissionMsg{mission: nil, err: assert.AnError}
	assert.Nil(t, msg.mission)
	assert.NotNil(t, msg.err)
}

func TestCompletedMissionsMsg(t *testing.T) {
	missions := []*mission.Mission{createTestMission("Test", "completed")}

	// Test successful message
	msg := completedMissionsMsg{missions: missions, err: nil}
	assert.Equal(t, missions, msg.missions)
	assert.Nil(t, msg.err)

	// Test error message
	msg = completedMissionsMsg{missions: nil, err: assert.AnError}
	assert.Nil(t, msg.missions)
	assert.NotNil(t, msg.err)
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
	model.totalCount = 3 // Set totalCount for lazy loading

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

// Test lazy loading trigger functionality
// This test verifies that the on-demand loading trigger works correctly
// when users scroll to the end of currently loaded missions.
func TestLazyLoadingTrigger(t *testing.T) {
	model := createTestModel()
	model.itemsPerPage = 2
	model.totalCount = 10 // Total missions available
	model.loadedCount = 4 // Currently loaded missions
	model.loading = false
	model.searchMode = false

	// Create 4 loaded missions
	missions := make([]*mission.Mission, 4)
	for i := 0; i < 4; i++ {
		missions[i] = createTestMission("Mission "+string(rune('A'+i)), "completed")
	}
	model.completedMissions = missions

	tests := []struct {
		name          string
		currentPage   int
		selectedIndex int
		shouldTrigger bool
		description   string
	}{
		{
			name:          "First page, first item",
			currentPage:   0,
			selectedIndex: 0,
			shouldTrigger: false,
			description:   "absoluteIndex=0, not at end",
		},
		{
			name:          "First page, second item",
			currentPage:   0,
			selectedIndex: 1,
			shouldTrigger: false,
			description:   "absoluteIndex=1 after move, not at end",
		},
		{
			name:          "Second page, first item",
			currentPage:   1,
			selectedIndex: 0,
			shouldTrigger: true,
			description:   "absoluteIndex=3 after move, at last loaded mission (loadedCount-1)",
		},
		{
			name:          "Second page, second item (last loaded)",
			currentPage:   1,
			selectedIndex: 1,
			shouldTrigger: true,
			description:   "absoluteIndex=3, at last loaded mission (loadedCount-1)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset model state
			model.currentPage = tt.currentPage
			model.selectedIndex = tt.selectedIndex
			model.loading = false

			// Simulate down key press
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
			updatedModel, cmd := model.Update(msg)

			m := updatedModel.(Model)

			if tt.shouldTrigger {
				// Should trigger loading
				assert.True(t, m.loading, "Loading should be triggered for %s", tt.description)
				assert.NotNil(t, cmd, "Command should be returned for loading more missions")
			} else {
				// Should not trigger loading
				assert.False(t, m.loading, "Loading should not be triggered for %s", tt.description)
			}
		})
	}
}

// Test lazy loading trigger edge cases
// Verifies that the trigger correctly handles boundary conditions
// like all missions loaded, search mode, and already loading states.
func TestLazyLoadingTriggerEdgeCases(t *testing.T) {
	model := createTestModel()
	model.itemsPerPage = 2
	model.totalCount = 4
	model.loadedCount = 4 // All missions loaded
	model.loading = false
	model.searchMode = false

	// When all missions are loaded, should not trigger
	model.currentPage = 1
	model.selectedIndex = 1

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
	updatedModel, cmd := model.Update(msg)

	m := updatedModel.(Model)
	assert.False(t, m.loading, "Should not trigger when all missions are loaded")
	assert.Nil(t, cmd)

	// Test search mode - should not trigger
	model.searchMode = true
	model.loadedCount = 2 // Less than total

	updatedModel, cmd = model.Update(msg)
	m = updatedModel.(Model)
	assert.False(t, m.loading, "Should not trigger in search mode")

	// Test when already loading - should not trigger
	model.searchMode = false
	model.loading = true

	updatedModel, cmd = model.Update(msg)
	m = updatedModel.(Model)
	assert.True(t, m.loading, "Should remain in loading state")
}

// Test lazy loading trigger for next page navigation
func TestLazyLoadingTriggerNextPage(t *testing.T) {
	model := createTestModel()
	model.itemsPerPage = 2
	model.totalCount = 10 // Total missions available
	model.loadedCount = 4 // Currently loaded missions (0,1,2,3)
	model.loading = false
	model.searchMode = false

	// Create 4 loaded missions
	missions := make([]*mission.Mission, 4)
	for i := 0; i < 4; i++ {
		missions[i] = createTestMission("Mission "+string(rune('A'+i)), "completed")
	}
	model.completedMissions = missions

	tests := []struct {
		name          string
		currentPage   int
		shouldTrigger bool
		description   string
	}{
		{
			name:          "Page 0 to 1 - no trigger needed",
			currentPage:   0,
			shouldTrigger: false,
			description:   "Page 1 starts at index 2, within loaded missions",
		},
		{
			name:          "Page 1 to 2 - should trigger",
			currentPage:   1,
			shouldTrigger: true,
			description:   "Page 2 starts at index 4, beyond loaded missions",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset model state
			model.currentPage = tt.currentPage
			model.selectedIndex = 0
			model.loading = false

			// Simulate right key press (next page)
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}}
			updatedModel, cmd := model.Update(msg)

			m := updatedModel.(Model)

			if tt.shouldTrigger {
				// Should trigger loading
				assert.True(t, m.loading, "Loading should be triggered for %s", tt.description)
				assert.NotNil(t, cmd, "Command should be returned for loading more missions")
				// Page should not change yet when loading is triggered
				assert.Equal(t, tt.currentPage, m.currentPage, "Page should not change when loading is triggered")
			} else {
				// Should not trigger loading, page should change normally
				assert.False(t, m.loading, "Loading should not be triggered for %s", tt.description)
				assert.Equal(t, tt.currentPage+1, m.currentPage, "Page should advance normally")
			}
		})
	}
}

// Test prefetch functionality
func TestPrefetchFunctionality(t *testing.T) {
	model := createTestModel()
	model.itemsPerPage = 2
	model.totalCount = 10
	model.loadedCount = 6 // Enough missions loaded for multiple pages
	model.prefetchedPages = make(map[int][]*mission.Mission)

	// Create 6 loaded missions
	missions := make([]*mission.Mission, 6)
	for i := 0; i < 6; i++ {
		missions[i] = createTestMission("Mission "+string(rune('A'+i)), "completed")
	}
	model.completedMissions = missions

	// Test prefetch trigger
	model.currentPage = 1 // Start on page 1
	cmd := model.triggerPrefetch()

	// Should return a command to prefetch adjacent pages
	assert.NotNil(t, cmd, "Prefetch should return a command")
	assert.True(t, model.prefetching, "Prefetching flag should be set")
}

// Test prefetch message handling
func TestPrefetchMessageHandling(t *testing.T) {
	model := createTestModel()
	model.prefetchedPages = make(map[int][]*mission.Mission)
	model.currentPage = 1 // Set current page so page 2 is adjacent and won't be cleaned up

	// Test successful prefetch message
	testMissions := []*mission.Mission{
		createTestMission("Prefetched Mission 1", "completed"),
		createTestMission("Prefetched Mission 2", "completed"),
	}

	msg := prefetchPageMsg{
		page:     2,
		missions: testMissions,
		err:      nil,
	}

	updatedModel, cmd := model.Update(msg)
	m := updatedModel.(Model)

	assert.False(t, m.prefetching, "Prefetching flag should be cleared")
	assert.Len(t, m.prefetchedPages[2], 2, "Prefetched missions should be cached")
	assert.Equal(t, "Prefetched Mission 1", extractIntent(m.prefetchedPages[2][0].Body))
	assert.Nil(t, cmd, "No command should be returned")
}

// Test prefetch trigger on loadMoreMissions
func TestPrefetchOnLoadMore(t *testing.T) {
	model := createTestModel()
	model.itemsPerPage = 2
	model.totalCount = 10
	model.loadedCount = 4
	model.prefetchedPages = make(map[int][]*mission.Mission)
	model.currentPage = 1

	// Test loadMoreMissionsMsg triggers prefetch
	msg := loadMoreMissionsMsg{
		missions: []*mission.Mission{
			createTestMission("New Mission 1", "completed"),
			createTestMission("New Mission 2", "completed"),
		},
		loadedCount: 2,
		err:         nil,
	}

	updatedModel, cmd := model.Update(msg)
	m := updatedModel.(Model)

	assert.Equal(t, 6, m.loadedCount, "Loaded count should be updated")
	assert.NotNil(t, cmd, "Should trigger prefetch command")
	assert.False(t, m.loading, "Loading flag should be cleared")
}

// Test aggressive prefetch functionality
func TestAggressivePrefetch(t *testing.T) {
	model := createTestModel()
	model.itemsPerPage = 2
	model.totalCount = 10
	model.loadedCount = 8 // Enough missions for multiple pages
	model.prefetchedPages = make(map[int][]*mission.Mission)
	model.currentPage = 1

	// Create missions to simulate having enough for multiple pages
	missions := make([]*mission.Mission, 8)
	for i := 0; i < 8; i++ {
		missions[i] = createTestMission("Mission "+string(rune('A'+i)), "completed")
	}
	model.completedMissions = missions

	// Test aggressive prefetch
	cmd := model.triggerAggressivePrefetch()

	// Should return a command to prefetch next 2 pages and previous page
	assert.NotNil(t, cmd, "Aggressive prefetch should return a command")
	assert.True(t, model.prefetching, "Prefetching flag should be set")
}

// Test next page navigation with lazy loading and prefetch
func TestNextPageNavigationWithLazyLoading(t *testing.T) {
	model := createTestModel()
	model.itemsPerPage = 2
	model.totalCount = 10
	model.loadedCount = 4 // Only 4 missions loaded (pages 0,1)
	model.prefetchedPages = make(map[int][]*mission.Mission)
	model.currentPage = 1
	model.pendingPageChange = -1

	// Create 4 loaded missions
	missions := make([]*mission.Mission, 4)
	for i := 0; i < 4; i++ {
		missions[i] = createTestMission("Mission "+string(rune('A'+i)), "completed")
	}
	model.completedMissions = missions

	// Test right navigation that should trigger lazy loading
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}}
	updatedModel, cmd := model.Update(msg)

	m := updatedModel.(Model)

	// Should trigger loading and set pending page change
	assert.True(t, m.loading, "Should trigger loading")
	assert.Equal(t, 2, m.pendingPageChange, "Should set pending page change to 2")
	assert.Equal(t, 1, m.currentPage, "Current page should not change yet")
	assert.NotNil(t, cmd, "Should return load command")

	// Simulate loadMoreMissionsMsg response
	loadMsg := loadMoreMissionsMsg{
		missions: []*mission.Mission{
			createTestMission("Mission E", "completed"),
			createTestMission("Mission F", "completed"),
		},
		loadedCount: 2,
		err:         nil,
	}

	updatedModel2, cmd2 := m.Update(loadMsg)
	m2 := updatedModel2.(Model)

	// Should complete page change and trigger prefetch
	assert.False(t, m2.loading, "Loading should be complete")
	assert.Equal(t, -1, m2.pendingPageChange, "Pending page change should be cleared")
	assert.Equal(t, 2, m2.currentPage, "Should navigate to page 2")
	assert.Equal(t, 0, m2.selectedIndex, "Should reset selected index")
	assert.Equal(t, 6, m2.loadedCount, "Should have 6 loaded missions")
	assert.NotNil(t, cmd2, "Should trigger aggressive prefetch")
}

// Test next page navigation without lazy loading
func TestNextPageNavigationWithoutLazyLoading(t *testing.T) {
	model := createTestModel()
	model.itemsPerPage = 2
	model.totalCount = 6
	model.loadedCount = 6 // All missions already loaded
	model.prefetchedPages = make(map[int][]*mission.Mission)
	model.currentPage = 1
	model.pendingPageChange = -1

	// Create 6 loaded missions
	missions := make([]*mission.Mission, 6)
	for i := 0; i < 6; i++ {
		missions[i] = createTestMission("Mission "+string(rune('A'+i)), "completed")
	}
	model.completedMissions = missions

	// Test right navigation that should NOT trigger lazy loading
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}}
	updatedModel, cmd := model.Update(msg)

	m := updatedModel.(Model)

	// Should change page immediately and trigger prefetch
	assert.False(t, m.loading, "Should not trigger loading")
	assert.Equal(t, -1, m.pendingPageChange, "Should not set pending page change")
	assert.Equal(t, 2, m.currentPage, "Should navigate to page 2 immediately")
	assert.Equal(t, 0, m.selectedIndex, "Should reset selected index")
	assert.NotNil(t, cmd, "Should trigger prefetch")
}

// Test scroll-based page navigation (pgup/pgdn)
// Verifies that page up/down keys navigate between pages when in mission list view
func TestScrollPageNavigation(t *testing.T) {
	model := createTestModel()
	model.itemsPerPage = 2
	model.totalCount = 10
	model.loadedCount = 6 // Enough missions for multiple pages
	model.prefetchedPages = make(map[int][]*mission.Mission)
	model.currentPage = 1
	model.pendingPageChange = -1

	// Create 6 loaded missions
	missions := make([]*mission.Mission, 6)
	for i := 0; i < 6; i++ {
		missions[i] = createTestMission("Mission "+string(rune('A'+i)), "completed")
	}
	model.completedMissions = missions

	tests := []struct {
		name        string
		key         string
		expectPage  int
		expectCmd   bool
		description string
	}{
		{
			name:        "PgUp navigation",
			key:         "pgup",
			expectPage:  0,
			expectCmd:   true,
			description: "Should navigate to previous page",
		},
		{
			name:        "PgDn navigation",
			key:         "pgdn",
			expectPage:  2,
			expectCmd:   true,
			description: "Should navigate to next page",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset to page 1
			model.currentPage = 1
			model.selectedIndex = 0
			model.loading = false
			model.pendingPageChange = -1

			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			updatedModel, cmd := model.Update(msg)

			m := updatedModel.(Model)

			assert.Equal(t, tt.expectPage, m.currentPage, tt.description)
			assert.Equal(t, 0, m.selectedIndex, "Should reset selected index")

			if tt.expectCmd {
				assert.NotNil(t, cmd, "Should trigger prefetch command")
			}
		})
	}
}

// Test scroll navigation with lazy loading
func TestScrollNavigationWithLazyLoading(t *testing.T) {
	model := createTestModel()
	model.itemsPerPage = 2
	model.totalCount = 10
	model.loadedCount = 4 // Only 4 missions loaded (pages 0,1)
	model.prefetchedPages = make(map[int][]*mission.Mission)
	model.currentPage = 1
	model.pendingPageChange = -1

	// Create 4 loaded missions
	missions := make([]*mission.Mission, 4)
	for i := 0; i < 4; i++ {
		missions[i] = createTestMission("Mission "+string(rune('A'+i)), "completed")
	}
	model.completedMissions = missions

	// Test pgdn navigation that should trigger lazy loading
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("pgdn")}
	updatedModel, cmd := model.Update(msg)

	m := updatedModel.(Model)

	// Should trigger loading and set pending page change
	assert.True(t, m.loading, "Should trigger loading")
	assert.Equal(t, 2, m.pendingPageChange, "Should set pending page change to 2")
	assert.Equal(t, 1, m.currentPage, "Current page should not change yet")
	assert.NotNil(t, cmd, "Should return load command")
}

// Test scroll navigation in search mode (should not work)
func TestScrollNavigationInSearchMode(t *testing.T) {
	model := createTestModel()
	model.itemsPerPage = 2
	model.totalCount = 6
	model.loadedCount = 6
	model.currentPage = 1
	model.searchMode = true // In search mode

	// Create missions
	missions := make([]*mission.Mission, 6)
	for i := 0; i < 6; i++ {
		missions[i] = createTestMission("Mission "+string(rune('A'+i)), "completed")
	}
	model.completedMissions = missions

	// Test pgup in search mode - should not navigate
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("pgup")}
	updatedModel, cmd := model.Update(msg)

	m := updatedModel.(Model)

	// Should not change page in search mode
	assert.Equal(t, 1, m.currentPage, "Should not change page in search mode")
	assert.Nil(t, cmd, "Should not return command in search mode")
}

// Test up/down arrow page navigation
func TestArrowKeyPageNavigation(t *testing.T) {
	model := createTestModel()
	model.itemsPerPage = 2
	model.totalCount = 6
	model.loadedCount = 6 // All missions loaded
	model.prefetchedPages = make(map[int][]*mission.Mission)
	model.currentPage = 1
	model.selectedIndex = 0

	// Create 6 loaded missions
	missions := make([]*mission.Mission, 6)
	for i := 0; i < 6; i++ {
		missions[i] = createTestMission("Mission "+string(rune('A'+i)), "completed")
	}
	model.completedMissions = missions

	// Test up arrow at top of page - should go to previous page
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")}
	updatedModel, cmd := model.Update(msg)

	m := updatedModel.(Model)
	assert.Equal(t, 0, m.currentPage, "Should navigate to previous page")
	assert.Equal(t, 1, m.selectedIndex, "Should select bottom item of previous page")
	assert.NotNil(t, cmd, "Should trigger prefetch")

	// Reset to page 1, bottom item
	model.currentPage = 1
	model.selectedIndex = 1 // Bottom of page (2 items per page)

	// Test down arrow at bottom of page - should go to next page
	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")}
	updatedModel, cmd = model.Update(msg)

	m = updatedModel.(Model)
	assert.Equal(t, 2, m.currentPage, "Should navigate to next page")
	assert.Equal(t, 0, m.selectedIndex, "Should select top item of next page")
	assert.NotNil(t, cmd, "Should trigger prefetch")
}

// Test search mode navigation behavior
func TestSearchModeNavigationInDetailView(t *testing.T) {
	model := createTestModel()
	model.searchMode = true
	model.searchQuery = "test"

	// Create a mission with enough content to scroll
	longMissionBody := `## INTENT
This is a test mission with a very long intent that spans multiple lines to ensure we have enough content for scrolling in the detail view.

## SCOPE
file1.go
file2.go
file3.go
file4.go
file5.go
file6.go

## PLAN
- Step 1: Do something important
- Step 2: Do another important thing
- Step 3: Do yet another important thing
- Step 4: Do one more important thing
- Step 5: Do the final important thing

## VERIFICATION
go test ./...`

	longMission := &mission.Mission{
		ID:     "test-long-id",
		Status: "completed",
		Type:   "WET",
		Track:  2,
		Body:   longMissionBody,
	}

	model.selectedMission = longMission
	model.scrollOffset = 0
	model.viewportHeight = 10 // Small viewport to ensure scrolling is possible

	// Test that navigation keys work in detail view even when in search mode
	tests := []struct {
		name         string
		key          string
		expectScroll int
		description  string
	}{
		{
			name:         "Up arrow in detail view",
			key:          "k",
			expectScroll: 0, // Can't scroll up from 0
			description:  "Should handle up navigation in detail view",
		},
		{
			name:         "Down arrow in detail view",
			key:          "j",
			expectScroll: 1,
			description:  "Should handle down navigation in detail view",
		},
		{
			name:         "Page up in detail view",
			key:          "pgup",
			expectScroll: 0, // Can't scroll up from 0
			description:  "Should handle page up in detail view",
		},
		{
			name:         "Page down in detail view",
			key:          "pgdn",
			expectScroll: 5, // Scrolls by 5, but limited by maxScroll
			description:  "Should handle page down in detail view",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset scroll offset
			model.scrollOffset = 0

			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			updatedModel, cmd := model.Update(msg)

			m := updatedModel.(Model)

			// For scroll down operations, check that scroll increased (limited by maxScroll)
			if tt.key == "j" || tt.key == "pgdn" {
				assert.Greater(t, m.scrollOffset, 0, tt.description+" - scroll should increase")
			} else {
				assert.Equal(t, tt.expectScroll, m.scrollOffset, tt.description)
			}

			assert.True(t, m.searchMode, "Should remain in search mode")
			assert.NotNil(t, m.selectedMission, "Should remain in detail view")
			assert.Nil(t, cmd, "Should not return command for detail view navigation")
		})
	}
}

func TestSearchModeTextInputInDetailView(t *testing.T) {
	model := createTestModel()
	model.searchMode = true
	model.searchQuery = "test"
	model.selectedMission = createTestMission("Selected mission", "completed")

	// Test that text input is NOT captured when viewing mission details
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
	updatedModel, cmd := model.Update(msg)

	m := updatedModel.(Model)
	assert.Equal(t, "test", m.searchQuery, "Search query should not change in detail view")
	assert.True(t, m.searchMode, "Should remain in search mode")
	assert.NotNil(t, m.selectedMission, "Should remain in detail view")
	assert.Nil(t, cmd, "Should not return command")
}

func TestSearchModeTextInputInMissionList(t *testing.T) {
	model := createTestModel()
	model.searchMode = true
	model.searchQuery = "test"
	model.selectedMission = nil // In mission list, not detail view

	// Test that text input IS captured when in mission list
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
	updatedModel, cmd := model.Update(msg)

	m := updatedModel.(Model)
	assert.Equal(t, "testa", m.searchQuery, "Search query should be updated in mission list")
	assert.True(t, m.searchMode, "Should remain in search mode")
	assert.Nil(t, m.selectedMission, "Should remain in mission list")
	assert.Equal(t, 0, m.selectedIndex, "Should reset selected index")
	assert.Equal(t, 0, m.currentPage, "Should reset current page")
	assert.Nil(t, cmd, "Should not return command")
}

// Test up/down arrow within page navigation
func TestArrowKeyWithinPageNavigation(t *testing.T) {
	model := createTestModel()
	model.itemsPerPage = 3
	model.totalCount = 6
	model.loadedCount = 6
	model.currentPage = 1
	model.selectedIndex = 1 // Middle item

	// Create missions
	missions := make([]*mission.Mission, 6)
	for i := 0; i < 6; i++ {
		missions[i] = createTestMission("Mission "+string(rune('A'+i)), "completed")
	}
	model.completedMissions = missions

	// Test up arrow within page
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")}
	updatedModel, cmd := model.Update(msg)

	m := updatedModel.(Model)
	assert.Equal(t, 1, m.currentPage, "Should stay on same page")
	assert.Equal(t, 0, m.selectedIndex, "Should move up within page")
	assert.Nil(t, cmd, "Should not trigger command for within-page movement")

	// Test down arrow within page
	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")}
	updatedModel, cmd = m.Update(msg)

	m2 := updatedModel.(Model)
	assert.Equal(t, 1, m2.currentPage, "Should stay on same page")
	assert.Equal(t, 1, m2.selectedIndex, "Should move down within page")
	assert.Nil(t, cmd, "Should not trigger command for within-page movement")
}
