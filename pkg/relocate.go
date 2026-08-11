package pkg

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Relocate(dryRun bool) {
	// log the used configuration file
	log.Info("use config file: ", viper.ConfigFileUsed())

	// iterate over relocate definitions from configuration
	for _, item := range CurrentConfiguration.Relocate {

		// expand any environment variables
		srcDirectory := os.ExpandEnv(item.Source)
		dstDirectory := os.ExpandEnv(item.Destination)

		// get files from the source directory
		relocateFiles := FilesFromDirectory(srcDirectory, item.Patterns)

		// skip the item if there are no files to relocate
		if relocateFiles == nil {
			continue
		}

		// check if the destination directory exists and skip if it is missing
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
			if _, err := os.Stat(dstPath); err == nil {
				log.WithFields(log.Fields{
					"src":  srcPath,
					"dst":  dstDirectory,
					"mode": item.Mode,
				}).Warnf("skip file: %v", srcFile)
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
			}).Infof("%v file: %v", item.Mode, srcFile)

		}
	}
}
