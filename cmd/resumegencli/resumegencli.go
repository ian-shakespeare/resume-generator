package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"resumegenerator/internal/cli"
	"resumegenerator/pkg/generator"
	"resumegenerator/pkg/resume"
	"resumegenerator/pkg/utils"
	"strings"

	"github.com/google/uuid"
)

const CLI_VERSION string = "1.0.0"
const CLI_HELP string = `usage: resumegen [resumepath] [-hv] [-t template]
  ------- Listing options -------
  -t template Set output template.
  ------- Miscellaneous options -------
  -v version  Print version.
  -h help     Print usage and this help message.
`
const TEMPLATE_DIR string = "./assets/templates"

func main() {
	so := log.New(os.Stdout, "", 0)
	se := log.New(os.Stderr, "", 0)

	p, err := cli.NewArgParser([]cli.Flag{
		{Name: "version", Markers: []string{"-v", "--version"}, HasValue: false},
		{Name: "help", Markers: []string{"-h", "--help"}, HasValue: false},
		{Name: "template", Markers: []string{"-t", "--template"}, HasValue: true},
	})
	if err != nil {
		se.Fatal(err)
	}

	args, err := p.Parse(os.Args)
	if err != nil {
		se.Fatal(err)
	}

	if args.Flags["help"] == "true" {
		so.Print(CLI_HELP)
	}

	if args.Flags["version"] == "true" {
		so.Println(CLI_VERSION)
	}

	if len(args.Positionals) < 1 {
		se.Fatal("missing positional argument resumepath")
	}

	r, err := getResume(args.Positionals[0])
	if err != nil {
		so.Fatal(err)
	}

	selectedTmpl, exists := args.Flags["template"]
	if !exists {
		selectedTmpl = "default"
	}

	tmpl, err := os.ReadFile(TEMPLATE_DIR + selectedTmpl + ".html")
	if err != nil {
		se.Fatal(err.Error())
	}

	output, err := generator.GenerateHtml(&r, string(tmpl))
	if err != nil {
		se.Fatal(err.Error())
	}

	outputDir := os.TempDir()
	outputName := "resume-" + uuid.New().String()
	outputPath := outputDir + outputName

	err = os.WriteFile(outputPath, []byte(output), 644)
	if err != nil {
		se.Fatal(err.Error())
	}

	cli.OpenUrl(outputPath)

	so.Printf("Resume saved at %s\n", outputPath)
}

func getResume(resumePath string) (resume.Resume, error) {
	acceptedExt := []string{"json", "yaml", "yml"}

	resumePathParts := strings.Split(resumePath, ".")
	if len(resumePathParts) < 1 {
		return resume.Resume{}, errors.New("cannot identify file extension")
	}

	ext := resumePathParts[len(resumePathParts)-1]
	if hasValidExt := utils.Contains(acceptedExt, ext); hasValidExt {
		return resume.Resume{}, fmt.Errorf("invalid file extension %s", ext)
	}

	b, err := os.ReadFile(resumePath)
	if err != nil {
		return resume.Resume{}, err
	}

	return resume.New(b)
}
