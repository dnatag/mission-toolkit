package beads

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	backlog "github.com/dnatag/mission-toolkit/pkg/backlog"
	"github.com/dnatag/mission-toolkit/pkg/backlog/file"
)

// Provider implements backlog.BacklogProvider using the bd CLI.
type Provider struct {
	projectRoot   string
	commandRunner CommandRunner
	epicCache     *EpicCache
	cachePath     string
	hashPath      string
	reconciling   bool
}

// NewProvider creates a new beads Provider for the given project directory.
func NewProvider(projectRoot string) *Provider {
	missionDir := filepath.Join(projectRoot, ".mission")
	return &Provider{
		projectRoot:   projectRoot,
		commandRunner: NewBDCommandRunner(projectRoot),
		epicCache:     &EpicCache{Epics: make(map[string]string)},
		cachePath:     filepath.Join(missionDir, "beads-epics.json"),
		hashPath:      filepath.Join(missionDir, "beads-backlog-hash"),
	}
}

func (p *Provider) reconcile() {
	if p.reconciling {
		return
	}
	p.reconciling = true
	defer func() { p.reconciling = false }()

	backlogPath := filepath.Join(p.projectRoot, ".mission", "backlog.md")

	data, err := os.ReadFile(backlogPath)
	if err != nil {
		return
	}

	currentHash := fmt.Sprintf("%x", sha256.Sum256(data))
	storedHash, err := os.ReadFile(p.hashPath)
	if err == nil && string(storedHash) == currentHash {
		return
	}

	fileProvider := file.NewManager(filepath.Join(p.projectRoot, ".mission"))
	reconciler := backlog.NewReconciler(fileProvider, p)
	edits, err := reconciler.DetectRogueEdits()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: reconciliation detection failed: %v\n", err)
		return
	}

	if len(edits) > 0 {
		imported, err := reconciler.ImportToBeads(edits)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: reconciliation import failed: %v\n", err)
		}
		if imported > 0 {
			fmt.Fprintf(os.Stderr, "Reconciled %d rogue item(s) from backlog.md into Beads\n", imported)
		}
	}
}

// List returns backlog items from Beads.
func (p *Provider) List(include []string, exclude []string) ([]string, error) {
	p.reconcile()
	if err := p.ensureEpics(); err != nil {
		return nil, fmt.Errorf("failed to ensure epics: %w", err)
	}

	includeCompleted := backlog.ContainsString(include, "completed")
	excludeCompleted := backlog.ContainsString(exclude, "completed")
	typeInclude := backlog.RemoveString(include, "completed")
	typeExclude := backlog.RemoveString(exclude, "completed")

	epicTypes := []string{
		backlog.ItemTypeFeature,
		backlog.ItemTypeBugfix,
		backlog.ItemTypeDecomposed,
		backlog.ItemTypeRefactor,
		backlog.ItemTypeFuture,
	}

	queryTypes := backlog.FilterStringSlice(epicTypes, typeInclude, typeExclude)

	var results []string
	for _, itemType := range queryTypes {
		epicID, ok := p.epicCache.Epics[itemType]
		if !ok {
			continue
		}

		output, err := p.commandRunner.Run("list", "--parent", epicID, "--json")
		if err != nil {
			return nil, fmt.Errorf("failed to list %s items: %w", itemType, err)
		}

		items, err := p.parseListOutput(output)
		if err != nil {
			return nil, fmt.Errorf("failed to parse list output for %s: %w", itemType, err)
		}

		for _, item := range items {
			isClosed := strings.HasPrefix(item, backlog.ListItemClosedFormat)
			if includeCompleted && !isClosed {
				continue
			}
			if excludeCompleted && isClosed {
				continue
			}
			results = append(results, item)
		}
	}

	return results, nil
}

func (p *Provider) parseListOutput(output string) ([]string, error) {
	var items []struct {
		ID     string `json:"id"`
		Title  string `json:"title"`
		Status string `json:"status"`
	}

	if err := json.Unmarshal([]byte(output), &items); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	results := make([]string, 0, len(items))
	for _, item := range items {
		var prefix string
		switch item.Status {
		case "closed":
			prefix = backlog.ListItemClosedFormat
		default:
			prefix = backlog.ListItemOpenFormat
		}
		results = append(results, prefix+item.Title)
	}

	return results, nil
}

// Add adds a new item to the specified section.
func (p *Provider) Add(description, itemType string) error {
	p.reconcile()
	if err := p.addItem(description, itemType); err != nil {
		return err
	}
	p.GenerateSnapshot()
	return nil
}

func (p *Provider) addItem(description, itemType string) error {
	if err := p.ensureEpics(); err != nil {
		return fmt.Errorf("failed to ensure epics: %w", err)
	}

	if !backlog.IsValidType(itemType) {
		return fmt.Errorf("invalid type: %s. Valid types: feature, bugfix, decomposed, refactor, future", itemType)
	}

	epicID, ok := p.epicCache.Epics[itemType]
	if !ok {
		return fmt.Errorf("epic not found for type: %s", itemType)
	}

	output, err := p.commandRunner.Run("create", description, "-t", "task", "--parent", epicID, "--json")
	if err != nil {
		return fmt.Errorf("failed to create %s item '%s': %w", itemType, description, err)
	}

	var result struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		return fmt.Errorf("failed to parse create output: %w", err)
	}

	if result.ID == "" {
		return fmt.Errorf("create returned empty ID for %s item '%s'", itemType, description)
	}

	return nil
}

// AddWithPattern adds a new item with pattern tracking for refactor items.
func (p *Provider) AddWithPattern(description, itemType, patternID string) error {
	p.reconcile()
	if err := p.addItem(description, itemType); err != nil {
		return err
	}

	if itemType != backlog.ItemTypeRefactor || patternID == "" {
		p.GenerateSnapshot()
		return nil
	}

	currentCount, err := p.GetPatternCount(patternID)
	if err != nil {
		return fmt.Errorf("failed to get pattern count: %w", err)
	}

	epicID, ok := p.epicCache.Epics[itemType]
	if !ok {
		return fmt.Errorf("epic not found for type: %s", itemType)
	}

	output, err := p.commandRunner.Run("list", "--parent", epicID, "--json")
	if err != nil {
		return fmt.Errorf("failed to list %s items: %w", itemType, err)
	}

	var items []struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	}
	if err := json.Unmarshal([]byte(output), &items); err != nil {
		return fmt.Errorf("failed to parse list output: %w", err)
	}

	var taskID string
	for i := len(items) - 1; i >= 0; i-- {
		if items[i].Title == description {
			taskID = items[i].ID
			break
		}
	}

	if taskID == "" {
		return fmt.Errorf("failed to find newly created task for '%s'", description)
	}

	notes := fmt.Sprintf("PATTERN:%s COUNT:%d", patternID, currentCount+1)
	_, err = p.commandRunner.Run("update", taskID, "--notes", notes)
	if err != nil {
		return fmt.Errorf("failed to update task notes: %w", err)
	}

	p.GenerateSnapshot()
	return nil
}

// AddMultiple adds multiple items to the specified section.
func (p *Provider) AddMultiple(descriptions []string, itemType string) error {
	if len(descriptions) == 0 {
		return nil
	}

	p.reconcile()

	if err := p.ensureEpics(); err != nil {
		return fmt.Errorf("failed to ensure epics: %w", err)
	}

	if !backlog.IsValidType(itemType) {
		return fmt.Errorf("invalid type: %s. Valid types: feature, bugfix, decomposed, refactor, future", itemType)
	}

	epicID, ok := p.epicCache.Epics[itemType]
	if !ok {
		return fmt.Errorf("epic not found for type: %s", itemType)
	}

	for _, description := range descriptions {
		output, err := p.commandRunner.Run("create", description, "-t", "task", "--parent", epicID, "--json")
		if err != nil {
			return fmt.Errorf("failed to create %s item '%s': %w", itemType, description, err)
		}

		var result struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal([]byte(output), &result); err != nil {
			return fmt.Errorf("failed to parse create output: %w", err)
		}

		if result.ID == "" {
			return fmt.Errorf("create returned empty ID for %s item '%s'", itemType, description)
		}
	}

	p.GenerateSnapshot()
	return nil
}

// Complete marks an item as completed.
func (p *Provider) Complete(itemText string) error {
	p.reconcile()

	if err := p.ensureEpics(); err != nil {
		return fmt.Errorf("failed to ensure epics: %w", err)
	}

	epicTypes := []string{
		backlog.ItemTypeFeature,
		backlog.ItemTypeBugfix,
		backlog.ItemTypeDecomposed,
		backlog.ItemTypeRefactor,
		backlog.ItemTypeFuture,
	}

	for _, itemType := range epicTypes {
		epicID, ok := p.epicCache.Epics[itemType]
		if !ok {
			continue
		}

		output, err := p.commandRunner.Run("list", "--parent", epicID, "--json")
		if err != nil {
			return fmt.Errorf("failed to list %s items: %w", itemType, err)
		}

		var items []struct {
			ID    string `json:"id"`
			Title string `json:"title"`
		}
		if err := json.Unmarshal([]byte(output), &items); err != nil {
			return fmt.Errorf("failed to parse list output for %s: %w", itemType, err)
		}

		for _, item := range items {
			if item.Title == itemText {
				_, err := p.commandRunner.Run("close", item.ID, "--reason", "Completed")
				if err != nil {
					return fmt.Errorf("failed to close task '%s': %w", itemText, err)
				}
				p.GenerateSnapshot()
				return nil
			}
		}
	}

	return fmt.Errorf("item not found: %s", itemText)
}

// Cleanup counts closed items. Beads doesn't delete closed items.
func (p *Provider) Cleanup(itemType string) (int, error) {
	p.reconcile()

	if err := p.ensureEpics(); err != nil {
		return 0, fmt.Errorf("failed to ensure epics: %w", err)
	}

	var queryTypes []string
	if itemType != "" {
		if !backlog.IsValidType(itemType) {
			return 0, fmt.Errorf("invalid type: %s. Valid types: feature, bugfix, decomposed, refactor, future", itemType)
		}
		queryTypes = []string{itemType}
	} else {
		queryTypes = []string{
			backlog.ItemTypeFeature,
			backlog.ItemTypeBugfix,
			backlog.ItemTypeDecomposed,
			backlog.ItemTypeRefactor,
			backlog.ItemTypeFuture,
		}
	}

	closedCount := 0
	for _, itemType := range queryTypes {
		epicID, ok := p.epicCache.Epics[itemType]
		if !ok {
			continue
		}

		output, err := p.commandRunner.Run("list", "--parent", epicID, "--json")
		if err != nil {
			return 0, fmt.Errorf("failed to list %s items: %w", itemType, err)
		}

		var items []struct {
			ID     string `json:"id"`
			Status string `json:"status"`
		}
		if err := json.Unmarshal([]byte(output), &items); err != nil {
			return 0, fmt.Errorf("failed to parse list output for %s: %w", itemType, err)
		}

		for _, item := range items {
			if item.Status == "closed" {
				closedCount++
			}
		}
	}

	return closedCount, nil
}

// GetPatternCount returns the occurrence count for a pattern ID.
func (p *Provider) GetPatternCount(patternID string) (int, error) {
	p.reconcile()

	if err := p.ensureEpics(); err != nil {
		return 0, fmt.Errorf("failed to ensure epics: %w", err)
	}

	epicID, ok := p.epicCache.Epics[backlog.ItemTypeRefactor]
	if !ok {
		return 0, nil
	}

	output, err := p.commandRunner.Run("list", "--parent", epicID, "--json")
	if err != nil {
		return 0, fmt.Errorf("failed to list refactor items: %w", err)
	}

	var items []struct {
		ID    string `json:"id"`
		Notes string `json:"notes"`
	}
	if err := json.Unmarshal([]byte(output), &items); err != nil {
		return 0, fmt.Errorf("failed to parse list output: %w", err)
	}

	maxCount := 0
	patternPrefix := fmt.Sprintf("PATTERN:%s", patternID)
	for _, item := range items {
		if strings.Contains(item.Notes, patternPrefix) {
			parts := strings.Split(item.Notes, " ")
			for _, part := range parts {
				if strings.HasPrefix(part, "COUNT:") {
					countStr := strings.TrimPrefix(part, "COUNT:")
					var count int
					if _, err := fmt.Sscanf(countStr, "%d", &count); err == nil {
						if count > maxCount {
							maxCount = count
						}
					}
				}
			}
		}
	}

	return maxCount, nil
}

// Decompose adds multiple sub-intents with dependency tracking.
func (p *Provider) Decompose(jsonInput string) error {
	p.reconcile()

	var decompose struct {
		Action     string `json:"action"`
		SubIntents []struct {
			Intent         string   `json:"intent"`
			Rationale      string   `json:"rationale"`
			EstimatedFiles int      `json:"estimated_files"`
			Dependencies   []string `json:"dependencies"`
		} `json:"sub_intents"`
		DecompositionRationale string `json:"decomposition_rationale"`
	}

	if err := json.Unmarshal([]byte(jsonInput), &decompose); err != nil {
		return fmt.Errorf("parsing decompose JSON: %w", err)
	}

	if len(decompose.SubIntents) == 0 {
		return fmt.Errorf("no sub-intents found in decompose input")
	}

	if err := p.ensureEpics(); err != nil {
		return fmt.Errorf("failed to ensure epics: %w", err)
	}

	epicID, ok := p.epicCache.Epics[backlog.ItemTypeDecomposed]
	if !ok {
		return fmt.Errorf("decomposed epic not found in cache")
	}

	intentToTaskID := make(map[string]string, len(decompose.SubIntents))
	for _, subIntent := range decompose.SubIntents {
		output, err := p.commandRunner.Run("create", subIntent.Intent, "-t", "task", "--parent", epicID, "--json")
		if err != nil {
			return fmt.Errorf("failed to create task for sub-intent '%s': %w", subIntent.Intent, err)
		}

		var result struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal([]byte(output), &result); err != nil {
			return fmt.Errorf("failed to parse create output for '%s': %w", subIntent.Intent, err)
		}

		if result.ID == "" {
			return fmt.Errorf("bd create returned empty ID for sub-intent '%s'", subIntent.Intent)
		}

		intentToTaskID[subIntent.Intent] = result.ID
	}

	for _, subIntent := range decompose.SubIntents {
		if len(subIntent.Dependencies) == 0 {
			continue
		}
		taskID := intentToTaskID[subIntent.Intent]
		for _, depIntent := range subIntent.Dependencies {
			depTaskID, ok := intentToTaskID[depIntent]
			if !ok {
				return fmt.Errorf("dependency task not found: '%s' referenced by '%s'", depIntent, subIntent.Intent)
			}
			_, err := p.commandRunner.Run("dep", "add", taskID, depTaskID)
			if err != nil {
				return fmt.Errorf("failed to add dependency from '%s' to '%s': %w", subIntent.Intent, depIntent, err)
			}
		}
	}

	p.GenerateSnapshot()
	return nil
}
