package cli

import (
	"errors"
	"fmt"
	"resumegenerator/pkg/utils"
)

type ArgParser struct {
	flags []Flag
}

type Arguments struct {
	Executable  string
	Flags       map[string]string
	Positionals []string
}

func NewArgParser(flags []Flag) (*ArgParser, error) {
	existingNames := make([]string, 0)
	existingMarkers := make([]string, 0)

	for _, flag := range flags {
		if utils.Contains(existingNames, flag.Name) {
			return nil, fmt.Errorf("cannot create flag with duplicate name %s", flag.Name)
		}
		if len(flag.Markers) < 1 {
			return nil, fmt.Errorf("cannot create flag %s with less than 1 marker", flag.Name)
		}
		for _, marker := range flag.Markers {
			if utils.Contains(existingMarkers, marker) {
				return nil, fmt.Errorf("cannot create flag with duplicate marker %s", marker)
			}
			existingMarkers = append(existingMarkers, marker)
		}
		existingNames = append(existingNames, flag.Name)
	}

	return &ArgParser{flags}, nil
}

func (a *ArgParser) Parse(args []string) (Arguments, error) {
	if len(args) < 1 {
		return Arguments{}, errors.New("cannot process invalid arguments")
	}

	executable := args[0]
	flags := make(map[string]string)
	positionals := make([]string, 0)

	i := 1
	for i < len(args) {
		arg := args[i]
		if len(arg) < 1 {
			return Arguments{}, errors.New("cannot process empty argument")
		}

		if isFlag(arg) {
			flagIndex := utils.Find(a.flags, func(f Flag) bool {
				return utils.Some(f.Markers, func(marker string) bool {
					return marker == arg
				})
			})
			if flagIndex < 0 {
				return Arguments{}, fmt.Errorf("unrecognized flag %s", arg)
			}

			flagName := a.flags[flagIndex].Name

			if !a.flags[flagIndex].HasValue {
				flags[flagName] = "true"
				i += 1
			} else {
				if i+1 >= len(args) {
					return Arguments{}, fmt.Errorf("%s must have a value", arg)
				}

				nextArg := args[i+1]
				if isFlag(nextArg) {
					return Arguments{}, fmt.Errorf("%s expected argument, received %s", arg, nextArg)
				}

				if _, exists := flags[flagName]; exists {
					return Arguments{}, errors.New("redeclared flag")
				}

				flags[flagName] = nextArg
				i += 2
			}
		} else {
			positionals = append(positionals, arg)
			i += 1
		}
	}

	return Arguments{Executable: executable, Flags: flags, Positionals: positionals}, nil
}
