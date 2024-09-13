package cli

import (
	"fmt"
	"strings"
)

type Flag struct {
	Name        string
	Markers     []string
	Description string
	HasValue    bool
}

func formatBooleanFlags(flags []Flag) (string, bool) {
	s := ""

	for i := 0; i < len(flags); i += 1 {
		if !flags[i].HasValue {
			firstMarkerChar := flags[i].Markers[0][1]
			s += string(firstMarkerChar)
		}
	}

	if len(s) >= 1 {
		s = "[-" + s + "] "
	}

	return s, len(s) >= 1
}

func formatStringFlags(flags []Flag) (string, bool) {
	s := []string{}

	for i := 0; i < len(flags); i += 1 {
		if flags[i].HasValue {
			s = append(s, fmt.Sprintf("[%s %s]", flags[i].Markers[0], flags[i].Name))
		}
	}

	return strings.Join(s, " "), len(s) >= 1
}

func isFlag(s string) bool {
	if len(s) < 1 {
		return false
	}
	c := strings.Split(s, "")
	return c[0] == "-"
}
