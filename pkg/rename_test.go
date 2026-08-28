package pkg

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestRenameFiles(t *testing.T) {
	tests := []struct {
		name          string
		filesToCreate map[string]string // filename -> content
		pattern       string
		prefixLength  int
		expectedFiles []string
		expectError   bool
	}{
		{
			name: "successful rename of multiple files",
			filesToCreate: map[string]string{
				"test1.jpg": "content1",
				"test2.jpg": "content2",
				"skip.txt":  "content3",
			},
			pattern: "*.jpg",
			expectedFiles: []string{
				fmt.Sprintf("00001-%s.jpg", fmt.Sprintf("%x", sha256.Sum256([]byte("content1")))),
				fmt.Sprintf("00002-%s.jpg", fmt.Sprintf("%x", sha256.Sum256([]byte("content2")))),
			},
		},
		{
			name: "no files matching pattern",
			filesToCreate: map[string]string{
				"test1.txt": "content1",
			},
			pattern:       "*.jpg",
			expectedFiles: []string{},
		},
		{
			name: "rename single file",
			filesToCreate: map[string]string{
				"image.png": "image-data",
			},
			pattern:      "*.png",
			prefixLength: 4,
			expectedFiles: []string{
				fmt.Sprintf("0001-%s.png", fmt.Sprintf("%x", sha256.Sum256([]byte("image-data")))),
			},
		},
		{
			name: "rename multiple files with custom prefix length",
			filesToCreate: map[string]string{
				"test1.jpg": "content1",
				"test2.jpg": "content2",
				"test3.jpg": "content3",
			},
			pattern:      "*.jpg",
			prefixLength: 4,
			expectedFiles: []string{
				fmt.Sprintf("0001-%s.jpg", fmt.Sprintf("%x", sha256.Sum256([]byte("content1")))),
				fmt.Sprintf("0002-%s.jpg", fmt.Sprintf("%x", sha256.Sum256([]byte("content2")))),
				fmt.Sprintf("0003-%s.jpg", fmt.Sprintf("%x", sha256.Sum256([]byte("content3")))),
			},
		},
		{
			name:          "invalid prefix length",
			filesToCreate: map[string]string{},
			pattern:       "*.jpg",
			prefixLength:  -1,
			expectedFiles: []string{},
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDir, err := os.MkdirTemp("", "brot-rename-tests-")
			if err != nil {
				t.Fatalf("failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(testDir)

			for name, content := range tt.filesToCreate {
				err := os.WriteFile(filepath.Join(testDir, name), []byte(content), 0644)
				if err != nil {
					t.Fatalf("failed to create test file %s: %v", name, err)
				}
			}

			prefixLength := tt.prefixLength
			if prefixLength == 0 {
				prefixLength = 5
			}

			err = RenameFiles(testDir, tt.pattern, prefixLength)
			if (err != nil) != tt.expectError {
				t.Errorf("RenameFiles() error = %v, expectError %v", err, tt.expectError)
				return
			}

			files, err := os.ReadDir(testDir)
			if err != nil {
				t.Fatalf("failed to read dir: %v", err)
			}

			var actualFiles []string
			for _, f := range files {
				actualFiles = append(actualFiles, f.Name())
			}

			// We only check if the expected renamed files exist in the directory
			for _, expected := range tt.expectedFiles {
				found := false
				for _, actual := range actualFiles {
					if actual == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected renamed file %s not found in directory", expected)
				}
			}

			// Ensure files that shouldn't be renamed still exist if they didn't match pattern
			for name := range tt.filesToCreate {
				isTarget := false
				// check if file matches the glob pattern
				matched, _ := filepath.Match(tt.pattern, name)
				if matched {
					isTarget = true
				}

				if !isTarget {
					if _, err := os.Stat(filepath.Join(testDir, name)); err != nil {
						t.Errorf("file %s should not have been renamed", name)
					}
				}
			}
		})
	}
}
