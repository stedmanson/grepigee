package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
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

func parseDateTime(dateTime string) string {
	dt, err := time.ParseInLocation("01/02/2006 15:04", dateTime, time.Local)
	if err != nil {
		fmt.Printf("Invalid format. Expected: MM/DD/YYYY HH:MM , Received: %s\n", dateTime)
		os.Exit(1)
	}
	loc, err := time.LoadLocation("UTC")
	if err != nil {
		os.Exit(1)
	}
	//convert to utc time
	return dt.In(loc).Format("01/02/2006 15:04")
}
