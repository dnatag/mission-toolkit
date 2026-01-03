package mission

import (
	"bytes"
	"fmt"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

// Reader handles reading and parsing mission.md files
type Reader struct {
	fs afero.Fs
}

// NewReader creates a new mission reader
func NewReader(fs afero.Fs) *Reader {
	return &Reader{fs: fs}
}

// Read reads and parses a mission file into a Mission struct
func (r *Reader) Read(path string) (*Mission, error) {
	data, err := afero.ReadFile(r.fs, path)
	if err != nil {
		return nil, fmt.Errorf("failed to read mission file: %w", err)
	}

	return r.parse(data)
}

// parse parses mission data with frontmatter and body
func (r *Reader) parse(data []byte) (*Mission, error) {
	if !bytes.HasPrefix(data, []byte("---\n")) {
		return nil, fmt.Errorf("no frontmatter found in mission file")
	}

	parts := bytes.SplitN(data[4:], []byte("\n---\n"), 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid frontmatter format")
	}

	mission := &Mission{Body: string(parts[1])}
	if err := yaml.Unmarshal(parts[0], mission); err != nil {
		return nil, fmt.Errorf("failed to unmarshal frontmatter: %w", err)
	}

	return mission, nil
}

// GetMissionID reads the mission ID from a mission file
func (r *Reader) GetMissionID(path string) (string, error) {
	mission, err := r.Read(path)
	if err != nil {
		return "", err
	}
	return mission.ID, nil
}

// GetMissionStatus reads the mission status from a mission file
func (r *Reader) GetMissionStatus(path string) (string, error) {
	mission, err := r.Read(path)
	if err != nil {
		return "", err
	}
	return mission.Status, nil
}
