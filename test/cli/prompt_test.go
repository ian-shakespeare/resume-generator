package cli_test

import (
	"bytes"
	"resumegenerator/internal/cli"
	"resumegenerator/test/expect"
	"strings"
	"testing"
	"time"
)

func TestPrompt(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		input := ""
		prompt := "test prompt"

		in := bytes.NewBuffer([]byte(input))
		out := bytes.NewBuffer(nil)

		_, err := cli.Prompt(in, out, prompt)
		expect.Error(t, err)
	})

	t.Run("nonEmpty", func(t *testing.T) {
		input := "test input"
		prompt := "test prompt"

		in := bytes.NewBuffer([]byte(input))
		out := bytes.NewBuffer(nil)

		s, err := cli.Prompt(in, out, prompt)
		expect.NilError(t, err)
		expect.Equal(t, input, s)

		formattedPrompt := prompt + string(cli.RESET_ESCAPE_CODE)
		expect.Equal(t, formattedPrompt, out.String())
	})
}

func TestPromptDate(t *testing.T) {
	t.Run("invalid", func(t *testing.T) {
		input := "NOT A DATE"
		prompt := "test prompt"

		in := bytes.NewBuffer([]byte(input))
		out := bytes.NewBuffer(nil)

		_, err := cli.PromptDate(in, out, prompt)
		expect.Error(t, err)
	})

	t.Run("valid", func(t *testing.T) {
		input := "2000-01-01"
		prompt := "test prompt"

		in := bytes.NewBuffer([]byte(input))
		out := bytes.NewBuffer(nil)

		expected, err := time.Parse(time.DateOnly, input)
		expect.NilError(t, err)

		received, err := cli.PromptDate(in, out, prompt)
		expect.NilError(t, err)
		expect.True(t, expected.Equal(received))
	})
}

func TestPromptList(t *testing.T) {
	input := "one\ntwo\nthree\nfour"
	prompt := "test prompt"
	exitSeq := "q"

	in := bytes.NewBuffer([]byte(input + "\n" + exitSeq))
	out := bytes.NewBuffer(nil)

	l, err := cli.PromptList(in, out, prompt, exitSeq)
	expect.NilError(t, err)

	inputElems := strings.Split(input, "\n")
	expect.Equal(t, len(inputElems), len(l))

	for i, elem := range l {
		expect.Equal(t, inputElems[i], elem)
	}
}
