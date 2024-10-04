package cli_test

import (
	"fmt"
	"resumegenerator/internal/cli"
	"testing"
)

func TestColorMessageString(t *testing.T) {
	colors := []cli.Color{
		cli.RESET,
		cli.RED,
		cli.GREEN,
		cli.YELLOW,
		cli.BLUE,
		cli.MAGENTA,
		cli.CYAN,
		cli.GRAY,
		cli.WHITE,
	}

	for _, color := range colors {
		baseMessage := "test"
		message := cli.NewColorMessage(color, baseMessage)

		expected := fmt.Sprintf("%s%s%s", string(color), baseMessage, string(cli.RESET))
		received := message.String()

		if expected != received {
			t.Fatalf("expected %s, received %s", expected, received)
		}
	}
}
