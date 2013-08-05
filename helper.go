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

func percentile(data []float64, p float64) float64 {
	i := int(p * float64(len(data)) + 0.5)
    return data[i-1]
}

func data2Percentage(data []float64) []float64 {
	n := len(data)
	out := make([]float64, 0, n)
	var total float64

	for i := 0; i < n; i++ {
		total += data[i]
	}

	for i := 0; i < n; i++ {
		out = append(out, (data[i]/total)*100)
	}
	return out
}

