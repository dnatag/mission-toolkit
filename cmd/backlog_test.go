package cmd

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBacklogAdd(t *testing.T) {
	// Setup
	fs := afero.NewMemMapFs()
	missionFs = fs
	missionDir = ".mission"

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Execute
	cmd := backlogAddCmd
	cmd.SetArgs([]string{"Test item"})
	err := cmd.Execute()

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)

	// Verify
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "Added to backlog: Test item")

	// Check file content
	content, err := afero.ReadFile(fs, filepath.Join(missionDir, "backlog.md"))
	require.NoError(t, err)
	assert.Contains(t, string(content), "- [ ] Test item")
}

func TestBacklogCheck(t *testing.T) {
	// Setup
	fs := afero.NewMemMapFs()
	missionFs = fs
	missionDir = ".mission"
	backlogPath := filepath.Join(missionDir, "backlog.md")

	err := fs.MkdirAll(missionDir, 0755)
	require.NoError(t, err)
	err = afero.WriteFile(fs, backlogPath, []byte("# Mission Backlog\n\n- [ ] Existing item\n"), 0644)
	require.NoError(t, err)

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Execute
	cmd := backlogCheckCmd
	err = cmd.Execute()

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)

	// Verify
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "# Mission Backlog")
	assert.Contains(t, buf.String(), "- [ ] Existing item")
}
