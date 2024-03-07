/*
Copyright © 2021-2023 Siwei Luo <siwei@lu0.org>

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
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// helper function to create files with dummy content
func createTestFile(t *testing.T, file string) {
	err := os.WriteFile(file, []byte("TESTDATA"), 0644)
	if err != nil {
		t.Errorf("error - creating temporary file at: %q", file)
	}
}

// helper function to create directories
func createTestDir(t *testing.T, dir string) {
	err := os.Mkdir(dir, 0755)
	if err != nil {
		t.Errorf("error - creating temporary direcory at: %q", dir)
	}
}

func initTestDirectory(t *testing.T) string {
	dir, err := os.MkdirTemp("", "brot-unittests-")
	if err != nil {
		t.Errorf("error - creating temporary working directory for tests at: %q", dir)
	}

	createTestDir(t, filepath.Join(dir, "src"))
	createTestDir(t, filepath.Join(dir, "dst"))

	createTestFile(t, filepath.Join(dir, "src", "test_1.txt"))
	createTestFile(t, filepath.Join(dir, "src", "test_2.txt"))
	createTestFile(t, filepath.Join(dir, "src", "test_3.txt"))
	createTestFile(t, filepath.Join(dir, "src", "test_4.txt"))
	createTestFile(t, filepath.Join(dir, "src", "_test_5.txt"))

	return dir
}

func TestFilesFromDirectory(t *testing.T) {
	// create a temporary directory with test data and remove everything afterward again
	testDir := initTestDirectory(t)
	defer func() {
		if err := os.RemoveAll(testDir); err != nil {
			t.Errorf("error - removing test directory at: %q", testDir)
		}
	}()

	// match single file
	testResult1 := FilesFromDirectory(testDir, []string{"test_1.txt"})
	if len(testResult1) != 1 {
		t.Errorf("failed - got %q but expected %q files", testResult1, 1)
	}

	// match multiple files with pattern;
	testResult2 := FilesFromDirectory(testDir, []string{"test_*.txt"})
	if len(testResult2) != 4 {
		t.Errorf("failed - got %q but expected %q files", testResult2, 4)
	}

	// match multiple files with pattern; only file ending
	testResult3 := FilesFromDirectory(testDir, []string{"*.txt"})
	if len(testResult3) != 5 {
		t.Errorf("failed - got %q but expected %q files", testResult3, 5)
	}

	// no match
	testResult4 := FilesFromDirectory(testDir, []string{"*.foo"})
	if len(testResult4) != 0 {
		t.Errorf("failed - got %q but expected %q files", testResult4, 0)
	}
}

func TestFileCopy(t *testing.T) {
	// create a temporary directory with test data and remove everything afterward again
	testDir := initTestDirectory(t)
	defer func() {
		if err := os.RemoveAll(testDir); err != nil {
			t.Errorf("error - removing test directory at: %q", testDir)
		}
	}()

	src := filepath.Join(testDir, "src", "test_1.txt")
	dst := filepath.Join(testDir, "dst", "test_1.txt")

	err := FileCopy(src, dst)
	if err != nil {
		t.Errorf("error - copying %q to %q", src, dst)
	}

	// check if source file is still present
	if _, err := os.Stat(src); errors.Is(err, os.ErrNotExist) {
		t.Errorf("failed - source file missing: %q", src)
	}

	// check if destination file is present
	if _, err := os.Stat(dst); errors.Is(err, os.ErrNotExist) {
		t.Errorf("failed - destination file missing: %q", dst)
	}
}

func TestFileMove(t *testing.T) {
	// create a temporary directory with test data and remove everything afterward again
	testDir := initTestDirectory(t)
	defer func() {
		if err := os.RemoveAll(testDir); err != nil {
			t.Errorf("error - removing test directory at: %q", testDir)
		}
	}()

	src := filepath.Join(testDir, "src", "test_1.txt")
	dst := filepath.Join(testDir, "dst", "test_1.txt")

	err := FileMove(src, dst)
	if err != nil {
		t.Errorf("error - copying %q to %q", src, dst)
	}

	// check if source file is absent
	if _, err := os.Stat(src); err == nil {
		t.Errorf("failed - source was not moved: %q", src)
	}

	// check if destination file is present
	if _, err := os.Stat(dst); errors.Is(err, os.ErrNotExist) {
		t.Errorf("failed - destination file missing: %q", dst)
	}
}

func TestFileRemove(t *testing.T) {
	// create a temporary directory with test data and remove everything afterward again
	testDir := initTestDirectory(t)
	defer func() {
		if err := os.RemoveAll(testDir); err != nil {
			t.Errorf("error - removing test directory at: %q", testDir)
		}
	}()

	file := filepath.Join(testDir, "src", "test_1.txt")

	err := FileRemove(file)
	if err != nil {
		t.Errorf("error - deleting file: %q", file)
	}

	// check if source file is absent
	if _, err := os.Stat(file); err == nil {
		t.Errorf("failed - file was not deleted: %q", file)
	}
}
