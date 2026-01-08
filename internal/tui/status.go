package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dnatag/mission-toolkit/internal/mission"
	"github.com/spf13/afero"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	activeStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#04B575"))

	completedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262"))

	failedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF5F87"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262"))

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 2)
)

type Model struct {
	currentMission    *mission.Mission
	completedMissions []*mission.Mission
	selectedIndex     int
	selectedMission   *mission.Mission
	scrollOffset      int
	viewportHeight    int
	currentPage       int
	itemsPerPage      int
	searchMode        bool
	searchQuery       string
	filteredMissions  []*mission.Mission
	width             int
	height            int
	// Lazy loading state
	totalCount  int  // Total number of completed missions
	loadedCount int  // Number of missions currently loaded
	loading     bool // Whether a load operation is in progress
	loadError   error
	// Prefetch state
	prefetchedPages map[int][]*mission.Mission // Cache of prefetched pages (limited to adjacent pages)
	prefetching     bool                       // Whether a prefetch operation is in progress
	// Navigation state
	pendingPageChange int // Page to navigate to after loading completes (-1 = no pending change)
}

func NewModel() Model {
	return Model{
		selectedIndex:     0,
		currentPage:       0,
		itemsPerPage:      5,
		prefetchedPages:   make(map[int][]*mission.Mission),
		pendingPageChange: -1,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		loadCurrentMission,
		loadInitialMissions,
	)
}

func loadCurrentMission() tea.Msg {
	fs := afero.NewOsFs()
	reader := mission.NewReader(fs)
	m, err := reader.Read(".mission/mission.md")
	if err != nil {
		return currentMissionMsg{err: err}
	}
	return currentMissionMsg{mission: m}
}

// loadInitialMissions loads the first batch of completed missions for lazy loading
// Uses a batch size of 5 to balance performance with memory usage
func loadInitialMissions() tea.Msg {
	return loadCompletedMissionsBatch(0, 5)
}

// loadCompletedMissionsBatch loads a batch of missions from the filesystem
func loadCompletedMissionsBatch(offset, limit int) tea.Msg {
	fs := afero.NewOsFs()
	reader := mission.NewReader(fs)

	completedDir := ".mission/completed"
	entries, err := os.ReadDir(completedDir)
	if err != nil {
		return initialMissionsMsg{err: err}
	}

	// Collect and sort mission file paths by name (newest first)
	var missionFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), "-mission.md") {
			missionFiles = append(missionFiles, entry.Name())
		}
	}

	// Sort filenames in descending order (newest first)
	sort.Slice(missionFiles, func(i, j int) bool {
		return missionFiles[i] > missionFiles[j]
	})

	// Load missions starting from offset until we have enough or run out
	var missions []*mission.Mission
	loaded := 0
	fileIndex := 0

	// Skip files until we reach the offset
	for fileIndex < len(missionFiles) && loaded < offset {
		path := filepath.Join(completedDir, missionFiles[fileIndex])
		_, err := reader.Read(path)
		if err == nil {
			loaded++
		}
		fileIndex++
	}

	// Load the requested batch
	batchLoaded := 0
	for fileIndex < len(missionFiles) && batchLoaded < limit {
		path := filepath.Join(completedDir, missionFiles[fileIndex])
		m, err := reader.Read(path)
		if err == nil {
			missions = append(missions, m)
			batchLoaded++
		}
		fileIndex++
	}

	// Count total loadable missions
	totalLoadable := 0
	for _, filename := range missionFiles {
		path := filepath.Join(completedDir, filename)
		_, err := reader.Read(path)
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

// loadMoreMissions loads the next batch of missions
func loadMoreMissions(offset, limit int) tea.Msg {
	result := loadCompletedMissionsBatch(offset, limit)
	if msg, ok := result.(initialMissionsMsg); ok {
		return loadMoreMissionsMsg{
			missions:    msg.missions,
			loadedCount: msg.loadedCount,
		}
	}
	return loadMoreMissionsMsg{err: fmt.Errorf("failed to load more missions")}
}

// prefetchPage loads missions for a specific page in the background
func prefetchPage(page, itemsPerPage int) tea.Msg {
	offset := page * itemsPerPage
	result := loadCompletedMissionsBatch(offset, 5) // Use same batch size as current loading
	if msg, ok := result.(initialMissionsMsg); ok {
		return prefetchPageMsg{
			page:     page,
			missions: msg.missions,
		}
	}
	return prefetchPageMsg{
		page: page,
		err:  fmt.Errorf("failed to prefetch page %d", page),
	}
}

// triggerPrefetch starts prefetching for adjacent pages if not already cached
// This improves user experience by preloading content for likely navigation targets
func (m *Model) triggerPrefetch() tea.Cmd {
	if m.prefetching || m.searchMode {
		return nil
	}

	var cmds []tea.Cmd

	// Prefetch next page if not cached and within bounds
	nextPage := m.currentPage + 1
	if _, exists := m.prefetchedPages[nextPage]; !exists && nextPage < m.getTotalPages() {
		cmds = append(cmds, func() tea.Msg { return prefetchPage(nextPage, m.itemsPerPage) })
	}

	// Prefetch previous page if not cached and within bounds
	prevPage := m.currentPage - 1
	if _, exists := m.prefetchedPages[prevPage]; !exists && prevPage >= 0 {
		cmds = append(cmds, func() tea.Msg { return prefetchPage(prevPage, m.itemsPerPage) })
	}

	if len(cmds) > 0 {
		m.prefetching = true
		return tea.Batch(cmds...)
	}

	return nil
}

// triggerAggressivePrefetch prefetches more aggressively after loading new missions
// This ensures pages are ready even when total page count increases
func (m *Model) triggerAggressivePrefetch() tea.Cmd {
	if m.prefetching || m.searchMode {
		return nil
	}

	var cmds []tea.Cmd

	// Prefetch next 2 pages to account for newly available pages
	for i := 1; i <= 2; i++ {
		nextPage := m.currentPage + i
		if _, exists := m.prefetchedPages[nextPage]; !exists && nextPage < m.getTotalPages() {
			cmds = append(cmds, func() tea.Msg { return prefetchPage(nextPage, m.itemsPerPage) })
		}
	}

	// Prefetch previous page if not cached and within bounds
	prevPage := m.currentPage - 1
	if _, exists := m.prefetchedPages[prevPage]; !exists && prevPage >= 0 {
		cmds = append(cmds, func() tea.Msg { return prefetchPage(prevPage, m.itemsPerPage) })
	}

	if len(cmds) > 0 {
		m.prefetching = true
		return tea.Batch(cmds...)
	}

	return nil
}

// cleanupPrefetchCache removes cached pages that are no longer adjacent to current page
// This prevents unbounded memory growth while maintaining performance for likely navigation
func (m *Model) cleanupPrefetchCache() {
	for page := range m.prefetchedPages {
		// Keep pages within 2 positions of current page to account for aggressive prefetch
		if page < m.currentPage-1 || page > m.currentPage+2 {
			delete(m.prefetchedPages, page)
		}
	}
}

type currentMissionMsg struct {
	mission *mission.Mission
	err     error
}

type completedMissionsMsg struct {
	missions []*mission.Mission
	err      error
}

// Lazy loading message types
type initialMissionsMsg struct {
	missions    []*mission.Mission
	totalCount  int
	loadedCount int
	offset      int
	err         error
}

type loadMoreMissionsMsg struct {
	missions    []*mission.Mission
	loadedCount int
	err         error
}

// Prefetch message type
type prefetchPageMsg struct {
	page     int
	missions []*mission.Mission
	err      error
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Reserve space for title, separator, help text (approximately 10 lines)
		m.viewportHeight = max(5, msg.Height-10)
		return m, nil

	case currentMissionMsg:
		if msg.err == nil {
			m.currentMission = msg.mission
		}
		return m, nil

	case completedMissionsMsg:
		if msg.err == nil {
			m.completedMissions = msg.missions
		}
		return m, nil

	case initialMissionsMsg:
		if msg.err == nil {
			m.completedMissions = msg.missions
			m.totalCount = msg.totalCount
			m.loadedCount = msg.loadedCount
			m.loading = false
			m.loadError = nil

			// Trigger initial prefetch for adjacent pages
			return m, m.triggerPrefetch()
		} else {
			m.loadError = msg.err
		}
		return m, nil

	case loadMoreMissionsMsg:
		if msg.err == nil {
			m.completedMissions = append(m.completedMissions, msg.missions...)
			m.loadedCount = m.loadedCount + msg.loadedCount
			m.loading = false
			m.loadError = nil

			// Complete pending page change if any
			// This ensures navigation completes after lazy loading finishes
			if m.pendingPageChange >= 0 {
				m.currentPage = m.pendingPageChange
				m.selectedIndex = 0
				m.pendingPageChange = -1
			}

			// Trigger aggressive prefetch when new missions are loaded
			// This ensures next pages are ready even when total page count increases
			return m, m.triggerAggressivePrefetch()
		} else {
			m.loadError = msg.err
			// Clear pending page change on error
			m.pendingPageChange = -1
		}
		return m, nil

	case prefetchPageMsg:
		if msg.err == nil {
			// Store prefetched missions in cache
			m.prefetchedPages[msg.page] = msg.missions
			m.prefetching = false

			// Clean up cache to prevent memory growth - keep only adjacent pages
			m.cleanupPrefetchCache()
		}
		return m, nil

	case tea.KeyMsg:
		// Handle text input for search only when in search mode AND in mission list (not detail view)
		if m.searchMode && m.selectedMission == nil {
			key := msg.String()
			// Accept all printable ASCII characters (space through tilde)
			if len(key) == 1 && key[0] >= 32 && key[0] <= 126 {
				m.searchQuery += key
				m.filteredMissions = m.filterMissions(m.searchQuery)
				m.selectedIndex = 0
				m.currentPage = 0
				return m, nil
			}
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "r":
			if !m.searchMode && m.selectedMission == nil {
				// Reload missions with lazy loading
				return m, tea.Batch(loadCurrentMission, loadInitialMissions)
			}
		case "/":
			if m.selectedMission == nil {
				// Enter search mode
				m.searchMode = true
				m.searchQuery = ""
				m.selectedIndex = 0
				m.currentPage = 0
			}
		case "esc":
			if m.selectedMission != nil {
				m.selectedMission = nil
				m.scrollOffset = 0 // Reset scroll when exiting detail view
			} else if m.searchMode {
				// Exit search mode
				m.searchMode = false
				m.searchQuery = ""
				m.filteredMissions = nil
				m.selectedIndex = 0
				m.currentPage = 0
			}
		case "backspace":
			if m.searchMode && len(m.searchQuery) > 0 {
				m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
				m.filteredMissions = m.filterMissions(m.searchQuery)
				m.selectedIndex = 0
				m.currentPage = 0
			}
		case "up", "k":
			if m.selectedMission != nil {
				// Scroll up in detail view
				m.scrollOffset = max(0, m.scrollOffset-1)
			} else {
				missions := m.getActiveMissions()
				if len(missions) > 0 {
					if m.selectedIndex > 0 {
						// Move up within current page
						m.selectedIndex = m.selectedIndex - 1
					} else if !m.searchMode && m.currentPage > 0 {
						// At top of page, go to previous page (bottom item)
						m.currentPage = m.currentPage - 1
						pageSize := m.getPageSize()
						m.selectedIndex = pageSize - 1

						// Trigger prefetch for adjacent pages
						return m, m.triggerPrefetch()
					}
				}
			}
		case "down", "j":
			if m.selectedMission != nil {
				// Scroll down in detail view with bounds checking
				maxScroll := m.getMaxScrollOffset()
				m.scrollOffset = min(maxScroll, m.scrollOffset+1)
			} else {
				missions := m.getActiveMissions()
				if len(missions) > 0 {
					pageSize := m.getPageSize()

					if m.selectedIndex < pageSize-1 {
						// Check if moving down would trigger lazy loading of more missions
						// This must be done before updating selectedIndex to avoid off-by-one errors
						if !m.searchMode && !m.loading && m.loadedCount < m.totalCount {
							// Calculate the absolute position after the proposed move
							newSelectedIndex := m.selectedIndex + 1
							absoluteIndex := newSelectedIndex + (m.currentPage * m.itemsPerPage)

							// Trigger loading when we reach the last currently loaded mission
							if absoluteIndex >= m.loadedCount-1 {
								m.loading = true
								m.selectedIndex = newSelectedIndex
								return m, func() tea.Msg { return loadMoreMissions(m.loadedCount, 5) }
							}
						}

						// Move down within current page
						m.selectedIndex = m.selectedIndex + 1
					} else if !m.searchMode {
						// At bottom of page, go to next page (top item)
						totalPages := m.getTotalPages()
						newPage := m.currentPage + 1

						if newPage < totalPages {
							// Check if moving to next page would trigger lazy loading
							if !m.loading && m.loadedCount < m.totalCount {
								// Calculate if the new page would need more missions
								firstItemOnNewPage := newPage * m.itemsPerPage
								if firstItemOnNewPage >= m.loadedCount {
									// Store pending page change and load more missions
									// This ensures page navigation completes after loading finishes
									m.pendingPageChange = newPage
									m.loading = true
									return m, func() tea.Msg { return loadMoreMissions(m.loadedCount, 5) }
								}
							}

							// No lazy loading needed, change page immediately
							m.currentPage = newPage
							m.selectedIndex = 0

							// Trigger prefetch for adjacent pages
							return m, m.triggerPrefetch()
						}
					}
				}
			}
		case "left", "h":
			if m.selectedMission == nil && !m.searchMode {
				missions := m.getActiveMissions()
				if len(missions) > 0 {
					// Previous page
					m.currentPage = max(0, m.currentPage-1)
					m.selectedIndex = 0

					// Trigger prefetch for adjacent pages
					return m, m.triggerPrefetch()
				}
			}
		case "right", "l":
			if m.selectedMission == nil && !m.searchMode {
				missions := m.getActiveMissions()
				if len(missions) > 0 {
					// Next page
					totalPages := m.getTotalPages()
					newPage := min(totalPages-1, m.currentPage+1)

					// Check if moving to next page would trigger lazy loading
					if !m.loading && m.loadedCount < m.totalCount {
						// Calculate if the new page would need more missions
						firstItemOnNewPage := newPage * m.itemsPerPage
						if firstItemOnNewPage >= m.loadedCount {
							// Store pending page change and load more missions
							// This ensures page navigation completes after loading finishes
							m.pendingPageChange = newPage
							m.loading = true
							return m, func() tea.Msg { return loadMoreMissions(m.loadedCount, 5) }
						}
					}

					// No lazy loading needed, change page immediately
					m.currentPage = newPage
					m.selectedIndex = 0

					// Trigger prefetch for adjacent pages
					return m, m.triggerPrefetch()
				}
			}
		case "pgup":
			if m.selectedMission != nil {
				// Scroll up in mission detail view
				m.scrollOffset = max(0, m.scrollOffset-5)
			} else if !m.searchMode {
				// Navigate to previous page in mission list
				missions := m.getActiveMissions()
				if len(missions) > 0 {
					m.currentPage = max(0, m.currentPage-1)
					m.selectedIndex = 0

					// Trigger prefetch for adjacent pages
					return m, m.triggerPrefetch()
				}
			}
		case "pgdn":
			if m.selectedMission != nil {
				// Scroll down in mission detail view
				maxScroll := m.getMaxScrollOffset()
				m.scrollOffset = min(maxScroll, m.scrollOffset+5)
			} else if !m.searchMode {
				// Navigate to next page in mission list with lazy loading support
				missions := m.getActiveMissions()
				if len(missions) > 0 {
					totalPages := m.getTotalPages()
					newPage := min(totalPages-1, m.currentPage+1)

					// Check if moving to next page would trigger lazy loading
					if !m.loading && m.loadedCount < m.totalCount {
						// Calculate if the new page would need more missions
						firstItemOnNewPage := newPage * m.itemsPerPage
						if firstItemOnNewPage >= m.loadedCount {
							// Store pending page change and load more missions
							// This ensures page navigation completes after loading finishes
							m.pendingPageChange = newPage
							m.loading = true
							return m, func() tea.Msg { return loadMoreMissions(m.loadedCount, 5) }
						}
					}

					// No lazy loading needed, change page immediately
					m.currentPage = newPage
					m.selectedIndex = 0

					// Trigger prefetch for adjacent pages
					return m, m.triggerPrefetch()
				}
			}
		case "enter":
			if m.selectedMission == nil {
				pageMissions := m.getCurrentPageMissions()
				if len(pageMissions) > 0 && m.selectedIndex < len(pageMissions) {
					m.selectedMission = pageMissions[m.selectedIndex]
					m.scrollOffset = 0 // Reset scroll when entering detail view
				}
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	var sections []string

	// Title
	sections = append(sections, titleStyle.Render("Mission Toolkit Status"))
	sections = append(sections, "")

	// Current Mission (Top Section)
	if m.currentMission != nil {
		sections = append(sections, m.renderCurrentMission())
	} else {
		sections = append(sections, m.renderNoMission())
	}

	sections = append(sections, "")
	sections = append(sections, strings.Repeat("─", 60))
	sections = append(sections, "")

	// Completed Missions (Bottom Section)
	missions := m.getActiveMissions()
	var headerText string
	if m.searchMode {
		if m.searchQuery == "" {
			headerText = fmt.Sprintf("Search: _ (%d total missions)", m.totalCount)
		} else {
			headerText = fmt.Sprintf("Search: %s (%d of %d missions)", m.searchQuery, len(missions), m.totalCount)
		}
	} else {
		headerText = fmt.Sprintf("Completed Missions (%d)", m.totalCount)
	}
	sections = append(sections, headerText)
	sections = append(sections, "")
	sections = append(sections, m.renderCompletedMissions())

	// Help
	sections = append(sections, "")
	if m.selectedMission != nil {
		sections = append(sections, helpStyle.Render("↑/↓: scroll • PgUp/PgDn: fast scroll • Esc: back to list • q: quit"))
	} else if m.searchMode {
		sections = append(sections, helpStyle.Render("Type to search • Backspace: delete • Esc: exit search • Enter: view details"))
	} else {
		sections = append(sections, helpStyle.Render("↑/↓: navigate • ←/→: page • /: search • r: reload • Enter: view details • q: quit"))
	}

	return strings.Join(sections, "\n")
}

// getActiveMissions returns either filtered missions or all completed missions
func (m Model) getActiveMissions() []*mission.Mission {
	if m.searchMode {
		if m.searchQuery == "" {
			return []*mission.Mission{} // Show empty list when search is active but no query
		}
		return m.filteredMissions
	}
	return m.completedMissions
}

// filterMissions performs fuzzy search on missions
func (m Model) filterMissions(query string) []*mission.Mission {
	if query == "" {
		return []*mission.Mission{} // Return empty slice for empty query, not all missions
	}

	query = strings.ToLower(query)
	var filtered []*mission.Mission

	for _, mission := range m.completedMissions {
		if m.matchesFuzzy(mission, query) {
			filtered = append(filtered, mission)
		}
	}

	return filtered
}

// matchesFuzzy performs fuzzy matching against mission content
func (m Model) matchesFuzzy(mission *mission.Mission, query string) bool {
	// Check intent
	intent := extractIntent(mission.Body)
	if strings.Contains(strings.ToLower(intent), query) {
		return true
	}

	// Check status
	if strings.Contains(strings.ToLower(mission.Status), query) {
		return true
	}

	// Check type
	if strings.Contains(strings.ToLower(mission.Type), query) {
		return true
	}

	// Check track
	trackStr := fmt.Sprintf("%d", mission.Track)
	if strings.Contains(trackStr, query) {
		return true
	}

	return false
}

// getTotalPages calculates the total number of pages
func (m Model) getTotalPages() int {
	// In search mode, use actual filtered count
	if m.searchMode {
		missions := m.getActiveMissions()
		if len(missions) == 0 {
			return 1
		}
		return (len(missions) + m.itemsPerPage - 1) / m.itemsPerPage
	}
	// In normal mode, use totalCount for lazy loading pagination
	if m.totalCount == 0 {
		return 1
	}
	return (m.totalCount + m.itemsPerPage - 1) / m.itemsPerPage
}

// getPageSize returns the number of items on the current page
func (m Model) getPageSize() int {
	missions := m.getActiveMissions()
	totalItems := len(missions)
	if totalItems == 0 {
		return 0
	}

	start := m.currentPage * m.itemsPerPage
	if start >= totalItems {
		return 0
	}

	end := min(totalItems, start+m.itemsPerPage)
	return end - start
}

// getCurrentPageMissions returns the missions for the current page
func (m Model) getCurrentPageMissions() []*mission.Mission {
	missions := m.getActiveMissions()
	totalItems := len(missions)
	if totalItems == 0 {
		return nil
	}

	start := m.currentPage * m.itemsPerPage
	if start >= totalItems {
		return nil
	}

	end := min(totalItems, start+m.itemsPerPage)
	return missions[start:end]
}

// getMaxScrollOffset calculates the maximum scroll offset for the current mission details
func (m Model) getMaxScrollOffset() int {
	if m.selectedMission == nil || m.viewportHeight <= 0 {
		return 0
	}

	// Calculate total lines in mission details
	totalLines := 0
	totalLines += 3 // Status, completed time, empty line
	totalLines += 1 // Intent line

	scope := extractScope(m.selectedMission.Body)
	if len(scope) > 0 {
		totalLines += 2 + len(scope) // "Scope:" + scope items
	}

	plan := extractPlan(m.selectedMission.Body)
	if len(plan) > 0 {
		totalLines += 2 + len(plan) // "Plan:" + plan items
	}

	verification := extractVerification(m.selectedMission.Body)
	if verification != "" {
		totalLines += 2 // empty line + verification
	}

	return max(0, totalLines-m.viewportHeight)
}

func (m Model) renderMissionDetails(mission *mission.Mission) string {
	var statusStyle lipgloss.Style
	switch mission.Status {
	case "completed":
		statusStyle = completedStyle
	case "failed":
		statusStyle = failedStyle
	default:
		statusStyle = lipgloss.NewStyle()
	}

	intent := extractIntent(mission.Body)
	scope := extractScope(mission.Body)
	plan := extractPlan(mission.Body)
	verification := extractVerification(mission.Body)

	var sections []string
	sections = append(sections, fmt.Sprintf("%s %s (Track %d)",
		statusStyle.Render(strings.ToUpper(mission.Status)),
		mission.Type,
		mission.Track))
	sections = append(sections, "")
	sections = append(sections, fmt.Sprintf("Intent: %s", intent))

	if len(scope) > 0 {
		sections = append(sections, "")
		sections = append(sections, "Scope:")
		for _, s := range scope {
			sections = append(sections, fmt.Sprintf("  %s", s))
		}
	}

	if len(plan) > 0 {
		sections = append(sections, "")
		sections = append(sections, "Plan:")
		for _, p := range plan {
			sections = append(sections, fmt.Sprintf("  %s", p))
		}
	}

	if verification != "" {
		sections = append(sections, "")
		sections = append(sections, fmt.Sprintf("Verification: %s", verification))
	}

	// Apply viewport scrolling
	totalLines := len(sections)
	if m.viewportHeight > 0 && totalLines > m.viewportHeight {
		// Ensure scroll offset doesn't exceed content
		maxScroll := max(0, totalLines-m.viewportHeight)
		scrollOffset := min(m.scrollOffset, maxScroll)

		// Get visible lines
		start := scrollOffset
		end := min(totalLines, scrollOffset+m.viewportHeight)
		visibleSections := sections[start:end]

		// Add scroll indicators
		var result []string
		if scrollOffset > 0 {
			result = append(result, helpStyle.Render("↑ more above"))
		}
		result = append(result, visibleSections...)
		if end < totalLines {
			result = append(result, helpStyle.Render("↓ more below"))
		}

		return boxStyle.Render(strings.Join(result, "\n"))
	}

	return boxStyle.Render(strings.Join(sections, "\n"))
}

func (m Model) renderCurrentMission() string {
	mission := m.currentMission

	var statusStyle lipgloss.Style
	var nextSteps string

	switch mission.Status {
	case "planned":
		statusStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00CED1"))
		nextSteps = "Next: Run '/m.apply' to execute the mission"
	case "active":
		statusStyle = activeStyle
		nextSteps = "Next: Run '/m.complete' to finalize the mission"
	case "completed":
		statusStyle = completedStyle
		nextSteps = "Mission completed successfully"
	case "failed":
		statusStyle = failedStyle
		nextSteps = "Next: Create a new mission with smaller scope using '/m.plan'"
	default:
		statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
		nextSteps = "Unknown status - use '/m.plan' to create a new mission"
	}

	intent := extractIntent(mission.Body)
	content := fmt.Sprintf("%s %s (Track %d)\n\n%s\n\n%s",
		statusStyle.Render(strings.ToUpper(mission.Status)),
		mission.Type,
		mission.Track,
		intent,
		nextSteps,
	)

	return boxStyle.Render(content)
}

func (m Model) renderNoMission() string {
	noMissionStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#626262"))

	content := fmt.Sprintf("%s\n\n%s",
		noMissionStyle.Render("NO ACTIVE MISSION"),
		"Use 'm.plan' to start with your intent",
	)

	return boxStyle.Render(content)
}

func (m Model) renderCompletedMissions() string {
	if !m.searchMode && len(m.completedMissions) == 0 {
		return "No completed missions found"
	}

	// Show detailed view if a mission is selected
	if m.selectedMission != nil {
		return m.renderMissionDetails(m.selectedMission)
	}

	// Handle search mode
	if m.searchMode {
		if m.searchQuery == "" {
			return "Type to search missions..."
		}

		missions := m.filteredMissions
		if len(missions) == 0 {
			return "No missions match your search"
		}

		// Show search results with pagination
		pageMissions := m.getCurrentPageMissions()
		if len(pageMissions) == 0 {
			return "No missions on this page"
		}

		var items []string
		for i, mission := range pageMissions {
			prefix := "  "
			if i == m.selectedIndex {
				prefix = "▶ "
			}

			intent := extractIntent(mission.Body)
			item := fmt.Sprintf("%s%s [%s] %s",
				prefix,
				mission.ID,
				mission.Type,
				truncate(intent, 50),
			)
			items = append(items, item)
		}

		// Add page indicator if multiple pages
		if m.getTotalPages() > 1 {
			totalPages := m.getTotalPages()
			pageIndicator := fmt.Sprintf("\nPage %d of %d", m.currentPage+1, totalPages)
			items = append(items, helpStyle.Render(pageIndicator))
		}

		return strings.Join(items, "\n")
	}

	// Normal mode - show all missions with pagination
	pageMissions := m.getCurrentPageMissions()
	if len(pageMissions) == 0 {
		return "No missions on this page"
	}

	var items []string
	for i, mission := range pageMissions {
		prefix := "  "
		if i == m.selectedIndex {
			prefix = "▶ "
		}

		intent := extractIntent(mission.Body)
		item := fmt.Sprintf("%s%s [%s] %s",
			prefix,
			mission.ID,
			mission.Type,
			truncate(intent, 50),
		)
		items = append(items, item)
	}

	// Add page indicator
	totalPages := m.getTotalPages()
	pageIndicator := fmt.Sprintf("\nPage %d of %d", m.currentPage+1, totalPages)
	items = append(items, helpStyle.Render(pageIndicator))

	// Add loading indicator if more missions are being loaded
	if m.loading {
		items = append(items, helpStyle.Render("Loading more missions..."))
	} else if m.loadedCount < m.totalCount {
		items = append(items, helpStyle.Render("Scroll down to load more..."))
	}

	return strings.Join(items, "\n")
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Helper functions to extract sections from mission body
func extractIntent(body string) string {
	re := regexp.MustCompile(`(?s)## INTENT\s*\n(.*?)(?:\n##|$)`)
	matches := re.FindStringSubmatch(body)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}

func extractScope(body string) []string {
	re := regexp.MustCompile(`(?s)## SCOPE\s*\n(.*?)(?:\n##|$)`)
	matches := re.FindStringSubmatch(body)
	if len(matches) > 1 {
		lines := strings.Split(strings.TrimSpace(matches[1]), "\n")
		var scope []string
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" {
				scope = append(scope, line)
			}
		}
		return scope
	}
	return nil
}

func extractPlan(body string) []string {
	re := regexp.MustCompile(`(?s)## PLAN\s*\n(.*?)(?:\n##|$)`)
	matches := re.FindStringSubmatch(body)
	if len(matches) > 1 {
		lines := strings.Split(strings.TrimSpace(matches[1]), "\n")
		var plan []string
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" {
				plan = append(plan, line)
			}
		}
		return plan
	}
	return nil
}

func extractVerification(body string) string {
	re := regexp.MustCompile(`(?s)## VERIFICATION\s*\n(.*?)(?:\n##|$)`)
	matches := re.FindStringSubmatch(body)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}

// RunStatusTUI starts the TUI for status display
func RunStatusTUI() error {
	p := tea.NewProgram(NewModel(), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
