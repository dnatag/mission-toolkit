package mission

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

// Reader handles reading and parsing mission.md files
type Reader struct {
	*BaseService
}

// NewReader creates a new mission reader
func NewReader(fs afero.Fs, path string) *Reader {
	return &Reader{
		BaseService: NewBaseServiceWithPath(fs, "", path),
	}
}

// Read reads and parses a mission file into a Mission struct
func (r *Reader) Read() (*Mission, error) {
	data, err := afero.ReadFile(r.FS(), r.MissionPath())
	if err != nil {
		return nil, fmt.Errorf("failed to read mission file: %w", err)
	}

	return r.parse(data)
}

// parse parses mission data with frontmatter and body
func (r *Reader) parse(data []byte) (*Mission, error) {
	// Handle empty files
	if len(data) == 0 {
		return nil, fmt.Errorf("empty mission file")
	}

	// Check for YAML frontmatter format (starts with ---)
	if bytes.HasPrefix(data, []byte("---\n")) {
		// Find the closing --- on its own line (may be followed by empty lines)
		// This handles cases where --- appears in the middle of YAML
		re := regexp.MustCompile(`(?m)^\s*---\s*$`)
		matches := re.FindAllIndex(data, -1)

		if len(matches) < 2 {
			return nil, fmt.Errorf("invalid frontmatter format: expected opening and closing ---")
		}

		// Second match is the closing ---
		closingDash := matches[1]

		// Frontmatter is everything between the opening --- and the closing ---
		frontmatter := data[4:closingDash[0]]

		// Body starts after the closing ---
		bodyStart := closingDash[1]
		// Skip any leading whitespace/newlines after the closing ---
		for bodyStart < len(data) && (data[bodyStart] == '\n' || data[bodyStart] == ' ') {
			bodyStart++
		}
		body := data[bodyStart:]

		mission := &Mission{Body: string(body)}
		if err := yaml.Unmarshal(frontmatter, mission); err != nil {
			return nil, fmt.Errorf("failed to unmarshal frontmatter: %w", err)
		}

		return mission, nil
	}

	// Check for legacy # MISSION format (including # MISSION ARCHIVE and # MISSION: title)
	if bytes.HasPrefix(data, []byte("# MISSION")) {
		return r.parseLegacy(data)
	}

	return nil, fmt.Errorf("no frontmatter found in mission file")
}

// parseLegacy parses the legacy # MISSION format and its variants
func (r *Reader) parseLegacy(data []byte) (*Mission, error) {
	lines := bytes.Split(data, []byte("\n"))
	var bodyLines []string
	var metadataLines []string
	bodyStarted := false

	for _, line := range lines {
		trimmed := bytes.TrimSpace(line)
		lineStr := string(line)

		// Body starts with ## INTENT or other ## headers
		if len(trimmed) >= 2 && trimmed[0] == '#' && trimmed[1] == '#' {
			bodyStarted = true
		}

		if bodyStarted {
			bodyLines = append(bodyLines, lineStr)
		} else if len(trimmed) > 0 {
			// Skip all # MISSION header variants (# MISSION, # MISSION ARCHIVE, # MISSION: Title)
			if !bytes.HasPrefix(trimmed, []byte("# MISSION")) {
				// Include metadata lines that contain colons (key: value format)
				if strings.Contains(lineStr, ":") {
					metadataLines = append(metadataLines, lineStr)
				}
			}
		}
	}

	mission := &Mission{Body: strings.Join(bodyLines, "\n")}

	// Parse metadata lines with consistent error handling
	for _, line := range metadataLines {
		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "type:"):
			mission.Type = strings.TrimSpace(strings.TrimPrefix(line, "type:"))
		case strings.HasPrefix(line, "track:"):
			fmt.Sscanf(strings.TrimPrefix(line, "track:"), "%d", &mission.Track)
		case strings.HasPrefix(line, "iteration:"):
			fmt.Sscanf(strings.TrimPrefix(line, "iteration:"), "%d", &mission.Iteration)
		case strings.HasPrefix(line, "status:"):
			mission.Status = strings.TrimSpace(strings.TrimPrefix(line, "status:"))
		case strings.HasPrefix(line, "id:"):
			mission.ID = strings.TrimSpace(strings.TrimPrefix(line, "id:"))
		case strings.HasPrefix(line, "completed_at:"), strings.HasPrefix(line, "duration_minutes:"):
			// Parse additional metadata fields if present but don't fail on parsing errors
			// These fields are informational and not critical for mission functionality
		}
	}

	// Ensure mission has a valid ID
	if mission.ID == "" {
		mission.ID = "legacy-mission"
	}

	return mission, nil
}

// GetMissionID reads the mission ID from a mission file
func (r *Reader) GetMissionID() (string, error) {
	mission, err := r.Read()
	if err != nil {
		return "", err
	}
	return mission.ID, nil
}

// GetMissionStatus reads the mission status from a mission file
func (r *Reader) GetMissionStatus() (string, error) {
	mission, err := r.Read()
	if err != nil {
		return "", err
	}
	return mission.Status, nil
}

// ReadIntent extracts the INTENT section from a mission file
func (r *Reader) ReadIntent() (string, error) {
	mission, err := r.Read()
	if err != nil {
		return "", fmt.Errorf("reading mission file: %w", err)
	}
	intent := mission.GetIntent()
	if intent == "" {
		return "", fmt.Errorf("no intent found in mission")
	}
	return intent, nil
}

// ReadScope extracts the SCOPE section from a mission file
func (r *Reader) ReadScope() (string, error) {
	mission, err := r.Read()
	if err != nil {
		return "", fmt.Errorf("reading mission file: %w", err)
	}
	scope := mission.GetScope()
	if len(scope) == 0 {
		return "", fmt.Errorf("no scope found in mission")
	}
	return strings.Join(scope, "\n"), nil
}
