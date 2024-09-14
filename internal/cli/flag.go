package cli

import "strings"

type Flag struct {
	Name        string
	Markers     []string
	HasValue    bool
}

func isFlag(s string) bool {
	if len(s) < 1 {
		return false
	}
	c := strings.Split(s, "")
	return c[0] == "-"
}
