package cli_test

import (
	"bytes"
	"resumegenerator/internal/cli"
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
		if err == nil {
			t.Fatal("expected error, received nil")
		}
	})

	t.Run("nonEmpty", func(t *testing.T) {
		input := "test input"
		prompt := "test prompt"

		in := bytes.NewBuffer([]byte(input))
		out := bytes.NewBuffer(nil)

		s, err := cli.Prompt(in, out, prompt)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		if s != input {
			t.Fatalf("expected %s, received %s", input, s)
		}

		formattedPrompt := string(cli.FormatPrompt(prompt)) + string(cli.RESET)
		if out.String() != formattedPrompt {
			t.Fatalf("expected %s, received %s", formattedPrompt, out.String())
		}
	})
}

func TestPromptDate(t *testing.T) {
	t.Run("invalid", func(t *testing.T) {
		input := "NOT A DATE"
		prompt := "test prompt"

		in := bytes.NewBuffer([]byte(input))
		out := bytes.NewBuffer(nil)

		_, err := cli.PromptDate(in, out, prompt)
		if err == nil {
			t.Fatal("expected error, received nil")
		}
	})

	t.Run("valid", func(t *testing.T) {
		input := "2000-01-01"
		prompt := "test prompt"

		in := bytes.NewBuffer([]byte(input))
		out := bytes.NewBuffer(nil)

		expected, err := time.Parse(time.DateOnly, input)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		received, err := cli.PromptDate(in, out, prompt)
		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		if !expected.Equal(received) {
			t.Fatalf("expected %v, received %v", expected, received)
		}
	})
}

func TestPromptList(t *testing.T) {
	input := "one\ntwo\nthree\nfour"
	prompt := "test prompt"
	exitSeq := "q"

	in := bytes.NewBuffer([]byte(input + "\n" + exitSeq))
	out := bytes.NewBuffer(nil)

	l, err := cli.PromptList(in, out, prompt, exitSeq)
	if err != nil {
		t.Fatalf("expected %s, received %s", "nil", err.Error())
	}

	inputElems := strings.Split(input, "\n")
	if len(l) != len(inputElems) {
		t.Fatalf("expected %d, received %d", len(inputElems), len(l))
	}

	for i, elem := range l {
		if elem != inputElems[i] {
			t.Fatalf("expected %s, received %s", inputElems[i], elem)
		}
	}
}
