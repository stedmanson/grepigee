// File: searcher/searcher.go

package searcher

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Found represents a found match within a file, including the file name, line number, proxy name, and the text that matched.
type Found struct {
	FileName   string
	LineNum    int
	FolderName string
	MatchText  string
	Error      error
}

// SearchInFile searches for a pattern in a file and returns a slice of Found detailing where the pattern was found.
func SearchInFile(filePath, pattern, folderName string) ([]Found, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	var results []Found
	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		if re.MatchString(line) {
			results = append(results, Found{
				FileName:   filePath,
				LineNum:    lineNum,
				FolderName: folderName,
				MatchText:  strings.TrimSpace(line),
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return results, err
	}

	return results, nil
}

// SearchInDirectory searches for a pattern in all files within a directory (recursively) and returns a slice of Found structs detailing where the pattern was found.
func SearchInDirectory(directory, pattern string) ([]Found, error) {
	var matched []Found
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil // Skip directories
		}

		// Extract the proxy name (second-level folder under 'directory')
		parts := strings.Split(strings.TrimPrefix(path, directory), string(os.PathSeparator))
		folderName := ""
		if len(parts) > 2 {
			folderName = parts[1] // The proxy name should be the second part
		}

		matches, err := SearchInFile(path, pattern, folderName)
		if err != nil {
			matched = append(matched, Found{Error: err})
			return nil // Continue searching; one error shouldn't stop the entire search.
		}

		matched = append(matched, matches...)

		return nil
	})

	return matched, err
}
