package mission

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Mission represents a parsed mission
type Mission struct {
	Type          string
	Track         string
	Iteration     string
	Status        string
	Intent        string
	Scope         []string
	Plan          []string
	Verification  string
	CompletedAt   *time.Time
	ParentMission string
	FilePath      string
}

// ReadCurrentMission reads the current mission from .mission/mission.md
func ReadCurrentMission() (*Mission, error) {
	return readMissionFile(".mission/mission.md")
}

// ReadCompletedMissions reads all completed missions from .mission/completed/
func ReadCompletedMissions() ([]*Mission, error) {
	completedDir := ".mission/completed"

	entries, err := os.ReadDir(completedDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read completed missions directory: %w", err)
	}

	var missions []*Mission
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), "-mission.md") {
			filePath := filepath.Join(completedDir, entry.Name())
			mission, err := readMissionFile(filePath)
			if err != nil {
				continue // Skip invalid files
			}
			missions = append(missions, mission)
		}
	}

	// Sort by completion time (newest first)
	sort.Slice(missions, func(i, j int) bool {
		if missions[i].CompletedAt == nil || missions[j].CompletedAt == nil {
			return false
		}
		return missions[i].CompletedAt.After(*missions[j].CompletedAt)
	})

	return missions, nil
}

// readMissionFile parses a mission.md file
func readMissionFile(filePath string) (*Mission, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open mission file: %w", err)
	}
	defer file.Close()

	mission := &Mission{FilePath: filePath}
	scanner := bufio.NewScanner(file)

	var currentSection string
	var planItems []string
	var scopeItems []string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "type:") {
			mission.Type = strings.TrimSpace(strings.TrimPrefix(line, "type:"))
		} else if strings.HasPrefix(line, "track:") {
			mission.Track = strings.TrimSpace(strings.TrimPrefix(line, "track:"))
		} else if strings.HasPrefix(line, "iteration:") {
			mission.Iteration = strings.TrimSpace(strings.TrimPrefix(line, "iteration:"))
		} else if strings.HasPrefix(line, "status:") {
			mission.Status = strings.TrimSpace(strings.TrimPrefix(line, "status:"))
		} else if strings.HasPrefix(line, "completed_at:") {
			timeStr := strings.TrimSpace(strings.TrimPrefix(line, "completed_at:"))
			if t, err := time.Parse(time.RFC3339, timeStr); err == nil {
				mission.CompletedAt = &t
			}
		} else if strings.HasPrefix(line, "parent_mission:") {
			mission.ParentMission = strings.TrimSpace(strings.TrimPrefix(line, "parent_mission:"))
		} else if line == "## INTENT" {
			currentSection = "intent"
		} else if line == "## SCOPE" {
			currentSection = "scope"
		} else if line == "## PLAN" {
			currentSection = "plan"
		} else if line == "## VERIFICATION" {
			currentSection = "verification"
		} else if strings.HasPrefix(line, "## ") {
			currentSection = ""
		} else if currentSection == "intent" && line != "" {
			mission.Intent = line
		} else if currentSection == "scope" && line != "" {
			scopeItems = append(scopeItems, line)
		} else if currentSection == "plan" && strings.HasPrefix(line, "- [") {
			planItems = append(planItems, line)
		} else if currentSection == "verification" && line != "" {
			mission.Verification = line
		}
	}

	mission.Plan = planItems
	mission.Scope = scopeItems

	return mission, scanner.Err()
}
