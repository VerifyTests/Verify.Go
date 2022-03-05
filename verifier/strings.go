package verifier

import (
	"strings"
)

func removeStringItem(slice []string, item string) []string {
	index := indexOfString(slice, item)
	if index != -1 {
		return removeStringIndex(slice, index)
	}
	return slice
}

func removeStringIndex(s []string, i int) []string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func indexOfString(slice []string, element string) int {
	for k, v := range slice {
		if element == v {
			return k
		}
	}
	return -1 // not found.
}

func fixNewlines(value string) string {
	value = strings.ReplaceAll(value, "\r\n", "\n")
	value = strings.ReplaceAll(value, "\r", "\n")
	return value
}
