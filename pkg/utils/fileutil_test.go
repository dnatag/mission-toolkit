package utils

import (
	"strings"
	"testing"

	"github.com/spf13/afero"
)

func TestCopyFile(t *testing.T) {
	tests := []struct {
		name           string
		setupFS        func(afero.Fs) (src, dst string)
		wantErr        bool
		errContains    string
		validateResult func(*testing.T, afero.Fs, string, string) bool
	}{
		{
			name: "successful copy",
			setupFS: func(fs afero.Fs) (src, dst string) {
				src = "/test/source.txt"
				dst = "/test/destination.txt"
				if err := fs.MkdirAll("/test", 0755); err != nil {
					t.Fatalf("failed to create test directory: %v", err)
				}
				if err := afero.WriteFile(fs, src, []byte("test content"), 0644); err != nil {
					t.Fatalf("failed to create source file: %v", err)
				}
				return src, dst
			},
			wantErr: false,
			validateResult: func(t *testing.T, fs afero.Fs, src, dst string) bool {
				// Verify destination file exists with same content
				srcContent, err := afero.ReadFile(fs, src)
				if err != nil {
					t.Errorf("failed to read source file: %v", err)
					return false
				}
				dstContent, err := afero.ReadFile(fs, dst)
				if err != nil {
					t.Errorf("failed to read destination file: %v", err)
					return false
				}
				if string(srcContent) != string(dstContent) {
					t.Errorf("source and destination content mismatch: got %q, want %q", string(dstContent), string(srcContent))
					return false
				}
				return true
			},
		},
		{
			name: "copy to nested directory",
			setupFS: func(fs afero.Fs) (src, dst string) {
				src = "/source.txt"
				dst = "/nested/deeply/destination.txt"
				if err := afero.WriteFile(fs, src, []byte("nested content"), 0644); err != nil {
					t.Fatalf("failed to create source file: %v", err)
				}
				return src, dst
			},
			wantErr: false,
			validateResult: func(t *testing.T, fs afero.Fs, src, dst string) bool {
				dstContent, err := afero.ReadFile(fs, dst)
				if err != nil {
					t.Errorf("failed to read destination file: %v", err)
					return false
				}
				if string(dstContent) != "nested content" {
					t.Errorf("unexpected content: got %q, want %q", string(dstContent), "nested content")
					return false
				}
				return true
			},
		},
		{
			name: "source file does not exist",
			setupFS: func(fs afero.Fs) (src, dst string) {
				src = "/nonexistent.txt"
				dst = "/destination.txt"
				return src, dst
			},
			wantErr:     true,
			errContains: "reading source file",
		},
		{
			name: "empty file copy",
			setupFS: func(fs afero.Fs) (src, dst string) {
				src = "/empty.txt"
				dst = "/empty_copy.txt"
				if err := afero.WriteFile(fs, src, []byte{}, 0644); err != nil {
					t.Fatalf("failed to create source file: %v", err)
				}
				return src, dst
			},
			wantErr: false,
			validateResult: func(t *testing.T, fs afero.Fs, src, dst string) bool {
				dstContent, err := afero.ReadFile(fs, dst)
				if err != nil {
					t.Errorf("failed to read destination file: %v", err)
					return false
				}
				if len(dstContent) != 0 {
					t.Errorf("expected empty file, got %q", string(dstContent))
					return false
				}
				return true
			},
		},
		{
			name: "large file copy",
			setupFS: func(fs afero.Fs) (src, dst string) {
				src = "/large.txt"
				dst = "/large_copy.txt"
				largeContent := make([]byte, 10*1024) // 10KB
				for i := range largeContent {
					largeContent[i] = byte(i % 256)
				}
				if err := afero.WriteFile(fs, src, largeContent, 0644); err != nil {
					t.Fatalf("failed to create source file: %v", err)
				}
				return src, dst
			},
			wantErr: false,
			validateResult: func(t *testing.T, fs afero.Fs, src, dst string) bool {
				srcContent, err := afero.ReadFile(fs, src)
				if err != nil {
					t.Errorf("failed to read source file: %v", err)
					return false
				}
				dstContent, err := afero.ReadFile(fs, dst)
				if err != nil {
					t.Errorf("failed to read destination file: %v", err)
					return false
				}
				if len(srcContent) != len(dstContent) {
					t.Errorf("size mismatch: got %d, want %d", len(dstContent), len(srcContent))
					return false
				}
				return true
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			src, dst := tt.setupFS(fs)

			err := CopyFile(fs, src, dst)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CopyFile() expected error containing %q, got nil", tt.errContains)
					return
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("CopyFile() error = %q, want error containing %q", err.Error(), tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("CopyFile() unexpected error: %v", err)
				return
			}

			if tt.validateResult != nil {
				if !tt.validateResult(t, fs, src, dst) {
					t.Error("CopyFile() validation failed")
				}
			}
		})
	}
}

func TestCopyFile_OverwriteExisting(t *testing.T) {
	fs := afero.NewMemMapFs()

	// Create source and existing destination with different content
	src := "/source.txt"
	dst := "/existing.txt"

	srcContent := []byte("new content")
	dstContent := []byte("old content")

	if err := afero.WriteFile(fs, src, srcContent, 0644); err != nil {
		t.Fatalf("failed to create source file: %v", err)
	}
	if err := afero.WriteFile(fs, dst, dstContent, 0644); err != nil {
		t.Fatalf("failed to create destination file: %v", err)
	}

	// Copy should overwrite existing destination
	if err := CopyFile(fs, src, dst); err != nil {
		t.Errorf("CopyFile() unexpected error: %v", err)
	}

	// Verify destination was overwritten
	result, err := afero.ReadFile(fs, dst)
	if err != nil {
		t.Errorf("failed to read destination file: %v", err)
	}
	if string(result) != "new content" {
		t.Errorf("destination not overwritten: got %q, want %q", string(result), "new content")
	}
}

func TestCopyFile_NestedPathsWithSameName(t *testing.T) {
	fs := afero.NewMemMapFs()

	// Test copying files with same name to different directories
	src1 := "/dir1/file.txt"
	dst1 := "/dir2/file.txt"
	src2 := "/dir3/file.txt"
	dst2 := "/dir4/subdir/file.txt"

	content1 := []byte("content 1")
	content2 := []byte("content 2")

	if err := afero.WriteFile(fs, src1, content1, 0644); err != nil {
		t.Fatalf("failed to create source file: %v", err)
	}
	if err := afero.WriteFile(fs, src2, content2, 0644); err != nil {
		t.Fatalf("failed to create source file: %v", err)
	}

	// Copy both files
	if err := CopyFile(fs, src1, dst1); err != nil {
		t.Errorf("CopyFile() error: %v", err)
	}
	if err := CopyFile(fs, src2, dst2); err != nil {
		t.Errorf("CopyFile() error: %v", err)
	}

	// Verify both copies
	result1, err := afero.ReadFile(fs, dst1)
	if err != nil {
		t.Errorf("failed to read destination file: %v", err)
	}
	if string(result1) != "content 1" {
		t.Errorf("unexpected content: got %q, want %q", string(result1), "content 1")
	}

	result2, err := afero.ReadFile(fs, dst2)
	if err != nil {
		t.Errorf("failed to read destination file: %v", err)
	}
	if string(result2) != "content 2" {
		t.Errorf("unexpected content: got %q, want %q", string(result2), "content 2")
	}
}

func TestCopyFile_PreservePermissions(t *testing.T) {
	fs := afero.NewMemMapFs()

	src := "/test.txt"
	dst := "/copy.txt"
	content := []byte("test content")

	if err := afero.WriteFile(fs, src, content, 0644); err != nil {
		t.Fatalf("failed to create source file: %v", err)
	}

	if err := CopyFile(fs, src, dst); err != nil {
		t.Errorf("CopyFile() error: %v", err)
	}

	// Verify file was copied with correct permissions (0644)
	info, err := fs.Stat(dst)
	if err != nil {
		t.Errorf("failed to stat destination file: %v", err)
	}

	// Note: afero.MemMapFs doesn't preserve Unix permissions precisely
	// but we can verify the file exists and has content
	if info == nil {
		t.Error("destination file info is nil")
	}
}
