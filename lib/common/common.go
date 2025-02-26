package common

import "strings"

func DotJoin(opts ...string) string {
	return strings.Join(opts, ".")
}
