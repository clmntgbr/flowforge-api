package consumer

import (
	"strconv"
	"strings"
)

func parseMajorMinor(index string) (major int, minor int, hasMinor bool) {
	parts := strings.Split(index, ".")
	major, _ = strconv.Atoi(parts[0])
	if len(parts) < 2 {
		return major, 0, false
	}
	minor, _ = strconv.Atoi(parts[1])
	return major, minor, true
}
