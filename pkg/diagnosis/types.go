package diagnosis

import "time"

// Diagnosis represents a debug investigation with structured findings
type Diagnosis struct {
	ID         string    `yaml:"id"`
	Status     string    `yaml:"status"`
	Confidence string    `yaml:"confidence"`
	Created    time.Time `yaml:"created"`
	Symptom    string    `yaml:"-"`
	Body       string    `yaml:"-"`
}
