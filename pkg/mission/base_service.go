package mission

import (
	"path/filepath"

	"github.com/spf13/afero"
)

// BaseService provides common fields and initialization for mission services.
// It embeds the filesystem, mission directory, and mission path that are
// shared across multiple mission services (Reader, Writer, CheckService, etc.).
type BaseService struct {
	fs          afero.Fs
	missionDir  string
	missionPath string
}

// NewBaseService creates a new BaseService with common mission path initialization.
// The mission path is constructed as <missionDir>/mission.md.
func NewBaseService(fs afero.Fs, missionDir string) *BaseService {
	return &BaseService{
		fs:          fs,
		missionDir:  missionDir,
		missionPath: filepath.Join(missionDir, "mission.md"),
	}
}

// NewBaseServiceWithPath creates a new BaseService with an explicit mission file path.
// Use this constructor when the mission file is not at the default location.
func NewBaseServiceWithPath(fs afero.Fs, missionDir, path string) *BaseService {
	return &BaseService{
		fs:          fs,
		missionDir:  missionDir,
		missionPath: path,
	}
}

// FS returns the filesystem instance.
func (b *BaseService) FS() afero.Fs {
	return b.fs
}

// MissionDir returns the mission directory path.
func (b *BaseService) MissionDir() string {
	return b.missionDir
}

// MissionPath returns the full path to the mission file.
func (b *BaseService) MissionPath() string {
	return b.missionPath
}
