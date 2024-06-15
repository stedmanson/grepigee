package output

import (
	"fmt"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/stedmanson/grepigee/internal/searcher"
)

func DisplayAsTable(foundItems []searcher.Found) {
	if len(foundItems) == 0 {
		fmt.Println("No items found.")
		return
	}

	// Initialize the tablewriter with os.Stdout as the output
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Folder Name", "Revision", "File Name", "Line Number", "Match Text"})
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)

	for _, item := range foundItems {
		if item.Error != nil {
			continue // Skip entries with errors
		}

		name, revision, err := extractNameAndRevision(item.FolderName)
		if err != nil {
			fmt.Printf("Error extracting name and revision: %v\n", err)
			continue
		}

		// Add a row for each item
		table.Append([]string{name, revision, item.FileName, strconv.Itoa(item.LineNum), item.MatchText})
	}

	// Render the table to the output
	table.Render()
	fmt.Println("Found", len(foundItems), "items.")
}
