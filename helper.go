package main

import (
	"path/filepath"
	"regexp"
	"strings"
)

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func safeFilename(v string) string {
	v = strings.ToLower(v)

	re := regexp.MustCompile("[^a-z0-9]")
	v = re.ReplaceAllLiteralString(v, "-")

	re = regexp.MustCompile("[-]{2,}")
	v = re.ReplaceAllLiteralString(v, "-")
	return filepath.Clean(v)
}
