package output

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/stedmanson/grepigee/internal/searcher"
)

func SaveAsCSV(foundItems []searcher.Found, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Writing the header
	header := []string{"Folder Name", "Revision", "File Name", "Line Number", "Match Text"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("error writing header to CSV file: %v", err)
	}

	// Writing data
	for _, item := range foundItems {
		if item.Error != nil {
			// Skip entries with errors
			continue
		}

		name, revision, err := extractNameAndRevision(item.FolderName)
		if err != nil {
			fmt.Printf("Error extracting name and revision: %v\n", err)
			continue
		}

		record := []string{
			name,
			revision,
			item.FileName,
			fmt.Sprintf("%d", item.LineNum),
			item.MatchText,
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("error writing record to CSV file: %v", err)
		}
	}

	return nil
}
