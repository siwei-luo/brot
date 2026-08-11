package pkg

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Cleanup(dryRun bool) {
	// log used configuration file
	log.Info("use config file: ", viper.ConfigFileUsed())

	// iterate over cleanup definitions from configuration
	for _, item := range CurrentConfiguration.Cleanup {

		// expand any environment variables
		srcDirectory := os.ExpandEnv(item.Source)

		// get files from source directory
		cleanupFiles := FilesFromDirectory(srcDirectory, item.Patterns)

		// skip item if there are no files to cleanup
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
					continue
				}
			}

			log.WithFields(log.Fields{
				"src": srcPath,
			}).Infof("remove file: %v", srcPath)

		}
	}
}
