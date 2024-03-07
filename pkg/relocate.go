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
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

func Relocate(dryRun bool) {
	// log used configuration file
	log.Info("use config file: ", viper.ConfigFileUsed())

	// iterate over relocate definitions from configuration
	for _, item := range CurrentConfiguration.Relocate {

		// expand any environment variables
		srcDirectory := os.ExpandEnv(item.Source)
		dstDirectory := os.ExpandEnv(item.Destination)

		// get files from source directory
		relocateFiles := FilesFromDirectory(srcDirectory, item.Patterns)

		// skip item if there are no files to relocate
		if relocateFiles == nil {
			continue
		}

		// check if destination directory exists and skip if it is missing
		if _, err := os.Stat(dstDirectory); os.IsNotExist(err) {
			log.WithFields(log.Fields{
				"error": err,
			}).Warn("skip missing destination")
			continue
		}

		for _, srcPath := range relocateFiles {

			// assemble full destination path preserving the file's name
			srcFile := filepath.Base(srcPath)
			dstPath := filepath.Join(dstDirectory, srcFile)

			// check if a file with the same name exists in destination
			if _, err := os.Stat(dstPath); err == nil || os.IsExist(err) {
				log.WithFields(log.Fields{
					"src":  srcPath,
					"dst":  dstDirectory,
					"mode": item.Mode,
				}).Warnf("skip file: %v\n", srcFile)
				continue
			}

			switch item.Mode {
			case "move":
				if !dryRun {
					if err := FileMove(srcPath, dstPath); err != nil {
						log.WithFields(log.Fields{
							"error": err,
							"src":   srcPath,
							"dst":   dstDirectory,
						}).Error("error moving file")
					}
				}
			case "copy":
				if !dryRun {
					if err := FileCopy(srcPath, dstPath); err != nil {
						log.WithFields(log.Fields{
							"error": err,
							"src":   srcPath,
							"dst":   dstDirectory,
						}).Error("error copying file")
					}
				}
			}

			log.WithFields(log.Fields{
				"src":  srcPath,
				"dst":  dstDirectory,
				"mode": item.Mode,
			}).Infof("%v file: %v\n", item.Mode, srcFile)

		}
	}
}
