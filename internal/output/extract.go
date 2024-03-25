package output

import (
	"fmt"
	"strings"
)

func extractNameAndRevision(folderName string) (string, string, error) {
	parts := strings.Split(folderName, "-")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid folder name format: %s", folderName)
	}

	revision := parts[len(parts)-1]
	name := strings.Join(parts[:len(parts)-1], "-")

	return name, revision, nil
}
