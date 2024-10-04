package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"time"
)

func FormatPrompt(message string) []byte {
	leader := NewColorMessage(YELLOW, ">")
	return []byte(fmt.Sprintf("%s %s: %s", leader, message, CYAN))
}

func Prompt(r io.Reader, w io.Writer, message string) (string, error) {
	_, err := w.Write(FormatPrompt(message))
	if err != nil {
		return "", err
	}

	defer w.Write([]byte(RESET))

	scanner := bufio.NewScanner(r)

	if scanner.Scan() {
		return scanner.Text(), nil
	}

	return "", errors.New("No input")
}

func PromptDate(r io.Reader, w io.Writer, message string) (time.Time, error) {
	s, err := Prompt(r, w, message)
	if err != nil {
		return time.Unix(0, 0), err
	}

	return time.Parse(time.DateOnly, s)
}

func PromptList(r io.Reader, w io.Writer, message string, exitSequence string) ([]string, error) {
	w.Write(FormatPrompt(message))
	w.Write([]byte("\n" + NewColorMessage(YELLOW, "> ").String() + string(CYAN)))
	defer w.Write([]byte(RESET))

	scanner := bufio.NewScanner(r)

	elements := []string{}

	for scanner.Scan() {
		leader := NewColorMessage(YELLOW, "> ").String() + string(CYAN)
		w.Write([]byte(leader))

		s := scanner.Text()

		if s == exitSequence {
			break
		}

		elements = append(elements, s)
	}

	return elements, nil
}
