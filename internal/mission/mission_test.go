package mission

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMission_GetScope(t *testing.T) {
	testCases := []struct {
		name     string
		body     string
		expected []string
	}{
		{
			name: "Standard Scope",
			body: `
## INTENT
Some intent.

## SCOPE
- internal/checkpoint/service.go
- internal/checkpoint/service_test.go

## PLAN
- Step 1
`,
			expected: []string{
				"internal/checkpoint/service.go",
				"internal/checkpoint/service_test.go",
			},
		},
		{
			name: "Scope with Asterisks",
			body: `
## SCOPE
* cmd/checkpoint.go
* docs/design/m-apply.md
`,
			expected: []string{
				"cmd/checkpoint.go",
				"docs/design/m-apply.md",
			},
		},
		{
			name: "Mixed and Messy Scope",
			body: `
## SCOPE

- file1.go

*    file2.go
  - file3.go

`,
			expected: []string{
				"file1.go",
				"file2.go",
				"file3.go",
			},
		},
		{
			name:     "No Scope Section",
			body:     "## INTENT\nNo scope here.",
			expected: nil,
		},
		{
			name:     "Empty Scope Section",
			body:     "## SCOPE\n\n## PLAN",
			expected: nil,
		},
		{
			name: "Scope with No Files",
			body: `
## SCOPE
`,
			expected: nil,
		},
		{
			name: "Lowercase Scope Header",
			body: `
## scope
- file1.go
`,
			expected: []string{"file1.go"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := &Mission{Body: tc.body}
			scope := m.GetScope()
			assert.Equal(t, tc.expected, scope)
		})
	}
}

func TestMission_DomainsField(t *testing.T) {
	m := &Mission{}

	// Test setting domains field
	m.Domains = "security,performance"
	assert.Equal(t, "security,performance", m.Domains)

	// Test empty domains field
	m2 := &Mission{Domains: ""}
	assert.Equal(t, "", m2.Domains)
}
