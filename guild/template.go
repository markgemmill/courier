package guild

import (
	"fmt"
	"os"
	"strings"
)

func emptyString(str string) bool {
	return strings.TrimSpace(str) == ""
}

// LoadTemplateString does some basic checks to determine if the input string
// might be a file path, and if so attempts to open the file for reading.  If it is a
// readable file path, then it returns the contents of the file, otherwise returns
// the original string value.
func LoadTemplateString(mightBeFilePath string) string {
	if emptyString(mightBeFilePath) {
		return ""
	}
	// does it look like a filepath?
	if strings.Contains(mightBeFilePath, "\\") || strings.Contains(mightBeFilePath, "/") || len(mightBeFilePath) < 1024 {
		// is it an actual file?
		info, err := os.Stat(mightBeFilePath)
		if err == nil && !info.IsDir() {
			// okay, read the file
			content, err := os.ReadFile(mightBeFilePath)
			if err == nil {
				return string(content)
			}
		}
		fmt.Printf("LoadTemplateString: %s", err)
	}
	// else just return the string
	return mightBeFilePath
}
