package utils

import "strings"

func NewSlug(slug string) string {
	slugCandidate := strings.Split(slug, " ")

	return strings.Join(slugCandidate, "-")
}
