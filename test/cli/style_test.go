package cli_test

import (
	"fmt"
	"resumegenerator/internal/cli"
	"resumegenerator/test/expect"
	"testing"
)

func TestNewStyledString(t *testing.T) {
	t.Run("invalidEffect", func(t *testing.T) {
		defer expect.Panic(t)

		cli.NewStyledString("test", 0)
	})
}

func TestStyledStringString(t *testing.T) {
	t.Run("escaped", func(t *testing.T) {
		s := cli.NewStyledString("test", 1)
		expect.Equal(t, fmt.Sprintf("\033[1mtest%s", cli.RESET_ESCAPE_CODE), s.String())

		s = cli.NewStyledString("test", 1, 2, 3)
		expect.Equal(t, fmt.Sprintf("\033[1;2;3mtest%s", cli.RESET_ESCAPE_CODE), s.String())
	})

	t.Run("unescaped", func(t *testing.T) {
		s := cli.NewStyledString("test", 1)
		expect.Equal(t, "\033[1mtest", s.StringUnescape())

		s = cli.NewStyledString("test", 1, 2, 3)
		expect.Equal(t, "\033[1;2;3mtest", s.StringUnescape())
	})
}

func TestStyledStringBold(t *testing.T) {
	s := cli.NewStyledString("test").Bold()
	expect.Equal(t, fmt.Sprintf("\033[1mtest"), s.StringUnescape())
}

func TestStyledStringUnderline(t *testing.T) {
	s := cli.NewStyledString("test").Underline()
	expect.Equal(t, fmt.Sprintf("\033[4mtest"), s.StringUnescape())
}

func TestGetAnsiCode(t *testing.T) {
	t.Run("single", func(t *testing.T) {
		code := cli.GetAnsiCode(1)
		expect.Equal(t, code, "\033[1m")
	})

	t.Run("multiple", func(t *testing.T) {
		code := cli.GetAnsiCode(1, 2, 3)
		expect.Equal(t, code, "\033[1;2;3m")
	})
}
