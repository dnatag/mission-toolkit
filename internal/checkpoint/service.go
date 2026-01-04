package checkpoint

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/dnatag/mission-toolkit/internal/git"
	"github.com/dnatag/mission-toolkit/internal/mission"
	"github.com/spf13/afero"
)

// Service handles checkpoint operations
type Service struct {
	fs            afero.Fs
	missionDir    string
	missionReader *mission.Reader
	git           git.GitClient
}

// NewService creates a new checkpoint service using CmdGitClient (production)
func NewService(fs afero.Fs, missionDir string) (*Service, error) {
	return &Service{
		fs:            fs,
		missionDir:    missionDir,
		missionReader: mission.NewReader(fs),
		git:           git.NewCmdGitClient("."),
	}, nil
}

// NewServiceWithGit creates a new checkpoint service with a specific GitClient (testing)
func NewServiceWithGit(fs afero.Fs, missionDir string, gitClient git.GitClient) *Service {
	return &Service{
		fs:            fs,
		missionDir:    missionDir,
		missionReader: mission.NewReader(fs),
		git:           gitClient,
	}
}

// Create creates a new checkpoint for the current mission
func (s *Service) Create(missionID string) (string, error) {
	scope, err := s.getScope()
	if err != nil {
		return "", err
	}

	// Only stage files that actually exist to avoid "pathspec did not match" errors
	existingFiles, err := s.filterExistingFiles(scope)
	if err != nil {
		return "", fmt.Errorf("filtering existing files: %w", err)
	}

	num, err := s.getNextCheckpointNumber(missionID)
	if err != nil {
		return "", fmt.Errorf("getting next checkpoint number: %w", err)
	}

	checkpointName := fmt.Sprintf("%s-%d", missionID, num)

	if err := s.git.Add(existingFiles); err != nil {
		return "", fmt.Errorf("staging files: %w", err)
	}

	commitHash, err := s.git.Commit(fmt.Sprintf("checkpoint: %s", checkpointName))
	if err != nil {
		if errors.Is(err, git.ErrNoChanges) {
			// If no changes, tag the current HEAD
			// We need to get HEAD hash. GitClient doesn't have GetHeadHash directly,
			// but GetTagCommit("HEAD") works with our implementation.
			commitHash, err = s.git.GetTagCommit("HEAD")
			if err != nil {
				return "", fmt.Errorf("getting HEAD hash: %w", err)
			}
		} else {
			return "", fmt.Errorf("creating checkpoint commit: %w", err)
		}
	}

	if err := s.git.CreateTag(checkpointName, commitHash); err != nil {
		return "", fmt.Errorf("creating checkpoint tag: %w", err)
	}

	return checkpointName, nil
}

// Restore reverts working directory to specified checkpoint
func (s *Service) Restore(checkpointName string) error {
	scope, err := s.getScope()
	if err != nil {
		return err
	}

	if err := s.git.Restore(checkpointName, scope); err != nil {
		return fmt.Errorf("restoring files: %w", err)
	}

	return nil
}

// Clear removes all checkpoint tags for the specified mission
func (s *Service) Clear(missionID string) (int, error) {
	tags, err := s.git.ListTags(missionID + "-")
	if err != nil {
		return 0, fmt.Errorf("listing tags: %w", err)
	}

	for i, tag := range tags {
		if err := s.git.DeleteTag(tag); err != nil {
			return i, fmt.Errorf("deleting tag %s: %w", tag, err)
		}
	}

	return len(tags), nil
}

// Consolidate creates a final commit with all changes from the mission and clears checkpoints.
func (s *Service) Consolidate(missionID, message string) (string, error) {
	scope, err := s.getScope()
	if err != nil {
		return "", err
	}

	existingFiles, err := s.filterExistingFiles(scope)
	if err != nil {
		return "", fmt.Errorf("filtering existing files: %w", err)
	}

	if err := s.git.Add(existingFiles); err != nil {
		return "", fmt.Errorf("staging final files: %w", err)
	}

	finalCommitHash, err := s.git.Commit(message)
	if err != nil {
		return "", fmt.Errorf("creating final commit: %w", err)
	}

	if _, err := s.Clear(missionID); err != nil {
		fmt.Printf("Warning: failed to clear all checkpoints: %v\n", err)
	}

	return finalCommitHash, nil
}

// getScope reads mission and returns scope files
func (s *Service) getScope() ([]string, error) {
	m, err := s.missionReader.Read(fmt.Sprintf("%s/mission.md", s.missionDir))
	if err != nil {
		return nil, fmt.Errorf("reading mission: %w", err)
	}
	scope := m.GetScope()
	if len(scope) == 0 {
		return nil, fmt.Errorf("no files in mission scope")
	}
	return scope, nil
}

// filterExistingFiles returns only the files from the list that exist on the filesystem
func (s *Service) filterExistingFiles(files []string) ([]string, error) {
	var existing []string
	for _, file := range files {
		exists, err := afero.Exists(s.fs, file)
		if err != nil {
			return nil, fmt.Errorf("checking file existence %s: %w", file, err)
		}
		if exists {
			existing = append(existing, file)
		}
	}
	return existing, nil
}

// getNextCheckpointNumber finds the next available checkpoint number
func (s *Service) getNextCheckpointNumber(missionID string) (int, error) {
	tags, err := s.git.ListTags(missionID + "-")
	if err != nil {
		return 0, fmt.Errorf("listing tags: %w", err)
	}

	maxNum := 0
	for _, tag := range tags {
		if idx := strings.LastIndex(tag, "-"); idx != -1 {
			if num, err := strconv.Atoi(tag[idx+1:]); err == nil && num > maxNum {
				maxNum = num
			}
		}
	}

	return maxNum + 1, nil
}
