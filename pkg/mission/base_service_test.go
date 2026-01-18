package mission

import (
	"testing"

	"github.com/spf13/afero"
)

func TestNewBaseService(t *testing.T) {
	fs := afero.NewMemMapFs()
	missionDir := "/test/mission"

	service := NewBaseService(fs, missionDir)

	if service == nil {
		t.Fatal("NewBaseService returned nil")
	}

	// Verify filesystem
	if service.FS() != fs {
		t.Errorf("FS() returned wrong filesystem: got %p, want %p", service.FS(), fs)
	}

	// Verify mission directory
	if service.MissionDir() != missionDir {
		t.Errorf("MissionDir() returned wrong path: got %q, want %q", service.MissionDir(), missionDir)
	}

	// Verify mission path is constructed correctly
	expectedPath := "/test/mission/mission.md"
	if service.MissionPath() != expectedPath {
		t.Errorf("MissionPath() returned wrong path: got %q, want %q", service.MissionPath(), expectedPath)
	}
}

func TestNewBaseServiceWithPath(t *testing.T) {
	fs := afero.NewMemMapFs()
	missionDir := "/test/mission"
	customPath := "/custom/location/mission.md"

	service := NewBaseServiceWithPath(fs, missionDir, customPath)

	if service == nil {
		t.Fatal("NewBaseServiceWithPath returned nil")
	}

	// Verify filesystem
	if service.FS() != fs {
		t.Errorf("FS() returned wrong filesystem: got %p, want %p", service.FS(), fs)
	}

	// Verify mission directory
	if service.MissionDir() != missionDir {
		t.Errorf("MissionDir() returned wrong path: got %q, want %q", service.MissionDir(), missionDir)
	}

	// Verify custom path is used
	if service.MissionPath() != customPath {
		t.Errorf("MissionPath() returned wrong path: got %q, want %q", service.MissionPath(), customPath)
	}
}

func TestBaseServiceAccessors(t *testing.T) {
	fs := afero.NewOsFs()
	missionDir := "/test/mission"

	service := NewBaseService(fs, missionDir)

	// Test that accessors return consistent values
	fs1 := service.FS()
	fs2 := service.FS()
	if fs1 != fs2 {
		t.Error("FS() returned different values on multiple calls")
	}

	dir1 := service.MissionDir()
	dir2 := service.MissionDir()
	if dir1 != dir2 {
		t.Error("MissionDir() returned different values on multiple calls")
	}

	path1 := service.MissionPath()
	path2 := service.MissionPath()
	if path1 != path2 {
		t.Error("MissionPath() returned different values on multiple calls")
	}
}

func TestBaseServicePathConstruction(t *testing.T) {
	tests := []struct {
		name         string
		missionDir   string
		expectedPath string
	}{
		{
			name:         "simple path",
			missionDir:   "/mission",
			expectedPath: "/mission/mission.md",
		},
		{
			name:         "nested path",
			missionDir:   "/project/.mission",
			expectedPath: "/project/.mission/mission.md",
		},
		{
			name:         "relative path",
			missionDir:   ".mission",
			expectedPath: ".mission/mission.md",
		},
		{
			name:         "trailing slash",
			missionDir:   "/mission/",
			expectedPath: "/mission/mission.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			service := NewBaseService(fs, tt.missionDir)

			if service.MissionPath() != tt.expectedPath {
				t.Errorf("MissionPath() returned wrong path: got %q, want %q", service.MissionPath(), tt.expectedPath)
			}
		})
	}
}
