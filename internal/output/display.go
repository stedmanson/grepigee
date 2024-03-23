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
	tbl := table.New("Folder Name", "File Name", "Line Number", "Match Text").WithWriter(os.Stdout)

	for _, item := range foundItems {
		if item.Error != nil {
			// Skip entries with errors
			continue
		}

		// Add a row for each item, respecting the specified order
		tbl.AddRow(item.FolderName, item.FileName, item.LineNum, item.MatchText)
	}

	// Print the table
	tbl.Print()
	fmt.Println()
}
