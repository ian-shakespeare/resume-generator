package cli

import (
	"bufio"
	"errors"
	"io"
	"strings"
	"time"
)

func Prompt(r io.Reader, w io.Writer, message string) (string, error) {
	_, err := w.Write([]byte(message))
	if err != nil {
		return "", err
	}

	defer w.Write([]byte(RESET_ESCAPE_CODE))

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

func PromptList(r io.Reader, w io.Writer, message string, exitSequence string, lineLeader ...string) ([]string, error) {
	w.Write([]byte(message))

	scanner := bufio.NewScanner(r)

	elements := []string{}
	leader := strings.Join(lineLeader, "")

	for scanner.Scan() {
		w.Write([]byte(leader))

		s := scanner.Text()

		if s == exitSequence {
			break
		}

		elements = append(elements, s)
	}

	return elements, nil
}
