package pkg

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var (
	reNonAlphanumeric = regexp.MustCompile("[^a-z0-9]+")
)

// Slugify converts a string into a URL-friendly slug.
// It handles accents, special characters, and multiple spaces.
func Slugify(text string) string {
	// 1. Normalize and remove accents
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, text)

	// 2. To lower case
	result = strings.ToLower(result)

	// 3. Replace non-alphanumeric with dashes
	result = reNonAlphanumeric.ReplaceAllString(result, "-")

	// 4. Trim leading/trailing dashes
	result = strings.Trim(result, "-")

	return result
}
