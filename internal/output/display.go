package output

import (
	"fmt"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/stedmanson/grepigee/internal/searcher"
)

func DisplayAsTable(headers []string, data [][]string) {
	if len(data) == 0 {
		fmt.Println("No items found.")
		return
	}

	table := getStandardTable()
	table.SetHeader(headers)

	for _, item := range data {
		table.Append(item)
	}

	table.Render()
	fmt.Println("\nFound", len(data), "items.")
}

func getStandardTable() *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
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

	return table
}

func FormatFoundData(items []searcher.Found) ([]string, [][]string) {
	var data [][]string

	for _, item := range items {
		name, revision, err := extractNameAndRevision(item.FolderName)
		if err != nil {
			continue
		}

		data = append(data, []string{name, revision, item.FileName, strconv.Itoa(item.LineNum), item.MatchText})
	}

	return []string{"Folder Name", "Revision", "File Name", "Line Number", "Match Text"}, data
}
