package mission

type MockGitClient struct {
	commitMessage string
	commitError   error
}

func (m *MockGitClient) Add(files []string) error {
	return nil
}

func (m *MockGitClient) Commit(message string) (string, error) {
	return "mock-hash", nil
}

func (m *MockGitClient) CommitNoVerify(message string) (string, error) {
	return "mock-hash", nil
}

func (m *MockGitClient) CreateTag(name string, commitHash string) error {
	return nil
}

func (m *MockGitClient) Restore(checkpointName string, files []string) error {
	return nil
}

func (m *MockGitClient) ListTags(prefix string) ([]string, error) {
	return []string{}, nil
}

func (m *MockGitClient) DeleteTag(name string) error {
	return nil
}

func (m *MockGitClient) GetTagCommit(tagName string) (string, error) {
	return "mock-commit-hash", nil
}

func (m *MockGitClient) SoftReset(commitHash string) error {
	return nil
}

func (m *MockGitClient) GetCommitMessage(commitHash string) (string, error) {
	return m.commitMessage, m.commitError
}

func (m *MockGitClient) IsTracked(path string) (bool, error) {
	return true, nil
}

func (m *MockGitClient) GetCommitParent(commitHash string) (string, error) {
	return "mock-parent-hash", nil
}

func (m *MockGitClient) GetUnstagedFiles() ([]string, error) {
	return []string{}, nil
}

func (m *MockGitClient) GetUntrackedFiles() ([]string, error) {
	return []string{}, nil
}
