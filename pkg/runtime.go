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
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
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
		}).Fatal("error reading directory")

		return nil
	}

	log.WithFields(log.Fields{
		"files": files,
	}).Debug("found files")

	return files
}

func FileCopy(src string, dst string) (err error) {

	// abort when source file is missing
	if _, err := os.Stat(src); errors.Is(err, os.ErrNotExist) {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("error reading file")
		return nil
	}

	// abort if a file with the same name exists in the destination
	if _, err := os.Stat(dst); err == nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("file already present in destination")
		return nil
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
	if _, err := os.Stat(dst); errors.Is(err, os.ErrNotExist) {
		return os.Rename(src, dst)
	}
	return err
}

func FileRemove(src string) (err error) {
	if _, err := os.Stat(src); err == nil {
		return os.Remove(src)
	}
	return err
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
			} else {
				log.WithFields(log.Fields{
					"file": path,
				}).Debug("ignored file")
			}
		}

		return nil
	}
}
