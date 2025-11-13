package str

import "strings"

func AngleWrap(strs []string) string {
	return "<" + strings.Join(strs, "> <") + ">"
}
