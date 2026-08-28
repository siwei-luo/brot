package pkg

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	log "github.com/sirupsen/logrus"
)

// Version `[[VERSION]]` is replaced during pipeline build with the respective string
const Version = "🍞 1.0.0"
const VersionMajor = 1

// struct representing the configuration file
type configuration struct {
	ApiVersion string `mapstructure:"apiVersion"`
	Defaults   struct {
		Loglevel  string `mapstructure:"loglevel"`
		Logformat string `mapstructure:"logformat"`
	} `mapstructure:"defaults"`
	Relocate []struct {
		Name        string   `mapstructure:"name"`
		Source      string   `mapstructure:"src"`
		Destination string   `mapstructure:"dst"`
		Patterns    []string `mapstructure:"patterns"`
		Mode        string   `mapstructure:"mode"`
	} `mapstructure:"relocate"`
	Cleanup []struct {
		Name     string   `mapstructure:"name"`
		Source   string   `mapstructure:"src"`
		Patterns []string `mapstructure:"patterns"`
	} `mapstructure:"cleanup"`
}

var CurrentConfiguration configuration

var Verbosity int

func FilesFromDirectory(directory string, patterns []string) []string {
	var files []string

	if err := filepath.Walk(directory, visit(patterns, &files)); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("error reading directory")

		return nil
	}

	sort.Strings(files)

	log.WithFields(log.Fields{
		"files": files,
	}).Debug("found files")

	return files
}

func FileCopy(src string, dst string) (err error) {

	if _, err := os.Stat(src); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("source file not found: %w", err)
	}

	if _, err := os.Stat(dst); err == nil {
		return fmt.Errorf("destination file already exists: %s", dst)
	}

	in, err := os.Open(filepath.Clean(src))
	if err != nil {
		return err
	}
	defer func(in *os.File) {
		_ = in.Close()
	}(in)

	out, err := os.Create(filepath.Clean(dst))
	if err != nil {
		return err
	}
	defer func(out *os.File) {
		_ = out.Close()
	}(out)

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func FileMove(src string, dst string) (err error) {
	if _, err := os.Stat(dst); err == nil {
		return fmt.Errorf("destination file already exists: %s", dst)
	} else if !errors.Is(err, os.ErrNotExist) {
		// some other stat error occurred
		return err
	}
	// destination does not exist, safe to rename
	return os.Rename(src, dst)
}

func FileRemove(src string) (err error) {
	if _, err := os.Stat(src); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("file not found: %w", err)
	} else if err != nil {
		return err
	}
	return os.Remove(src)
}

func FileChecksum(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func visit(patterns []string, files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Warn("skip reading directory")
			return nil
		}

		// iterate over all patterns
		for _, pattern := range patterns {

			// skip if pattern is empty and do not match any files
			if pattern == "" {
				continue
			}

			matched, err := filepath.Match(pattern, info.Name())
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Fatal("error in matching pattern")
			}
			if matched {
				*files = append(*files, path)

				log.WithFields(log.Fields{
					"file": path,
				}).Debug("matched file")
				// only process one pattern match per file
				break
			}
		}

		return nil
	}
}
