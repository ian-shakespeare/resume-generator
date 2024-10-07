package cli_test

import (
	"resumegenerator/internal/cli"
	"resumegenerator/test/expect"
	"testing"
)

func TestNewArgParser(t *testing.T) {
	t.Run("duplicateName", func(t *testing.T) {
		_, err := cli.NewArgParser([]cli.Flag{
			{Name: "input", Markers: []string{"-o", "--output"}, HasValue: true},
			{Name: "input", Markers: []string{"-i", "--input"}, HasValue: true},
		})

		expect.Error(t, err)
	})

	t.Run("noMarkers", func(t *testing.T) {
		_, err := cli.NewArgParser([]cli.Flag{
			{Name: "input", Markers: []string{}, HasValue: true},
		})

		expect.Error(t, err)
	})

	t.Run("duplicateMarker", func(t *testing.T) {
		_, err := cli.NewArgParser([]cli.Flag{
			{Name: "output", Markers: []string{"-o", "-o"}, HasValue: true},
			{Name: "input", Markers: []string{"-i", "--input"}, HasValue: true},
		})

		expect.Error(t, err)

		_, err = cli.NewArgParser([]cli.Flag{
			{Name: "output", Markers: []string{"-o", "--output"}, HasValue: true},
			{Name: "input", Markers: []string{"-o", "--input"}, HasValue: true},
		})

		expect.Error(t, err)
	})
}

func TestArgParserParse(t *testing.T) {
	p, err := cli.NewArgParser([]cli.Flag{
		{Name: "noValue", Markers: []string{"-n", "--no-value"}, HasValue: false},
		{Name: "input", Markers: []string{"-i", "--input"}, HasValue: true},
		{Name: "output", Markers: []string{"-o", "--output"}, HasValue: true},
	})

	expect.NilError(t, err)

	t.Run("empty", func(t *testing.T) {
		args := make([]string, 0)

		_, err := p.Parse(args)

		expect.Error(t, err)
	})

	t.Run("executable", func(t *testing.T) {
		args := []string{"executable"}

		a, err := p.Parse(args)

		expect.NilError(t, err)
		expect.Equal(t, "executable", a.Executable)
	})

	t.Run("flagNoValue", func(t *testing.T) {
		args := []string{"executable", "-n"}

		a, err := p.Parse(args)

		expect.NilError(t, err)
		expect.Equal(t, "true", a.Flags["noValue"])
	})

	t.Run("flagExpectedValue", func(t *testing.T) {
		args := []string{"executable", "-i"}

		_, err := p.Parse(args)

		expect.Error(t, err)
	})

	t.Run("flagToFlag", func(t *testing.T) {
		args := []string{"executable", "-i", "-o"}

		_, err := p.Parse(args)

		expect.Error(t, err)
	})

	t.Run("redefineFlag", func(t *testing.T) {
		args := []string{"executable", "-i", "value", "-i", "value"}

		_, err := p.Parse(args)

		expect.Error(t, err)
	})

	t.Run("flag", func(t *testing.T) {
		args := []string{"executable", "-i", "value"}

		a, err := p.Parse(args)

		expect.NilError(t, err)
		expect.Equal(t, "value", a.Flags["input"])
	})

	t.Run("flagMultiple", func(t *testing.T) {
		args := []string{"executable", "-i", "in", "-o", "out"}

		a, err := p.Parse(args)

		expect.NilError(t, err)
		expect.Equal(t, "in", a.Flags["input"])
		expect.Equal(t, "out", a.Flags["output"])
	})

	t.Run("positional", func(t *testing.T) {
		args := []string{"executable", "pos1"}

		a, err := p.Parse(args)

		expect.NilError(t, err)
		expect.Equal(t, "pos1", a.Positionals[0])
	})

	t.Run("positionalMultiple", func(t *testing.T) {
		args := []string{"executable", "pos1", "pos2"}

		a, err := p.Parse(args)

		expect.NilError(t, err)
		expect.Equal(t, "pos1", a.Positionals[0])
		expect.Equal(t, "pos2", a.Positionals[1])
	})

	t.Run("flagsAndPositional", func(t *testing.T) {
		args := []string{"executable", "-i", "value", "positional"}

		a, err := p.Parse(args)

		expect.NilError(t, err)
		expect.Equal(t, "value", a.Flags["input"])
		expect.Equal(t, "positional", a.Positionals[0])
	})

	t.Run("flagsPositionalMixed", func(t *testing.T) {
		args := []string{"executable", "-i", "in", "pos1", "-o", "out", "pos2"}

		a, err := p.Parse(args)

		expect.NilError(t, err)
		expect.Equal(t, "in", a.Flags["input"])
		expect.Equal(t, "out", a.Flags["output"])
		expect.Equal(t, "pos1", a.Positionals[0])
		expect.Equal(t, "pos2", a.Positionals[1])
	})
}
