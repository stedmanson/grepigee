package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func cleanupDirectory(directory string) {
	err := os.RemoveAll(directory)
	if err != nil {
		fmt.Printf("Error removing directory %s: %v\n", directory, err)
	}
}

func removeZipFiles(directory string) {
	// Find and remove zip files
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".zip") {
			err := os.Remove(path)
			if err != nil {
				fmt.Printf("Error removing file %s: %v\n", path, err)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking through directory %s: %v\n", directory, err)
	}
}
