package output

import (
	"fmt"
	"os"

	"github.com/stedmanson/grepigee/internal/searcher"

	"github.com/rodaine/table"
)

func DisplayAsTable(foundItems []searcher.Found) {
	if len(foundItems) == 0 {
		fmt.Println("No items found.")
		return
	}

	// Initialize the table with headers in the desired order
	tbl := table.New("Folder Name", "Revision", "File Name", "Line Number", "Match Text").WithWriter(os.Stdout)

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

		// Add a row for each item, respecting the specified order
		tbl.AddRow(name, revision, item.FileName, item.LineNum, item.MatchText)
	}

	// Print the table
	tbl.Print()
	fmt.Println()
	fmt.Println("Found", len(foundItems), "items.")
}
