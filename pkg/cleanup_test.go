/*
Copyright © 2021-2026 Siwei Luo <siwei@lu0.org>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package pkg

import (
	"os"
	"path/filepath"
	"testing"
)

func initCleanupTestDirectory(t *testing.T) string {
	dir, err := os.MkdirTemp("", "brot-cleanup-tests-")
	if err != nil {
		t.Errorf("error - creating temporary working directory for tests at: %q", dir)
	}

	createTestDir(t, filepath.Join(dir, "src"))

	createTestFile(t, filepath.Join(dir, "src", "file_1.txt"))
	createTestFile(t, filepath.Join(dir, "src", "file_2.txt"))
	createTestFile(t, filepath.Join(dir, "src", "file_3.txt"))
	createTestFile(t, filepath.Join(dir, "src", "keep_this.txt"))

	return dir
}

func setupCleanupConfig(srcDir string, patterns []string) {
	CurrentConfiguration.Cleanup = make([]struct {
		Name     string   `mapstructure:"name"`
		Source   string   `mapstructure:"src"`
		Patterns []string `mapstructure:"patterns"`
	}, 1)
	CurrentConfiguration.Cleanup[0].Name = "test-cleanup"
	CurrentConfiguration.Cleanup[0].Source = srcDir
	CurrentConfiguration.Cleanup[0].Patterns = patterns
}

func TestCleanupBasicRemoval(t *testing.T) {
	testDir := initCleanupTestDirectory(t)
	defer os.RemoveAll(testDir)

	srcDir := filepath.Join(testDir, "src")

	setupCleanupConfig(srcDir, []string{"file_*.txt"})
	Cleanup(false)

	// verify matching files are removed
	for i := 1; i <= 3; i++ {
		filename := "file_" + string(rune(48+i)) + ".txt"
		file := filepath.Join(srcDir, filename)
		if _, err := os.Stat(file); err == nil {
			t.Errorf("failed - file should be removed: %q", file)
		}
	}

	// verify non-matching files still exist
	keepFile := filepath.Join(srcDir, "keep_this.txt")
	if _, err := os.Stat(keepFile); err != nil {
		t.Errorf("failed - non-matching file should still exist: %q", keepFile)
	}
}

func TestCleanupDryRun(t *testing.T) {
	testDir := initCleanupTestDirectory(t)
	defer os.RemoveAll(testDir)

	srcDir := filepath.Join(testDir, "src")

	setupCleanupConfig(srcDir, []string{"file_*.txt"})
	Cleanup(true) // dryRun = true

	// verify all files still exist (no actual removal happened)
	for i := 1; i <= 3; i++ {
		filename := "file_" + string(rune(48+i)) + ".txt"
		file := filepath.Join(srcDir, filename)
		if _, err := os.Stat(file); err != nil {
			t.Errorf("failed - file should still exist in dryRun mode: %q", file)
		}
	}

	keepFile := filepath.Join(srcDir, "keep_this.txt")
	if _, err := os.Stat(keepFile); err != nil {
		t.Errorf("failed - file should still exist in dryRun mode: %q", keepFile)
	}
}

func TestCleanupMissingSourceDirectory(t *testing.T) {
	testDir := initCleanupTestDirectory(t)
	defer os.RemoveAll(testDir)

	srcDir := filepath.Join(testDir, "nonexistent")

	setupCleanupConfig(srcDir, []string{"file_*.txt"})
	Cleanup(false)

	// No files to check - should just skip gracefully
	// Test passes if no panic occurs
}

func TestCleanupNoMatchingFiles(t *testing.T) {
	testDir := initCleanupTestDirectory(t)
	defer os.RemoveAll(testDir)

	srcDir := filepath.Join(testDir, "src")

	setupCleanupConfig(srcDir, []string{"*.pdf"})
	Cleanup(false)

	// verify all files still exist (no pattern match)
	for i := 1; i <= 3; i++ {
		filename := "file_" + string(rune(48+i)) + ".txt"
		file := filepath.Join(srcDir, filename)
		if _, err := os.Stat(file); err != nil {
			t.Errorf("failed - file should still exist when pattern doesn't match: %q", file)
		}
	}

	keepFile := filepath.Join(srcDir, "keep_this.txt")
	if _, err := os.Stat(keepFile); err != nil {
		t.Errorf("failed - file should still exist: %q", keepFile)
	}
}

func TestCleanupMultipleFiles(t *testing.T) {
	testDir := initCleanupTestDirectory(t)
	defer os.RemoveAll(testDir)

	srcDir := filepath.Join(testDir, "src")

	setupCleanupConfig(srcDir, []string{"file_*.txt"})
	Cleanup(false)

	// verify all matching files are removed
	for i := 1; i <= 3; i++ {
		filename := "file_" + string(rune(48+i)) + ".txt"
		file := filepath.Join(srcDir, filename)
		if _, err := os.Stat(file); err == nil {
			t.Errorf("failed - file should be removed: %q", file)
		}
	}
}

func TestCleanupEnvironmentVariables(t *testing.T) {
	testDir := initCleanupTestDirectory(t)
	defer os.RemoveAll(testDir)

	srcDir := filepath.Join(testDir, "src")

	// set environment variable
	testEnvVar := "BROT_TEST_CLEANUP_DIR"
	os.Setenv(testEnvVar, srcDir)
	defer os.Unsetenv(testEnvVar)

	// setup configuration with environment variable in path
	CurrentConfiguration.Cleanup = make([]struct {
		Name     string   `mapstructure:"name"`
		Source   string   `mapstructure:"src"`
		Patterns []string `mapstructure:"patterns"`
	}, 1)
	CurrentConfiguration.Cleanup[0].Name = "test-env-vars"
	CurrentConfiguration.Cleanup[0].Source = "$" + testEnvVar
	CurrentConfiguration.Cleanup[0].Patterns = []string{"file_1.txt"}

	Cleanup(false)

	// verify file was removed (environment variable was expanded correctly)
	file := filepath.Join(srcDir, "file_1.txt")
	if _, err := os.Stat(file); err == nil {
		t.Errorf("failed - environment variable was not expanded correctly: %q", file)
	}

	// verify other files still exist
	for i := 2; i <= 3; i++ {
		filename := "file_" + string(rune(48+i)) + ".txt"
		file := filepath.Join(srcDir, filename)
		if _, err := os.Stat(file); err != nil {
			t.Errorf("failed - non-matching file should still exist: %q", file)
		}
	}
}

func TestCleanupSpecificFile(t *testing.T) {
	testDir := initCleanupTestDirectory(t)
	defer os.RemoveAll(testDir)

	srcDir := filepath.Join(testDir, "src")

	// only remove file_1.txt
	setupCleanupConfig(srcDir, []string{"file_1.txt"})
	Cleanup(false)

	// verify file_1.txt is removed
	file := filepath.Join(srcDir, "file_1.txt")
	if _, err := os.Stat(file); err == nil {
		t.Errorf("failed - file should be removed: %q", file)
	}

	// verify other files still exist
	file2 := filepath.Join(srcDir, "file_2.txt")
	if _, err := os.Stat(file2); err != nil {
		t.Errorf("failed - file should still exist: %q", file2)
	}

	file3 := filepath.Join(srcDir, "file_3.txt")
	if _, err := os.Stat(file3); err != nil {
		t.Errorf("failed - file should still exist: %q", file3)
	}

	keepFile := filepath.Join(srcDir, "keep_this.txt")
	if _, err := os.Stat(keepFile); err != nil {
		t.Errorf("failed - file should still exist: %q", keepFile)
	}
}

func TestCleanupMultiplePatterns(t *testing.T) {
	testDir := initCleanupTestDirectory(t)
	defer os.RemoveAll(testDir)

	srcDir := filepath.Join(testDir, "src")

	// setup configuration with multiple patterns
	CurrentConfiguration.Cleanup = make([]struct {
		Name     string   `mapstructure:"name"`
		Source   string   `mapstructure:"src"`
		Patterns []string `mapstructure:"patterns"`
	}, 1)
	CurrentConfiguration.Cleanup[0].Name = "test-multi-pattern"
	CurrentConfiguration.Cleanup[0].Source = srcDir
	CurrentConfiguration.Cleanup[0].Patterns = []string{"file_1.txt", "file_2.txt", "keep_this.txt"}

	Cleanup(false)

	// verify file_1.txt and file_2.txt are removed
	file1 := filepath.Join(srcDir, "file_1.txt")
	if _, err := os.Stat(file1); err == nil {
		t.Errorf("failed - file should be removed: %q", file1)
	}

	file2 := filepath.Join(srcDir, "file_2.txt")
	if _, err := os.Stat(file2); err == nil {
		t.Errorf("failed - file should be removed: %q", file2)
	}

	keepFile := filepath.Join(srcDir, "keep_this.txt")
	if _, err := os.Stat(keepFile); err == nil {
		t.Errorf("failed - file should be removed: %q", keepFile)
	}

	// verify file_3.txt still exists (not in any pattern)
	file3 := filepath.Join(srcDir, "file_3.txt")
	if _, err := os.Stat(file3); err != nil {
		t.Errorf("failed - file should still exist: %q", file3)
	}
}

func TestCleanupEmptyDirectory(t *testing.T) {
	testDir, err := os.MkdirTemp("", "brot-cleanup-empty-")
	if err != nil {
		t.Errorf("error - creating temporary directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	srcDir := filepath.Join(testDir, "src")
	os.Mkdir(srcDir, 0755)

	setupCleanupConfig(srcDir, []string{"*.txt"})
	Cleanup(false)

	// Test passes if no panic occurs and directory is still empty
	files, err := os.ReadDir(srcDir)
	if err != nil {
		t.Errorf("error - reading directory: %v", err)
	}
	if len(files) != 0 {
		t.Errorf("failed - directory should be empty, found %d files", len(files))
	}
}

