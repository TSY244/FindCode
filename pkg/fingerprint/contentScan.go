package fingerprint

import (
	"os"
	"path/filepath"
	"strings"
)

func IsContainContent(productPath, content string) bool {
	absPath, err := filepath.Abs(productPath)
	if err != nil {
		return false
	}

	entries, err := os.ReadDir(absPath)
	if err != nil {
		return false
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			filePath := absPath + "/" + entry.Name()
			fileContent, err := os.ReadFile(filePath)
			if err != nil {
				return false
			}
			if strings.Contains(string(fileContent), content) {
				return true
			}
		}
	}

	for _, entry := range entries {
		if entry.IsDir() {
			if strings.HasPrefix(entry.Name(), ".") {
				continue
			}
			subDirPath := filepath.Join(absPath, entry.Name())
			if ret := IsContainContent(subDirPath, content); ret {
				return ret
			}
		}
	}
	return false
}
