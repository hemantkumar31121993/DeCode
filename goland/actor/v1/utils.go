package actor

import (
	"regexp"
	"strings"
)

func makeID(name string) string {
	re := regexp.MustCompile(`[^a-z0-9]+`)
	id := re.ReplaceAllString(strings.ToLower(name), "")
	return id
}
