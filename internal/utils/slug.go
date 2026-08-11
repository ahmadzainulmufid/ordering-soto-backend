package utils

import "strings"

func GenerateSlug(value string) string {
	value = strings.TrimSpace(value)
	value = strings.ToLower(value)
	value = strings.Join(
		strings.Fields(value),
		"-",
	)

	return value
}