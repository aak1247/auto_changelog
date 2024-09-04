package utils

import (
	"regexp"
	"strconv"
	"strings"
)

func IsMultiline(s string) bool {
	return strings.Contains(s, "\n")
}

// CompareVersions compares two version strings
func CompareVersions(version1, version2 string) int {
	// Split the version strings into parts
	parts1 := splitVersion(version1)
	parts2 := splitVersion(version2)

	// Compare each part
	for i := 0; i < len(parts1) && i < len(parts2); i++ {
		if parts1[i] != parts2[i] {
			// If both parts are numeric, compare their numeric values
			if isNumeric(parts1[i]) && isNumeric(parts2[i]) {
				num1, _ := strconv.Atoi(parts1[i])
				num2, _ := strconv.Atoi(parts2[i])
				if num1 != num2 {
					return num1 - num2
				}
			} else {
				// Otherwise, compare their lexicographical order
				return strings.Compare(parts1[i], parts2[i])
			}
		}
	}

	// If all parts are equal, compare the length of the parts
	return len(parts1) - len(parts2)
}

// splitVersion splits a version string into parts
func splitVersion(version string) []string {
	re := regexp.MustCompile(`([a-zA-Z]+)|(\d+)`)
	return re.FindAllString(version, -1)
}

// isNumeric checks if a string is numeric
func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}
