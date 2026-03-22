package cache

import "strings"

func buildKey(args ...string) string {
	return strings.Join(args, ":")
}
