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

func initRelocateTestDirectory(t *testing.T) string {
	dir, err := os.MkdirTemp("", "brot-relocate-tests-")
	if err != nil {
		t.Errorf("error - creating temporary working directory for tests at: %q", dir)
	}

	createTestDir(t, filepath.Join(dir, "src"))
	createTestDir(t, filepath.Join(dir, "dst"))

	createTestFile(t, filepath.Join(dir, "src", "file_1.txt"))
	createTestFile(t, filepath.Join(dir, "src", "file_2.txt"))
	createTestFile(t, filepath.Join(dir, "src", "file_3.txt"))

	return dir
}

func setupRelocateConfig(srcDir, dstDir, mode string, patterns []string) {
	CurrentConfiguration.Relocate = make([]struct {
		Name        string   `mapstructure:"name"`
		Source      string   `mapstructure:"src"`
		Destination string   `mapstructure:"dst"`
		Patterns    []string `mapstructure:"patterns"`
		Mode        string   `mapstructure:"mode"`
	}, 1)
	CurrentConfiguration.Relocate[0].Name = "test-" + mode
	CurrentConfiguration.Relocate[0].Source = srcDir
	CurrentConfiguration.Relocate[0].Destination = dstDir
	CurrentConfiguration.Relocate[0].Patterns = patterns
	CurrentConfiguration.Relocate[0].Mode = mode
}

func TestRelocateCopyMode(t *testing.T) {
	testDir := initRelocateTestDirectory(t)
	defer os.RemoveAll(testDir)

	srcDir := filepath.Join(testDir, "src")
	dstDir := filepath.Join(testDir, "dst")

	setupRelocateConfig(srcDir, dstDir, "copy", []string{"file_*.txt"})
	Relocate(false)

	// verify source files still exist
	srcFile := filepath.Join(srcDir, "file_1.txt")
	if _, err := os.Stat(srcFile); err != nil {
		t.Errorf("failed - source file should still exist: %q", srcFile)
	}

	// verify destination files exist
	dstFile := filepath.Join(dstDir, "file_1.txt")
	if _, err := os.Stat(dstFile); err != nil {
		t.Errorf("failed - destination file missing: %q", dstFile)
	}
}

func TestRelocateMoveMode(t *testing.T) {
	testDir := initRelocateTestDirectory(t)
	defer os.RemoveAll(testDir)

	srcDir := filepath.Join(testDir, "src")
	dstDir := filepath.Join(testDir, "dst")

	setupRelocateConfig(srcDir, dstDir, "move", []string{"file_*.txt"})
	Relocate(false)

	// verify source files no longer exist
	srcFile := filepath.Join(srcDir, "file_1.txt")
	if _, err := os.Stat(srcFile); err == nil {
		t.Errorf("failed - source file should not exist: %q", srcFile)
	}

	// verify destination files exist
	dstFile := filepath.Join(dstDir, "file_1.txt")
	if _, err := os.Stat(dstFile); err != nil {
		t.Errorf("failed - destination file missing: %q", dstFile)
	}
}

func TestRelocateDryRun(t *testing.T) {
	testDir := initRelocateTestDirectory(t)
	defer os.RemoveAll(testDir)

	srcDir := filepath.Join(testDir, "src")
	dstDir := filepath.Join(testDir, "dst")

	setupRelocateConfig(srcDir, dstDir, "move", []string{"file_*.txt"})
	Relocate(true) // dryRun = true

	// verify source files still exist (no actual move happened)
	srcFile := filepath.Join(srcDir, "file_1.txt")
	if _, err := os.Stat(srcFile); err != nil {
		t.Errorf("failed - source file should still exist in dryRun mode: %q", srcFile)
	}

	// verify destination files do not exist
	dstFile := filepath.Join(dstDir, "file_1.txt")
	if _, err := os.Stat(dstFile); err == nil {
		t.Errorf("failed - destination file should not exist in dryRun mode: %q", dstFile)
	}
}

func TestRelocateMissingDestination(t *testing.T) {
	testDir := initRelocateTestDirectory(t)
	defer os.RemoveAll(testDir)

	srcDir := filepath.Join(testDir, "src")
	dstDir := filepath.Join(testDir, "dst")

	// remove the destination directory to test the skip logic
	if err := os.RemoveAll(dstDir); err != nil {
		t.Errorf("error - removing destination directory: %v", err)
	}

	setupRelocateConfig(srcDir, dstDir, "copy", []string{"file_*.txt"})
	Relocate(false)

	// verify source files still exist (no operation occurred)
	srcFile := filepath.Join(srcDir, "file_1.txt")
	if _, err := os.Stat(srcFile); err != nil {
		t.Errorf("failed - source file should still exist when destination is missing: %q", srcFile)
	}
}

func TestRelocateSkipExistingFile(t *testing.T) {
	testDir := initRelocateTestDirectory(t)
	defer os.RemoveAll(testDir)

	srcDir := filepath.Join(testDir, "src")
	dstDir := filepath.Join(testDir, "dst")

	// create a file in destination with same name to test skip logic
	createTestFile(t, filepath.Join(dstDir, "file_1.txt"))

	setupRelocateConfig(srcDir, dstDir, "move", []string{"file_1.txt"})
	Relocate(false)

	// verify source file still exists (skipped due to existing destination file)
	srcFile := filepath.Join(srcDir, "file_1.txt")
	if _, err := os.Stat(srcFile); err != nil {
		t.Errorf("failed - source file should exist when destination file already exists: %q", srcFile)
	}
}

func TestRelocateNoMatchingFiles(t *testing.T) {
	testDir := initRelocateTestDirectory(t)
	defer os.RemoveAll(testDir)

	srcDir := filepath.Join(testDir, "src")
	dstDir := filepath.Join(testDir, "dst")

	setupRelocateConfig(srcDir, dstDir, "copy", []string{"*.pdf"})
	Relocate(false)

	// verify source files still exist (no operation occurred)
	srcFile := filepath.Join(srcDir, "file_1.txt")
	if _, err := os.Stat(srcFile); err != nil {
		t.Errorf("failed - source file should exist when no files match pattern: %q", srcFile)
	}

	// verify destination is empty
	dstFiles, err := os.ReadDir(dstDir)
	if err != nil {
		t.Errorf("error - reading destination directory: %v", err)
	}
	if len(dstFiles) != 0 {
		t.Errorf("failed - destination directory should be empty, found %d files", len(dstFiles))
	}
}

func TestRelocateMultipleFiles(t *testing.T) {
	testDir := initRelocateTestDirectory(t)
	defer os.RemoveAll(testDir)

	srcDir := filepath.Join(testDir, "src")
	dstDir := filepath.Join(testDir, "dst")

	setupRelocateConfig(srcDir, dstDir, "copy", []string{"file_*.txt"})
	Relocate(false)

	// verify all three files were copied
	for i := 1; i <= 3; i++ {
		filename := "file_" + string(rune(48+i)) + ".txt"
		dstFile := filepath.Join(dstDir, filename)
		if _, err := os.Stat(dstFile); err != nil {
			t.Errorf("failed - destination file missing: %q", dstFile)
		}
	}
}

func TestRelocateEnvironmentVariables(t *testing.T) {
	testDir := initRelocateTestDirectory(t)
	defer os.RemoveAll(testDir)

	srcDir := filepath.Join(testDir, "src")
	dstDir := filepath.Join(testDir, "dst")

	// set environment variable
	testEnvVar := "BROT_TEST_SRC_DIR"
	os.Setenv(testEnvVar, srcDir)
	defer os.Unsetenv(testEnvVar)

	// setup configuration with environment variable in path
	CurrentConfiguration.Relocate = make([]struct {
		Name        string   `mapstructure:"name"`
		Source      string   `mapstructure:"src"`
		Destination string   `mapstructure:"dst"`
		Patterns    []string `mapstructure:"patterns"`
		Mode        string   `mapstructure:"mode"`
	}, 1)
	CurrentConfiguration.Relocate[0].Name = "test-env-vars"
	CurrentConfiguration.Relocate[0].Source = "$" + testEnvVar
	CurrentConfiguration.Relocate[0].Destination = dstDir
	CurrentConfiguration.Relocate[0].Patterns = []string{"file_1.txt"}
	CurrentConfiguration.Relocate[0].Mode = "copy"

	Relocate(false)

	// verify file was copied (environment variable was expanded correctly)
	dstFile := filepath.Join(dstDir, "file_1.txt")
	if _, err := os.Stat(dstFile); err != nil {
		t.Errorf("failed - environment variable was not expanded correctly: %q", dstFile)
	}
}





