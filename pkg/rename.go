package pkg

import (
	"fmt"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func RenameFiles(dir string, pattern string, prefixLength int) error {
	if prefixLength < 1 {
		return fmt.Errorf("prefix length must be at least 1, got %d", prefixLength)
	}

	files := FilesFromDirectory(dir, []string{pattern})
	if len(files) == 0 {
		log.Info("no files found matching pattern")
		return nil
	}

	for i, path := range files {
		checksum, err := FileChecksum(path)
		if err != nil {
			log.WithFields(log.Fields{
				"file":  path,
				"error": err,
			}).Error("error calculating checksum")
			continue
		}

		ext := filepath.Ext(path)
		newName := fmt.Sprintf("%0*d-%s%s", prefixLength, i+1, checksum, ext)
		newPath := filepath.Join(filepath.Dir(path), newName)

		if path == newPath {
			continue
		}

		if err := FileMove(path, newPath); err != nil {
			log.WithFields(log.Fields{
				"src":   path,
				"dst":   newPath,
				"error": err,
			}).Error("error renaming file")
			continue
		}

		log.WithFields(log.Fields{
			"src": path,
			"dst": newPath,
		}).Info("renamed file")
	}
	return nil
}
