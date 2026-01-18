package git

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

// EnsureEntry ensures an entry exists in a gitignore file.
// Creates the file if it doesn't exist, skips if entry already present.
func EnsureEntry(fs afero.Fs, dir, entry string) error {
	path := filepath.Join(dir, ".gitignore")

	content, err := afero.ReadFile(fs, path)
	if os.IsNotExist(err) {
		return afero.WriteFile(fs, path, []byte(entry+"\n"), 0644)
	}
	if err != nil {
		return err
	}

	// Check if entry already exists
	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) == strings.TrimSpace(entry) {
			return nil
		}
	}

	// Append entry with proper newline handling
	if len(content) > 0 && content[len(content)-1] != '\n' {
		entry = "\n" + entry
	}
	return afero.WriteFile(fs, path, append(content, []byte(entry+"\n")...), 0644)
}
