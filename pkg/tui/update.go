package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/dnatag/mission-toolkit/pkg/git"
	"github.com/dnatag/mission-toolkit/pkg/mission"
	"github.com/spf13/afero"
)

// Update handles all state updates for the dashboard
func (m DashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewportHeight = max(5, msg.Height-10)
		m.calculatePaneWidths()
		return m, nil

	case currentMissionMsg:
		if msg.err == nil && msg.mission != nil {
			m.currentMission = msg.mission
			// Start refresh ticker for active missions
			if msg.mission.Status == "active" {
				return m, m.startRefreshTicker()
			}
		}
		return m, nil

	case initialMissionsMsg:
		if msg.err == nil {
			m.completedMissions = msg.missions
			m.totalCount = msg.totalCount
			m.loadedCount = msg.loadedCount
			m.loading = false
			m.loadError = nil
		} else {
			m.loadError = msg.err
		}
		return m, nil

	case refreshTickMsg:
		// Refresh execution log for active missions
		if m.currentMission != nil && m.currentMission.Status == "active" {
			return m, tea.Batch(
				loadExecutionLog(m.currentMission.ID, true),
				tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
					return refreshTickMsg{time: t}
				}),
			)
		}
		return m, nil

	case executionLogMsg:
		if msg.err == nil {
			m.executionLog = msg.content
			m.executionLogLoaded = true
		}
		return m, nil

	case commitMsg:
		if msg.err == nil {
			m.commitMessage = msg.content
			m.commitLoaded = true
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			if m.selectedMission != nil {
				// Cycle through panes
				// MissionPane -> ExecutionLogPane -> CommitPane (if completed) -> MissionPane
				nextPane := m.currentPane

				if m.selectedMission.Status == "completed" {
					switch m.currentPane {
					case MissionPane:
						nextPane = ExecutionLogPane
					case ExecutionLogPane:
						nextPane = CommitPane
					case CommitPane:
						nextPane = MissionPane
					}
				} else {
					// For active missions, toggle between Mission and Log
					switch m.currentPane {
					case MissionPane:
						nextPane = ExecutionLogPane
					case ExecutionLogPane:
						nextPane = MissionPane
					default:
						nextPane = MissionPane
					}
				}

				m.currentPane = nextPane

				// Load content when switching to panes
				var cmd tea.Cmd
				if m.currentPane == ExecutionLogPane && !m.executionLogLoaded {
					cmd = loadExecutionLog(m.selectedMission.ID, m.selectedMission.Status == "active")
				} else if m.currentPane == CommitPane && !m.commitLoaded && m.selectedMission.Status == "completed" {
					cmd = loadCommitMessage(m.selectedMission.ID)
				}
				return m, cmd
			}
		case "esc":
			if m.selectedMission != nil {
				m.selectedMission = nil
				m.currentPane = MissionPane
				m.executionLogLoaded = false
				m.commitLoaded = false
				m.scrollOffset = 0
				// Reset scroll positions
				m.leftPaneScrollX = 0
				m.leftPaneScrollY = 0
				m.rightPaneScrollX = 0
				m.rightPaneScrollY = 0
			}
		case "enter":
			if m.selectedMission == nil {
				pageMissions := m.getCurrentPageMissions()
				if len(pageMissions) > 0 && m.selectedIndex < len(pageMissions) {
					m.selectedMission = pageMissions[m.selectedIndex]
					m.scrollOffset = 0
					m.currentPane = MissionPane
					// Reset scroll positions for new mission
					m.leftPaneScrollX = 0
					m.leftPaneScrollY = 0
					m.rightPaneScrollX = 0
					m.rightPaneScrollY = 0
				}
			}
		case "up", "k":
			if m.selectedMission == nil && len(m.completedMissions) > 0 {
				if m.selectedIndex > 0 {
					m.selectedIndex--
				} else if m.currentPage > 0 {
					// At top of page, go to previous page
					m.currentPage--
					pageMissions := m.getCurrentPageMissions()
					m.selectedIndex = len(pageMissions) - 1
				}
			} else if m.selectedMission != nil {
				// Vertical scroll up in active pane
				if m.currentPane == MissionPane {
					if m.leftPaneScrollY > 0 {
						m.leftPaneScrollY--
					}
				} else {
					if m.rightPaneScrollY > 0 {
						m.rightPaneScrollY--
					}
				}
			}
		case "down", "j":
			if m.selectedMission == nil {
				pageMissions := m.getCurrentPageMissions()
				if len(pageMissions) > 0 {
					if m.selectedIndex < len(pageMissions)-1 {
						m.selectedIndex++
					} else {
						// At bottom of page, go to next page
						totalPages := (len(m.completedMissions) + m.itemsPerPage - 1) / m.itemsPerPage
						if m.currentPage < totalPages-1 {
							m.currentPage++
							m.selectedIndex = 0
						}
					}
				}
			} else if m.selectedMission != nil {
				// Vertical scroll down in active pane - always allow scrolling down
				if m.currentPane == MissionPane {
					m.leftPaneScrollY++
				} else {
					m.rightPaneScrollY++
				}
			}
		case "left", "h":
			if m.selectedMission == nil && m.currentPage > 0 {
				m.currentPage--
				m.selectedIndex = 0
			} else if m.selectedMission != nil {
				// Horizontal scroll left in active pane
				if m.currentPane == MissionPane {
					if m.leftPaneScrollX > 0 {
						m.leftPaneScrollX--
					}
				} else {
					if m.rightPaneScrollX > 0 {
						m.rightPaneScrollX--
					}
				}
			}
		case "right", "l":
			if m.selectedMission == nil && len(m.completedMissions) > 0 {
				totalPages := (len(m.completedMissions) + m.itemsPerPage - 1) / m.itemsPerPage
				if m.currentPage < totalPages-1 {
					m.currentPage++
					m.selectedIndex = 0
				}
			} else if m.selectedMission != nil {
				// Horizontal scroll right in active pane - always allow scrolling right
				if m.currentPane == MissionPane {
					m.leftPaneScrollX++
				} else {
					m.rightPaneScrollX++
				}
			}
		}
	}

	return m, nil
}

// loadExecutionLog loads the execution log for the current or selected mission
func loadExecutionLog(missionID string, isActive bool) tea.Cmd {
	return func() tea.Msg {
		var logPath string
		if isActive {
			logPath = ".mission/execution.log"
		} else {
			logPath = fmt.Sprintf(".mission/completed/%s-execution.log", missionID)
		}

		content, err := os.ReadFile(logPath)
		if err != nil {
			return executionLogMsg{content: "No execution log available"}
		}

		return executionLogMsg{content: string(content)}
	}
}

// loadCommitMessage loads the commit message for a completed mission
func loadCommitMessage(missionID string) tea.Cmd {
	return func() tea.Msg {
		// For completed missions, try to read the archived commit message first
		commitPath := fmt.Sprintf(".mission/completed/%s-commit.msg", missionID)
		if content, err := os.ReadFile(commitPath); err == nil {
			return commitMsg{content: string(content)}
		}

		// Fallback: try to get from git using the final consolidated commit
		// The final commit should be tagged when the mission is completed
		gitClient := git.NewCmdGitClient(".")

		// Look for the mission's final commit by checking recent commits
		// that contain the mission ID in the commit message
		commitMessage, err := gitClient.GetCommitMessage("HEAD")
		if err == nil && strings.Contains(commitMessage, missionID) {
			return commitMsg{content: commitMessage}
		}

		// If that fails, provide a helpful message
		content := fmt.Sprintf("Mission %s\n\nCommit message not found\nThe mission may not have been completed with @m.complete",
			missionID)
		return commitMsg{content: content}
	}
}

// loadCurrentMission loads the current active mission
func loadCurrentMission() tea.Msg {
	fs := afero.NewOsFs()
	missionPath := ".mission/mission.md"
	reader := mission.NewReader(fs, missionPath)
	m, err := reader.Read()
	if err != nil {
		return currentMissionMsg{err: err}
	}
	return currentMissionMsg{mission: m}
}

// loadInitialMissions loads the first batch of completed missions
func loadInitialMissions() tea.Msg {
	return loadCompletedMissionsBatch(0, -1) // Load all missions
}

// loadCompletedMissionsBatch loads a batch of completed missions
func loadCompletedMissionsBatch(offset, limit int) tea.Msg {
	fs := afero.NewOsFs()
	completedDir := ".mission/completed"
	entries, err := os.ReadDir(completedDir)
	if err != nil {
		return initialMissionsMsg{err: err}
	}

	var missionFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), "-mission.md") {
			missionFiles = append(missionFiles, entry.Name())
		}
	}

	sort.Slice(missionFiles, func(i, j int) bool {
		return missionFiles[i] > missionFiles[j]
	})

	var missions []*mission.Mission
	loaded := 0
	fileIndex := 0

	for fileIndex < len(missionFiles) && loaded < offset {
		path := filepath.Join(completedDir, missionFiles[fileIndex])
		reader := mission.NewReader(fs, path)
		_, err := reader.Read()
		if err == nil {
			loaded++
		}
		fileIndex++
	}

	batchLoaded := 0
	for fileIndex < len(missionFiles) && (limit < 0 || batchLoaded < limit) {
		path := filepath.Join(completedDir, missionFiles[fileIndex])
		reader := mission.NewReader(fs, path)
		m, err := reader.Read()
		if err == nil {
			missions = append(missions, m)
			batchLoaded++
		}
		fileIndex++
	}

	totalLoadable := 0
	for _, filename := range missionFiles {
		path := filepath.Join(completedDir, filename)
		reader := mission.NewReader(fs, path)
		_, err := reader.Read()
		if err == nil {
			totalLoadable++
		}
	}

	return initialMissionsMsg{
		missions:    missions,
		totalCount:  totalLoadable,
		loadedCount: len(missions),
		offset:      offset,
	}
}
