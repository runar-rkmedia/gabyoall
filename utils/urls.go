package utils

import "strings"

func PrepareUrl(s string) string {
	if !strings.HasPrefix(s, "http") {
		s = "https://" + s
	}
	s = strings.TrimSuffix(s, "/")
	return s
}

func UrlsAreEqual(a, b string) bool {
	A := PrepareUrl(a)
	B := PrepareUrl(b)
	return A == B
}
