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
)

func Cleanup(dryRun bool) {
	// log used configuration file
	log.Info("use config file: ", viper.ConfigFileUsed())

	// iterate over relocate definitions from configuration
	for _, item := range CurrentConfiguration.Cleanup {

		// expand any environment variables
		srcDirectory := os.ExpandEnv(item.Source)

		// get files from source directory
		cleanupFiles := FilesFromDirectory(srcDirectory, item.Patterns)

		// skip item if there are no files to relocate
		if cleanupFiles == nil {
			continue
		}

		for _, srcPath := range cleanupFiles {

			if !dryRun {
				if err := FileRemove(srcPath); err != nil {
					log.WithFields(log.Fields{
						"error": err,
						"src":   srcPath,
					}).Error("error removing file")
				}
			}

			log.WithFields(log.Fields{
				"src": srcPath,
			}).Infof("remove file: %v\n", srcPath)

		}
	}
}
