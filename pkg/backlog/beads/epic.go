package beads

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	backlog "github.com/dnatag/mission-toolkit/pkg/backlog"
)

// Epic title constants.
const (
	EpicTitleFeatures                 = "Features"
	EpicTitleBugfixes                 = "Bugfixes"
	EpicTitleDecomposedIntents        = "Decomposed Intents"
	EpicTitleRefactoringOpportunities = "Refactoring Opportunities"
	EpicTitleFutureEnhancements       = "Future Enhancements"
)

// EpicCache stores epic ID mappings.
type EpicCache struct {
	Epics map[string]string `json:"epics"`
}

// EpicTypeDef pairs an item type with its epic title.
type EpicTypeDef struct {
	ItemType  string
	EpicTitle string
}

// OrderedEpicTypes returns epic type definitions in a deterministic order.
func OrderedEpicTypes() []EpicTypeDef {
	return []EpicTypeDef{
		{backlog.ItemTypeFeature, EpicTitleFeatures},
		{backlog.ItemTypeBugfix, EpicTitleBugfixes},
		{backlog.ItemTypeDecomposed, EpicTitleDecomposedIntents},
		{backlog.ItemTypeRefactor, EpicTitleRefactoringOpportunities},
		{backlog.ItemTypeFuture, EpicTitleFutureEnhancements},
	}
}

func (p *Provider) ensureEpics() error {
	if err := p.loadEpicCache(); err == nil {
		return nil
	}

	for _, et := range OrderedEpicTypes() {
		output, err := p.commandRunner.Run("create", et.EpicTitle, "-t", "epic", "--json")
		if err != nil {
			return fmt.Errorf("failed to create %s epic (%s): %w", et.ItemType, et.EpicTitle, err)
		}

		var result struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal([]byte(output), &result); err != nil {
			return fmt.Errorf("failed to parse epic creation output for %s (%s): %w", et.ItemType, et.EpicTitle, err)
		}

		p.epicCache.Epics[et.ItemType] = result.ID
	}

	return p.saveEpicCache()
}

func (p *Provider) loadEpicCache() error {
	data, err := os.ReadFile(p.cachePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, p.epicCache)
}

func (p *Provider) saveEpicCache() error {
	if err := os.MkdirAll(filepath.Dir(p.cachePath), 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	data, err := json.MarshalIndent(p.epicCache, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal epic cache: %w", err)
	}

	return os.WriteFile(p.cachePath, data, 0644)
}
