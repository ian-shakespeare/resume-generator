package cli_test

import (
	"resumegenerator/internal/cli"
	"testing"
)

func TestNewArgParser(t *testing.T) {
	t.Run("duplicateName", func(t *testing.T) {
		_, err := cli.NewArgParser([]cli.Flag{
			{Name: "input", Markers: []string{"-o", "--output"}, Description: "", HasValue: true},
			{Name: "input", Markers: []string{"-i", "--input"}, Description: "", HasValue: true},
		})

		if err == nil {
			t.Fatalf("expected %s, received %s", "error", "nil")
		}
	})

	t.Run("noMarkers", func(t *testing.T) {
		_, err := cli.NewArgParser([]cli.Flag{
			{Name: "input", Markers: []string{}, Description: "", HasValue: true},
		})

		if err == nil {
			t.Fatalf("expected %s, received %s", "error", "nil")
		}
	})

	t.Run("duplicateMarker", func(t *testing.T) {
		_, err := cli.NewArgParser([]cli.Flag{
			{Name: "output", Markers: []string{"-o", "-o"}, Description: "", HasValue: true},
			{Name: "input", Markers: []string{"-i", "--input"}, Description: "", HasValue: true},
		})

		if err == nil {
			t.Fatalf("expected %s, received %s", "error", "nil")
		}

		_, err = cli.NewArgParser([]cli.Flag{
			{Name: "output", Markers: []string{"-o", "--output"}, Description: "", HasValue: true},
			{Name: "input", Markers: []string{"-o", "--input"}, Description: "", HasValue: true},
		})

		if err == nil {
			t.Fatalf("expected %s, received %s", "error", "nil")
		}
	})
}

func TestArgParserParse(t *testing.T) {
	p, err := cli.NewArgParser([]cli.Flag{
		{Name: "noValue", Markers: []string{"-n", "--no-value"}, Description: "", HasValue: false},
		{Name: "input", Markers: []string{"-i", "--input"}, Description: "", HasValue: true},
		{Name: "output", Markers: []string{"-o", "--output"}, Description: "", HasValue: true},
	})

	if err != nil {
		t.Fatalf("expected %s, received %s", "nil", err.Error())
	}

	t.Run("empty", func(t *testing.T) {
		args := make([]string, 0)

		_, err := p.Parse(args)

		if err == nil {
			t.Fatalf("expected %s, received %s", "error", "nil")
		}
	})

	t.Run("executable", func(t *testing.T) {
		args := []string{"executable"}

		a, err := p.Parse(args)

		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		if a.Executable != "executable" {
			t.Fatalf("expected %s, received %s", "executable", a.Executable)
		}
	})

	t.Run("flagNoValue", func(t *testing.T) {
		args := []string{"executable", "-n"}

		a, err := p.Parse(args)

		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		if a.Flags["noValue"] != "true" {
			t.Fatalf("expected %s, received %s", "true", a.Flags["noValue"])
		}
	})

	t.Run("flagExpectedValue", func(t *testing.T) {
		args := []string{"executable", "-i"}

		_, err := p.Parse(args)

		if err == nil {
			t.Fatalf("expected %s, received %s", "error", "nil")
		}
	})

	t.Run("flagToFlag", func(t *testing.T) {
		args := []string{"executable", "-i", "-o"}

		_, err := p.Parse(args)

		if err == nil {
			t.Fatalf("expected %s, received %s", "error", "nil")
		}
	})

	t.Run("redefineFlag", func(t *testing.T) {
		args := []string{"executable", "-i", "value", "-i", "value"}

		_, err := p.Parse(args)

		if err == nil {
			t.Fatalf("expected %s, received %s", "error", "nil")
		}
	})

	t.Run("flag", func(t *testing.T) {
		args := []string{"executable", "-i", "value"}

		a, err := p.Parse(args)

		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		if a.Flags["input"] != "value" {
			t.Fatalf("expected %s, received %s", "value", a.Flags["input"])
		}
	})

	t.Run("flagMultiple", func(t *testing.T) {
		args := []string{"executable", "-i", "in", "-o", "out"}

		a, err := p.Parse(args)

		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		if a.Flags["input"] != "in" {
			t.Fatalf("expected %s, received %s", "in", a.Flags["input"])
		}

		if a.Flags["output"] != "out" {
			t.Fatalf("expected %s, received %s", "out", a.Flags["output"])
		}
	})

	t.Run("positional", func(t *testing.T) {
		args := []string{"executable", "pos1"}

		a, err := p.Parse(args)

		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		if a.Positionals[0] != "pos1" {
			t.Fatalf("expected %s, received %s", "pos1", a.Positionals[0])
		}
	})

	t.Run("positionalMultiple", func(t *testing.T) {
		args := []string{"executable", "pos1", "pos2"}

		a, err := p.Parse(args)

		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		if a.Positionals[0] != "pos1" {
			t.Fatalf("expected %s, received %s", "pos1", a.Positionals[0])
		}

		if a.Positionals[1] != "pos2" {
			t.Fatalf("expected %s, received %s", "pos2", a.Positionals[1])
		}
	})

	t.Run("flagsAndPositional", func(t *testing.T) {
		args := []string{"executable", "-i", "value", "positional"}

		a, err := p.Parse(args)

		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		if a.Flags["input"] != "value" {
			t.Fatalf("expected %s, received %s", "value", a.Flags["input"])
		}

		if a.Positionals[0] != "positional" {
			t.Fatalf("expected %s, received %s", "positional", a.Positionals[0])
		}
	})

	t.Run("flagsPositionalMixed", func(t *testing.T) {
		args := []string{"executable", "-i", "in", "pos1", "-o", "out", "pos2"}

		a, err := p.Parse(args)

		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		if a.Flags["input"] != "in" {
			t.Fatalf("expected %s, received %s", "value", a.Flags["input"])
		}

		if a.Flags["output"] != "out" {
			t.Fatalf("expected %s, received %s", "value", a.Flags["output"])
		}

		if a.Positionals[0] != "pos1" {
			t.Fatalf("expected %s, received %s", "pos1", a.Positionals[0])
		}

		if a.Positionals[1] != "pos2" {
			t.Fatalf("expected %s, received %s", "pos2", a.Positionals[0])
		}
	})
}

func TestArgParserUsage(t *testing.T) {
	t.Run("nonValueFlags", func(t *testing.T) {
		p, err := cli.NewArgParser([]cli.Flag{
			{Name: "output", Markers: []string{"-o", "--output"}, Description: ""},
			{Name: "input", Markers: []string{"-i", "--input"}, Description: ""},
		})

		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		expected := `usage: executable [-oi]`

		received := p.Usage("executable")

		if expected != received {
			t.Fatalf("expected %s, received %s", expected, received)
		}
	})

	t.Run("valueFlags", func(t *testing.T) {
		p, err := cli.NewArgParser([]cli.Flag{
			{Name: "output", Markers: []string{"-o", "--output"}, Description: "", HasValue: true},
			{Name: "input", Markers: []string{"-i", "--input"}, Description: "", HasValue: true},
		})

		if err != nil {
			t.Fatalf("expected %s, received %s", "nil", err.Error())
		}

		expected := `usage: executable [-o output] [-i input]`

		received := p.Usage("executable")

		if expected != received {
			t.Fatalf("expected %s, received %s", expected, received)
		}
	})
}
