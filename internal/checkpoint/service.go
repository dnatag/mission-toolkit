package checkpoint

import (
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/spf13/afero"
)

// Service handles checkpoint operations using go-git
type Service struct {
	fs         afero.Fs
	missionDir string
	repo       *git.Repository
}

// NewService creates a new checkpoint service
func NewService(fs afero.Fs, missionDir string) (*Service, error) {
	repo, err := git.PlainOpen(".")
	if err != nil {
		return nil, fmt.Errorf("opening git repository: %w", err)
	}

	return &Service{
		fs:         fs,
		missionDir: missionDir,
		repo:       repo,
	}, nil
}

// Create creates a new checkpoint for the current mission
func (s *Service) Create(missionID string) (string, error) {
	// Get next checkpoint number
	num, err := s.getNextCheckpointNumber(missionID)
	if err != nil {
		return "", fmt.Errorf("getting next checkpoint number: %w", err)
	}

	checkpointName := fmt.Sprintf("%s-%d", missionID, num)

	// Get current HEAD
	head, err := s.repo.Head()
	if err != nil {
		return "", fmt.Errorf("getting HEAD: %w", err)
	}

	// Get worktree
	wt, err := s.repo.Worktree()
	if err != nil {
		return "", fmt.Errorf("getting worktree: %w", err)
	}

	// Stage all changes
	if err := wt.AddWithOptions(&git.AddOptions{All: true}); err != nil {
		return "", fmt.Errorf("staging changes: %w", err)
	}

	// Create commit
	commitHash, err := wt.Commit(fmt.Sprintf("checkpoint: %s", checkpointName), &git.CommitOptions{})
	if err != nil {
		return "", fmt.Errorf("creating checkpoint commit: %w", err)
	}

	// Create ref
	refName := plumbing.ReferenceName(fmt.Sprintf("refs/checkpoints/%s", checkpointName))
	ref := plumbing.NewHashReference(refName, commitHash)
	if err := s.repo.Storer.SetReference(ref); err != nil {
		return "", fmt.Errorf("creating checkpoint ref: %w", err)
	}

	// Reset to original HEAD (keep working directory unchanged)
	if err := wt.Reset(&git.ResetOptions{
		Commit: head.Hash(),
		Mode:   git.HardReset,
	}); err != nil {
		return "", fmt.Errorf("resetting to HEAD: %w", err)
	}

	return checkpointName, nil
}

// Revert reverts working directory to specified checkpoint
func (s *Service) Revert(checkpointName string) error {
	// Resolve checkpoint ref
	refName := plumbing.ReferenceName(fmt.Sprintf("refs/checkpoints/%s", checkpointName))
	ref, err := s.repo.Reference(refName, true)
	if err != nil {
		return fmt.Errorf("checkpoint not found: %s", checkpointName)
	}

	// Get worktree
	wt, err := s.repo.Worktree()
	if err != nil {
		return fmt.Errorf("getting worktree: %w", err)
	}

	// Hard reset to checkpoint
	if err := wt.Reset(&git.ResetOptions{
		Commit: ref.Hash(),
		Mode:   git.HardReset,
	}); err != nil {
		return fmt.Errorf("reverting to checkpoint: %w", err)
	}

	// Delete checkpoint ref
	if err := s.repo.Storer.RemoveReference(refName); err != nil {
		return fmt.Errorf("removing checkpoint ref: %w", err)
	}

	return nil
}

// Clear removes all checkpoint refs for the specified mission
func (s *Service) Clear(missionID string) (int, error) {
	prefix := fmt.Sprintf("refs/checkpoints/%s-", missionID)
	count := 0

	// List all refs
	refs, err := s.repo.References()
	if err != nil {
		return 0, fmt.Errorf("listing refs: %w", err)
	}

	// Find and delete matching checkpoint refs
	err = refs.ForEach(func(ref *plumbing.Reference) error {
		if strings.HasPrefix(ref.Name().String(), prefix) {
			if err := s.repo.Storer.RemoveReference(ref.Name()); err != nil {
				return fmt.Errorf("removing ref %s: %w", ref.Name(), err)
			}
			count++
		}
		return nil
	})

	if err != nil {
		return count, fmt.Errorf("clearing checkpoints: %w", err)
	}

	return count, nil
}

// getNextCheckpointNumber finds the next available checkpoint number
func (s *Service) getNextCheckpointNumber(missionID string) (int, error) {
	prefix := fmt.Sprintf("refs/checkpoints/%s-", missionID)
	maxNum := 0

	refs, err := s.repo.References()
	if err != nil {
		return 0, fmt.Errorf("listing refs: %w", err)
	}

	err = refs.ForEach(func(ref *plumbing.Reference) error {
		refName := ref.Name().String()
		if strings.HasPrefix(refName, prefix) {
			// Extract number from ref name
			parts := strings.Split(refName, "-")
			if len(parts) > 0 {
				var num int
				fmt.Sscanf(parts[len(parts)-1], "%d", &num)
				if num > maxNum {
					maxNum = num
				}
			}
		}
		return nil
	})

	if err != nil {
		return 0, fmt.Errorf("finding max checkpoint number: %w", err)
	}

	return maxNum + 1, nil
}
