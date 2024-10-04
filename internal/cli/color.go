package cli

type Color string

const (
	RESET   Color = "\033[0m"
	RED     Color = "\033[31m"
	GREEN   Color = "\033[32m"
	YELLOW  Color = "\033[33m"
	BLUE    Color = "\033[34m"
	MAGENTA Color = "\033[35m"
	CYAN    Color = "\033[36m"
	GRAY    Color = "\033[37m"
	WHITE   Color = "\033[97m"
)

type ColorMessage struct {
	color   Color
	message string
}

func NewColorMessage(c Color, message string) *ColorMessage {
	return &ColorMessage{
		c,
		message,
	}
}

func (c *ColorMessage) String() string {
	return string(c.color) + c.message + string(RESET)
}
