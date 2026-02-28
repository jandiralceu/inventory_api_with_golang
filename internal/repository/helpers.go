package repository

import "strings"

// sanitizeLike escapes LIKE/ILIKE wildcard characters (% and _)
// so they are treated as literal characters in search queries.
func sanitizeLike(s string) string {
	s = strings.ReplaceAll(s, "%", "\\%")
	s = strings.ReplaceAll(s, "_", "\\_")
	return s
}
