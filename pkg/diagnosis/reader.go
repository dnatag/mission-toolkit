package diagnosis

import (
	"fmt"
	"os"

	"github.com/dnatag/mission-toolkit/pkg/md"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

// ReadDiagnosis reads and parses a diagnosis.md file using pkg/md abstraction.
func ReadDiagnosis(fs afero.Fs, diagnosisPath string) (*Diagnosis, error) {
	content, err := afero.ReadFile(fs, diagnosisPath)
	if err != nil {
		return nil, fmt.Errorf("reading diagnosis file: %w", err)
	}

	doc, err := md.Parse(content)
	if err != nil {
		return nil, fmt.Errorf("parsing diagnosis: %w", err)
	}

	// Convert frontmatter map to Diagnosis struct via YAML marshaling
	// This ensures type conversions (especially time.Time) are handled correctly
	yamlData, err := yaml.Marshal(doc.Frontmatter)
	if err != nil {
		return nil, fmt.Errorf("marshaling frontmatter: %w", err)
	}

	var diag Diagnosis
	if err := yaml.Unmarshal(yamlData, &diag); err != nil {
		return nil, fmt.Errorf("unmarshaling frontmatter: %w", err)
	}

	diag.Body = doc.Body
	return &diag, nil
}

// DiagnosisExists checks if a diagnosis.md file exists
func DiagnosisExists(fs afero.Fs, diagnosisPath string) (bool, error) {
	_, err := fs.Stat(diagnosisPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
