package checkpoint

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/dnatag/mission-toolkit/pkg/git"
	"github.com/dnatag/mission-toolkit/pkg/mission"
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
	missionPath := fmt.Sprintf("%s/mission.md", missionDir)
	return &Service{
		fs:            fs,
		missionDir:    missionDir,
		missionReader: mission.NewReader(fs, missionPath),
		git:           git.NewCmdGitClient("."),
	}, nil
}

// NewServiceWithGit creates a new checkpoint service with a specific GitClient (testing)
func NewServiceWithGit(fs afero.Fs, missionDir string, gitClient git.GitClient) *Service {
	missionPath := fmt.Sprintf("%s/mission.md", missionDir)
	return &Service{
		fs:            fs,
		missionDir:    missionDir,
		missionReader: mission.NewReader(fs, missionPath),
		git:           gitClient,
	}
}

// Create creates a new checkpoint for the current mission
func (s *Service) Create(missionID string) (string, error) {
	stagableFiles, err := s.getStagableScope()
	if err != nil {
		return "", err
	}

	num, err := s.getNextCheckpointNumber(missionID)
	if err != nil {
		return "", fmt.Errorf("getting next checkpoint number: %w", err)
	}

	checkpointName := fmt.Sprintf("%s-%d", missionID, num)

	if err := s.git.Add(stagableFiles); err != nil {
		return "", fmt.Errorf("staging files: %w", err)
	}

	commitHash, err := s.git.Commit(fmt.Sprintf("checkpoint: %s", checkpointName))
	if err != nil {
		if !errors.Is(err, git.ErrNoChanges) {
			return "", fmt.Errorf("creating checkpoint commit: %w", err)
		}
		// No changes, tag current HEAD
		if commitHash, err = s.git.GetTagCommit("HEAD"); err != nil {
			return "", fmt.Errorf("getting HEAD hash: %w", err)
		}
	}

	if err := s.git.CreateTag(checkpointName, commitHash); err != nil {
		return "", fmt.Errorf("creating checkpoint tag: %w", err)
	}

	// Create baseline tag on first checkpoint for easy diff viewing
	if num == 1 {
		if err := s.createBaselineTag(missionID, commitHash); err != nil {
			return "", fmt.Errorf("creating baseline tag: %w", err)
		}
	}

	return checkpointName, nil
}

// createBaselineTag creates a baseline tag for viewing cumulative mission changes
func (s *Service) createBaselineTag(missionID, commitHash string) error {
	baselineTag := fmt.Sprintf("%s-baseline", missionID)
	return s.git.CreateTag(baselineTag, commitHash)
}

// Restore reverts working directory to specified checkpoint
func (s *Service) Restore(checkpointName string) error {
	scope, err := s.getScope()
	if err != nil {
		return err
	}
	return s.git.Restore(checkpointName, scope)
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

// RestoreAll reverts the working directory to the baseline commit (before any checkpoints)
// and deletes all checkpoint tags. This is used by the --all flag to completely undo all
// mission changes and clean up checkpoint history.
func (s *Service) RestoreAll(missionID string) (int, error) {
	baselineTag := fmt.Sprintf("%s-baseline", missionID)
	baselineHash, err := s.git.GetTagCommit(baselineTag)
	if err != nil {
		return 0, fmt.Errorf("getting baseline commit: %w", err)
	}

	scope, err := s.getScope()
	if err != nil {
		return 0, fmt.Errorf("reading mission scope: %w", err)
	}

	if err := s.git.Restore(baselineTag, scope); err != nil {
		return 0, fmt.Errorf("restoring files to baseline: %w", err)
	}

	if err := s.git.SoftReset(baselineHash); err != nil {
		return 0, fmt.Errorf("resetting HEAD to baseline: %w", err)
	}

	count, err := s.Clear(missionID)
	if err != nil {
		return 0, fmt.Errorf("clearing checkpoint tags: %w", err)
	}

	return count, nil
}

// ConsolidateResult contains the result of a consolidate operation
type ConsolidateResult struct {
	CommitHash    string
	UnstagedFiles []string
}

// Consolidate creates a final commit with all changes from the mission and clears checkpoints.
func (s *Service) Consolidate(missionID, message string) (*ConsolidateResult, error) {
	targetHash, err := s.squashCheckpoints(missionID)
	if err != nil {
		return nil, err
	}

	if targetHash != "" {
		if err := s.git.SoftReset(targetHash); err != nil {
			return nil, fmt.Errorf("soft reset to %s failed: %w", targetHash, err)
		}
	}

	stagableFiles, err := s.getStagableScope()
	if err != nil {
		return nil, err
	}

	if err := s.git.Add(stagableFiles); err != nil {
		return nil, fmt.Errorf("staging final files: %w", err)
	}

	finalCommitHash, err := s.git.Commit(message)
	if err != nil {
		return nil, fmt.Errorf("creating final commit: %w", err)
	}

	if _, err := s.Clear(missionID); err != nil {
		fmt.Printf("Warning: failed to clear all checkpoints: %v\n", err)
	}

	unstaged, _ := s.git.GetUnstagedFiles()

	return &ConsolidateResult{
		CommitHash:    finalCommitHash,
		UnstagedFiles: unstaged,
	}, nil
}

// squashCheckpoints finds the initial checkpoint and determines the target commit for squashing.
// Returns the target hash to reset to, or empty string if no squashing is needed.
func (s *Service) squashCheckpoints(missionID string) (string, error) {
	initialCommitHash, err := s.git.GetTagCommit(fmt.Sprintf("%s-1", missionID))
	if err != nil {
		return "", nil // No initial checkpoint found, no squashing needed
	}

	msg, err := s.git.GetCommitMessage(initialCommitHash)
	if err != nil {
		return "", nil // Can't read commit message, skip squashing
	}

	targetHash := initialCommitHash
	if strings.HasPrefix(msg, "checkpoint:") {
		// If we created it, reset to its parent to squash it
		if parentHash, err := s.git.GetCommitParent(initialCommitHash); err == nil && parentHash != "" {
			targetHash = parentHash
		}
	}
	// If we didn't create it (it was a pre-existing commit we tagged), reset to it directly.
	// This preserves the pre-existing commit but squashes subsequent checkpoints.

	return targetHash, nil
}

// getStagableScope reads mission scope and filters to stagable files.
// It combines getScope() and filterStagableFiles() to reduce duplication.
func (s *Service) getStagableScope() ([]string, error) {
	scope, err := s.getScope()
	if err != nil {
		return nil, err
	}

	stagableFiles, err := s.filterStagableFiles(scope)
	if err != nil {
		return nil, fmt.Errorf("filtering stagable files: %w", err)
	}

	return stagableFiles, nil
}

// getScope reads mission and returns scope files
func (s *Service) getScope() ([]string, error) {
	m, err := s.missionReader.Read()
	if err != nil {
		return nil, fmt.Errorf("reading mission: %w", err)
	}
	scope := m.GetScope()
	if len(scope) == 0 {
		return nil, fmt.Errorf("no files in mission scope")
	}
	return scope, nil
}

// filterStagableFiles returns files that exist OR are tracked (for deletion)
func (s *Service) filterStagableFiles(files []string) ([]string, error) {
	var stagable []string
	for _, file := range files {
		exists, err := afero.Exists(s.fs, file)
		if err != nil {
			return nil, fmt.Errorf("checking file existence %s: %w", file, err)
		}

		if exists {
			stagable = append(stagable, file)
			continue
		}

		// If file doesn't exist, check if it's tracked (deleted)
		if tracked, err := s.git.IsTracked(file); err != nil {
			return nil, fmt.Errorf("checking if file is tracked %s: %w", file, err)
		} else if tracked {
			stagable = append(stagable, file)
		}
	}
	return stagable, nil
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
