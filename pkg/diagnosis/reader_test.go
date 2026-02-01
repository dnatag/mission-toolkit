package diagnosis

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestDiagnosisExists(t *testing.T) {
	tests := []struct {
		name       string
		setupFile  bool
		wantExists bool
		wantErr    bool
	}{
		{
			name:       "file exists",
			setupFile:  true,
			wantExists: true,
			wantErr:    false,
		},
		{
			name:       "file does not exist",
			setupFile:  false,
			wantExists: false,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			path := ".mission/diagnosis.md"

			if tt.setupFile {
				fs.MkdirAll(".mission", 0755)
				afero.WriteFile(fs, path, []byte("test content"), 0644)
			}

			exists, err := DiagnosisExists(fs, path)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantExists, exists)
		})
	}
}
