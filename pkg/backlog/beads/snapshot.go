package beads

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	backlog "github.com/dnatag/mission-toolkit/pkg/backlog"
)

// GenerateSnapshot writes the current Beads state to .mission/backlog.md.
// Errors are logged but don't fail operations.
func (p *Provider) GenerateSnapshot() {
	p.reconciling = true
	defer func() { p.reconciling = false }()

	var content strings.Builder

	content.WriteString("---\n")
	content.WriteString(fmt.Sprintf("last_updated: %s\n", time.Now().Format(time.RFC3339)))
	content.WriteString("source: beads\n")
	content.WriteString("---\n\n")

	content.WriteString("<!-- MANAGED BY BEADS — Source of truth is 'bd'. Use 'm backlog' commands. -->\n\n")

	content.WriteString("# Backlog\n\n")

	sections := []struct {
		title string
		desc  string
		typ   string
	}{
		{"FEATURES", "User-defined feature requests and enhancements", backlog.ItemTypeFeature},
		{"BUGFIXES", "Bug reports and issues to be fixed", backlog.ItemTypeBugfix},
		{"DECOMPOSED INTENTS", "Sub-intents from Track 4 Epic requests that need separate missions", backlog.ItemTypeDecomposed},
		{"REFACTORING OPPORTUNITIES", "Detected duplication patterns that need DRY missions", backlog.ItemTypeRefactor},
		{"FUTURE ENHANCEMENTS", "Ideas and improvements for later consideration", backlog.ItemTypeFuture},
		{"COMPLETED", "History of completed backlog items", "completed"},
	}

	for _, section := range sections {
		content.WriteString(fmt.Sprintf("## %s\n", section.title))
		content.WriteString(fmt.Sprintf("(%s)\n", section.desc))

		var items []string
		var err error

		if section.typ == "completed" {
			items, err = p.List([]string{"completed"}, nil)
		} else {
			items, err = p.List([]string{section.typ}, []string{"completed"})
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to list %s items for snapshot: %v\n", section.typ, err)
			items = []string{}
		}

		for _, item := range items {
			content.WriteString(item + "\n")
		}

		content.WriteString("\n")
	}

	backlogPath := filepath.Join(p.projectRoot, ".mission", "backlog.md")

	// Safety check: don't overwrite existing backlog with empty content.
	// If all sections returned 0 items, Beads may be in a bad state (e.g., DB lock).
	totalItems := strings.Count(content.String(), "- [ ] ") + strings.Count(content.String(), "- [x] ")
	if totalItems == 0 {
		if info, err := os.Stat(backlogPath); err == nil && info.Size() > 0 {
			fmt.Fprintf(os.Stderr, "Warning: snapshot skipped — Beads returned 0 items but backlog.md exists\n")
			return
		}
	}

	if err := os.MkdirAll(filepath.Dir(backlogPath), 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to create .mission directory for snapshot: %v\n", err)
		return
	}

	if err := os.WriteFile(backlogPath, []byte(content.String()), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to write snapshot to %s: %v\n", backlogPath, err)
		return
	}

	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(content.String())))
	if err := os.WriteFile(p.hashPath, []byte(hash), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to write hash to %s: %v\n", p.hashPath, err)
	}
}
