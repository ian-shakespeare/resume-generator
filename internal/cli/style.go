package cli

import (
	"fmt"
	"resumegenerator/internal/assert"
	"resumegenerator/pkg/utils"
	"strconv"
	"strings"
)

const RESET_ESCAPE_CODE string = "\033[0m"

const (
	BOLD           int = 1
	UNDERLINE      int = 4
	BLACK          int = 30
	RED            int = 31
	GREEN          int = 32
	YELLOW         int = 33
	BLUE           int = 34
	MAGENTA        int = 35
	CYAN           int = 36
	WHITE          int = 37
	BRIGHT_BLACK   int = 90
	BRIGHT_RED     int = 91
	BRIGHT_GREEN   int = 92
	BRIGHT_YELLOW  int = 93
	BRIGHT_BLUE    int = 94
	BRIGHT_MAGENTA int = 95
	BRIGHT_CYAN    int = 96
	BRIGHT_WHITE   int = 97
)

type StyledString struct {
	effects []int
	value   string
}

func NewStyledString(value string, effects ...int) *StyledString {
	for _, effect := range effects {
		assert.GreaterThan(effect, 0)
	}

	return &StyledString{
		effects,
		value,
	}
}

func NewBlackString(value string) *StyledString {
	return NewStyledString(value, BLACK)
}

func NewRedString(value string) *StyledString {
	return NewStyledString(value, RED)
}

func NewGreenString(value string) *StyledString {
	return NewStyledString(value, GREEN)
}

func NewYellowString(value string) *StyledString {
	return NewStyledString(value, YELLOW)
}

func NewBlueString(value string) *StyledString {
	return NewStyledString(value, BLUE)
}

func NewMagentaString(value string) *StyledString {
	return NewStyledString(value, MAGENTA)
}

func NewCyanString(value string) *StyledString {
	return NewStyledString(value, CYAN)
}

func NewWhiteString(value string) *StyledString {
	return NewStyledString(value, WHITE)
}

func NewBrightBlackString(value string) *StyledString {
	return NewStyledString(value, BRIGHT_BLACK)
}

func NewBrightRedString(value string) *StyledString {
	return NewStyledString(value, BRIGHT_RED)
}

func NewBrightGreenString(value string) *StyledString {
	return NewStyledString(value, BRIGHT_GREEN)
}

func NewBrightYellowString(value string) *StyledString {
	return NewStyledString(value, BRIGHT_YELLOW)
}

func NewBrightBlueString(value string) *StyledString {
	return NewStyledString(value, BRIGHT_BLUE)
}

func NewBrightMagentaString(value string) *StyledString {
	return NewStyledString(value, BRIGHT_MAGENTA)
}

func NewBrightCyanString(value string) *StyledString {
	return NewStyledString(value, BRIGHT_CYAN)
}

func NewBrightWhiteString(value string) *StyledString {
	return NewStyledString(value, BRIGHT_WHITE)
}

func (s *StyledString) String() string {
	return GetAnsiCode(s.effects...) + s.value + RESET_ESCAPE_CODE
}

func (s *StyledString) StringUnescape() string {
	return GetAnsiCode(s.effects...) + s.value
}

func (s *StyledString) Bold() *StyledString {
	s.effects = append(s.effects, BOLD)
	return s
}

func (s *StyledString) Underline() *StyledString {
	s.effects = append(s.effects, UNDERLINE)
	return s
}

func GetAnsiCode(effects ...int) string {
	codes := utils.Map(effects, func(code, _ int) string {
		return strconv.Itoa(code)
	})

	return fmt.Sprintf("\033[%sm", strings.Join(codes, ";"))
}

func isValidColorCode(code int) bool {
	return (30 <= code && code <= 37) || (90 <= code && code <= 97)
}
