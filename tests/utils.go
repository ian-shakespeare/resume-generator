package tests

import "fmt"

func FormatExpected(expected string, received string) string {
	return fmt.Sprintf("expected %s, received %s", expected, received)
}
