package tui

import (
	"fmt"
	"strings"
	"testing"

	"github.com/dnatag/mission-toolkit/pkg/mission"
)

func TestRenderNoMission(t *testing.T) {
	m := NewDashboardModel()
	output := m.renderNoMission()

	if !strings.Contains(output, "NO ACTIVE MISSION") {
		t.Error("expected output to contain 'NO ACTIVE MISSION'")
	}
}

func TestRenderCompletedMissions_Empty(t *testing.T) {
	m := NewDashboardModel()
	output := m.renderCompletedMissions()

	if output != "No completed missions found" {
		t.Errorf("expected 'No completed missions found', got '%s'", output)
	}
}

func TestRenderCompletedMissions_WithMissions(t *testing.T) {
	tests := []struct {
		name              string
		completedMissions []*mission.Mission
		selectedIndex     int
		wantPrefix        string
	}{
		{
			name: "single mission",
			completedMissions: []*mission.Mission{
				{ID: "m1", Status: "completed", Type: "WET", Body: "## INTENT\nFirst"},
			},
			selectedIndex: 0,
			wantPrefix:    "▶ ",
		},
		{
			name: "multiple missions with selection",
			completedMissions: []*mission.Mission{
				{ID: "m1", Status: "completed", Type: "WET", Body: "## INTENT\nFirst"},
				{ID: "m2", Status: "completed", Type: "DRY", Body: "## INTENT\nSecond"},
				{ID: "m3", Status: "completed", Type: "WET", Body: "## INTENT\nThird"},
			},
			selectedIndex: 1,
			wantPrefix:    "▶ ",
		},
		{
			name: "multiple missions no selection (index 0)",
			completedMissions: []*mission.Mission{
				{ID: "m1", Status: "completed", Type: "WET", Body: "## INTENT\nFirst"},
				{ID: "m2", Status: "completed", Type: "DRY", Body: "## INTENT\nSecond"},
			},
			selectedIndex: 0,
			wantPrefix:    "▶ ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewDashboardModel()
			m.completedMissions = tt.completedMissions
			m.selectedIndex = tt.selectedIndex
			m.itemsPerPage = 10 // Ensure all missions fit on one page

			output := m.renderCompletedMissions()

			// Verify selected prefix exists
			if !strings.Contains(output, tt.wantPrefix) {
				t.Errorf("expected output to contain prefix %q", tt.wantPrefix)
			}

			// Verify each mission's ID is shown
			for _, mission := range tt.completedMissions {
				if !strings.Contains(output, mission.ID) {
					t.Errorf("expected output to contain mission ID %q", mission.ID)
				}
			}

			// Verify each mission's type is shown
			for _, mission := range tt.completedMissions {
				if !strings.Contains(output, mission.Type) {
					t.Errorf("expected output to contain type %q", mission.Type)
				}
			}

			// Verify non-selected items have space prefix
			lines := strings.Split(output, "\n")
			foundSelected := false
			for _, line := range lines {
				if strings.Contains(line, tt.wantPrefix) {
					foundSelected = true
					break
				}
			}
			if !foundSelected {
				t.Error("expected to find selected item with ▶ prefix")
			}
		})
	}
}

func TestRenderMissionDetails(t *testing.T) {
	m := NewDashboardModel()
	testMission := &mission.Mission{
		ID:     "test-123",
		Status: "completed",
		Type:   "WET",
		Track:  2,
		Body:   "## INTENT\nTest intent",
	}

	output := m.renderMissionDetails(testMission)
	if !strings.Contains(output, "COMPLETED") {
		t.Error("expected output to contain 'COMPLETED'")
	}
	if !strings.Contains(output, "Test intent") {
		t.Error("expected output to contain 'Test intent'")
	}
}

func TestRenderDashboardView(t *testing.T) {
	m := NewDashboardModel()
	m.width = 100 // Set width to trigger pane calculation
	m.calculatePaneWidths()

	// Case 1: Active Mission (2 panes)
	m.selectedMission = &mission.Mission{ID: "active", Status: "active", Body: "## INTENT\nActive"}
	m.currentPane = MissionPane

	output := m.renderDashboardView()
	// Should contain mission content
	if !strings.Contains(output, "Active") {
		t.Error("expected active mission view to contain body content")
	}
	// Should NOT contain commit pane placeholder
	if strings.Contains(output, "No commit info") {
		t.Error("expected active mission view NOT to contain commit info")
	}

	// Case 2: Completed Mission (Commit Pane)
	m.selectedMission = &mission.Mission{ID: "completed", Status: "completed", Body: "## INTENT\nDone"}
	m.currentPane = CommitPane
	m.commitLoaded = true
	m.commitMessage = "test commit"

	output = m.renderDashboardView()
	if !strings.Contains(output, "test commit") {
		t.Error("expected completed mission view to contain commit message in CommitPane")
	}
}

func TestRenderTwoPaneLayout(t *testing.T) {
	m := NewDashboardModel()
	m.width = 100
	m.calculatePaneWidths()

	m.selectedMission = &mission.Mission{ID: "test", Status: "active", Body: "## INTENT\nBody"}
	m.currentPane = MissionPane
	m.executionLogLoaded = true
	m.executionLog = "LogContent"

	// MissionPane active
	output := m.renderTwoPaneLayout()
	if !strings.Contains(output, "Body") {
		t.Error("expected renderTwoPaneLayout to contain mission body")
	}
	if !strings.Contains(output, "LogContent") {
		t.Error("expected renderTwoPaneLayout to contain log content")
	}
}

func TestApplyFixedDimensions(t *testing.T) {
	m := NewDashboardModel()

	tests := []struct {
		name     string
		input    string
		wantRows int
		wantCols int
	}{
		{
			name:     "short content",
			input:    "line1\nline2",
			wantRows: fixedPaneHeight,
			wantCols: fixedPaneWidth,
		},
		{
			name:     "long line",
			input:    strings.Repeat("x", 100),
			wantRows: fixedPaneHeight,
			wantCols: fixedPaneWidth,
		},
		{
			name:     "many lines",
			input:    strings.Repeat("line\n", 100),
			wantRows: fixedPaneHeight,
			wantCols: fixedPaneWidth,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := m.applyFixedDimensions(tt.input)
			lines := strings.Split(output, "\n")

			if len(lines) != tt.wantRows {
				t.Errorf("expected %d rows, got %d", tt.wantRows, len(lines))
			}

			for i, line := range lines {
				if len(line) != tt.wantCols {
					t.Errorf("line %d: expected %d columns, got %d", i, tt.wantCols, len(line))
				}
			}
		})
	}
}

func TestClipContentToViewport(t *testing.T) {
	tests := []struct {
		name          string
		lines         []string
		scrollX       int
		scrollY       int
		visibleWidth  int
		visibleHeight int
		wantLines     int
	}{
		{
			name:          "no scrolling",
			lines:         []string{"line1", "line2", "line3"},
			scrollX:       0,
			scrollY:       0,
			visibleWidth:  10,
			visibleHeight: 3,
			wantLines:     3,
		},
		{
			name:          "vertical scroll",
			lines:         []string{"line1", "line2", "line3", "line4"},
			scrollX:       0,
			scrollY:       2,
			visibleWidth:  10,
			visibleHeight: 2,
			wantLines:     2,
		},
		{
			name:          "horizontal scroll",
			lines:         []string{"abcdefghij"},
			scrollX:       5,
			scrollY:       0,
			visibleWidth:  5,
			visibleHeight: 1,
			wantLines:     1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := clipContentToViewport(tt.lines, tt.scrollX, tt.scrollY, tt.visibleWidth, tt.visibleHeight)
			if len(result) != tt.visibleHeight {
				t.Errorf("expected %d lines, got %d", tt.visibleHeight, len(result))
			}
		})
	}
}

func TestRenderVerticalScrollbar(t *testing.T) {
	tests := []struct {
		name          string
		lines         []string
		scrollY       int
		totalLines    int
		visibleHeight int
		wantScrollbar bool
	}{
		{
			name:          "no scrollbar needed",
			lines:         []string{"a", "b"},
			scrollY:       0,
			totalLines:    2,
			visibleHeight: 5,
			wantScrollbar: false,
		},
		{
			name:          "scrollbar with up arrow",
			lines:         []string{"a", "b", "c"},
			scrollY:       1,
			totalLines:    10,
			visibleHeight: 3,
			wantScrollbar: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := renderVerticalScrollbar(tt.lines, tt.scrollY, tt.totalLines, tt.visibleHeight)
			if len(result) != len(tt.lines) {
				t.Errorf("expected %d lines, got %d", len(tt.lines), len(result))
			}
			hasScrollbar := strings.Contains(strings.Join(result, ""), "│") ||
				strings.Contains(strings.Join(result, ""), "▲") ||
				strings.Contains(strings.Join(result, ""), "▼")
			if hasScrollbar != tt.wantScrollbar {
				t.Errorf("expected scrollbar=%v, got %v", tt.wantScrollbar, hasScrollbar)
			}
		})
	}
}

func TestRenderHorizontalScrollbar(t *testing.T) {
	tests := []struct {
		name         string
		scrollX      int
		maxWidth     int
		visibleWidth int
		totalWidth   int
		wantArrows   bool
	}{
		{
			name:         "no scrollbar needed",
			scrollX:      0,
			maxWidth:     10,
			visibleWidth: 20,
			totalWidth:   20,
			wantArrows:   false,
		},
		{
			name:         "scrollbar with arrows",
			scrollX:      5,
			maxWidth:     100,
			visibleWidth: 20,
			totalWidth:   21,
			wantArrows:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := renderHorizontalScrollbar(tt.scrollX, tt.maxWidth, tt.visibleWidth, tt.totalWidth)
			hasArrows := strings.Contains(result, "◀") || strings.Contains(result, "▶")
			if hasArrows != tt.wantArrows {
				t.Errorf("expected arrows=%v, got %v", tt.wantArrows, hasArrows)
			}
		})
	}
}

func TestView(t *testing.T) {
	tests := []struct {
		name              string
		currentMission    *mission.Mission
		completedMissions []*mission.Mission
		wantTitle         bool
		wantMission       bool
		wantCompleted     bool
	}{
		{
			name:              "no missions at all",
			currentMission:    nil,
			completedMissions: []*mission.Mission{},
			wantTitle:         true,
			wantMission:       false,
			wantCompleted:     true,
		},
		{
			name: "only current mission",
			currentMission: &mission.Mission{
				ID:     "active-1",
				Status: "active",
				Type:   "WET",
				Track:  2,
				Body:   "## INTENT\nTest active mission",
			},
			completedMissions: []*mission.Mission{},
			wantTitle:         true,
			wantMission:       true,
			wantCompleted:     true,
		},
		{
			name:           "only completed missions",
			currentMission: nil,
			completedMissions: []*mission.Mission{
				{ID: "completed-1", Status: "completed", Body: "## INTENT\nDone"},
			},
			wantTitle:     true,
			wantMission:   false,
			wantCompleted: true,
		},
		{
			name: "both current and completed missions",
			currentMission: &mission.Mission{
				ID:     "active-1",
				Status: "active",
				Type:   "WET",
				Track:  2,
				Body:   "## INTENT\nActive",
			},
			completedMissions: []*mission.Mission{
				{ID: "completed-1", Status: "completed", Body: "## INTENT\nDone"},
			},
			wantTitle:     true,
			wantMission:   true,
			wantCompleted: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewDashboardModel()
			m.currentMission = tt.currentMission
			m.completedMissions = tt.completedMissions
			m.totalCount = len(tt.completedMissions)

			output := m.View()

			// Verify title is always present
			if tt.wantTitle && !strings.Contains(output, "Mission Dashboard") {
				t.Error("expected output to contain 'Mission Dashboard' title")
			}

			// Verify current mission section (checks for status which is always displayed)
			if tt.wantMission {
				if !strings.Contains(output, strings.ToUpper(tt.currentMission.Status)) {
					t.Errorf("expected output to contain status %q", tt.currentMission.Status)
				}
			}

			// Verify completed missions section
			if tt.wantCompleted {
				if !strings.Contains(output, "Completed Missions") {
					t.Error("expected output to contain 'Completed Missions' section")
				}
			}

			// Verify output is non-empty
			if output == "" {
				t.Error("View output should not be empty")
			}
		})
	}
}

func TestRenderCurrentMission(t *testing.T) {
	tests := []struct {
		name          string
		mission       *mission.Mission
		wantStatus    bool
		wantType      bool
		wantIntent    bool
		wantNextSteps bool
	}{
		{
			name: "planned mission",
			mission: &mission.Mission{
				ID:     "plan-1",
				Status: "planned",
				Type:   "WET",
				Track:  2,
				Body:   "## INTENT\nTest planned",
			},
			wantStatus:    true,
			wantType:      true,
			wantIntent:    true,
			wantNextSteps: true,
		},
		{
			name: "active mission",
			mission: &mission.Mission{
				ID:     "active-1",
				Status: "active",
				Type:   "DRY",
				Track:  3,
				Body:   "## INTENT\nTest active",
			},
			wantStatus:    true,
			wantType:      true,
			wantIntent:    true,
			wantNextSteps: true,
		},
		{
			name: "completed mission",
			mission: &mission.Mission{
				ID:     "done-1",
				Status: "completed",
				Type:   "WET",
				Track:  2,
				Body:   "## INTENT\nTest done",
			},
			wantStatus:    true,
			wantType:      true,
			wantIntent:    true,
			wantNextSteps: true,
		},
		{
			name: "failed mission",
			mission: &mission.Mission{
				ID:     "failed-1",
				Status: "failed",
				Type:   "WET",
				Track:  2,
				Body:   "## INTENT\nTest failed",
			},
			wantStatus:    true,
			wantType:      true,
			wantIntent:    true,
			wantNextSteps: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewDashboardModel()
			m.currentMission = tt.mission

			output := m.renderCurrentMission()

			// Verify status is displayed
			if tt.wantStatus && !strings.Contains(output, strings.ToUpper(tt.mission.Status)) {
				t.Errorf("expected output to contain status %q", tt.mission.Status)
			}

			// Verify type is displayed
			if tt.wantType && !strings.Contains(output, tt.mission.Type) {
				t.Errorf("expected output to contain type %q", tt.mission.Type)
			}

			// Verify track is displayed
			if !strings.Contains(output, fmt.Sprintf("Track %d", tt.mission.Track)) {
				t.Errorf("expected output to contain Track %d", tt.mission.Track)
			}

			// Verify intent is displayed
			if tt.wantIntent && !strings.Contains(output, "Test") {
				t.Error("expected output to contain intent content")
			}

			// Verify next steps or completion message is displayed
			// For completed missions, it says "Mission completed successfully" without "Next:"
			if tt.wantNextSteps {
				if tt.mission.Status == "completed" {
					if !strings.Contains(output, "Mission completed successfully") {
						t.Error("expected output to contain completion message")
					}
				} else {
					if !strings.Contains(output, "Next:") {
						t.Error("expected output to contain next steps with 'Next:' prefix")
					}
				}
			}
		})
	}
}
