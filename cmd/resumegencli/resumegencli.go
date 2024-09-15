package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"resumegenerator/internal/cli"
	"resumegenerator/pkg/generator"
	"resumegenerator/pkg/resume"
	"resumegenerator/pkg/utils"
	"strings"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

const CLI_VERSION string = "1.0.0"
const CLI_HELP string = `usage: resumegen [resumepath] [-ehvO] [-t template] [-o outputdir]
  ------- Generator options -------
  -e example    Use example data.
  -o outputdir  Set output directory.
  -t template   Set output template.
  -O open       Open the output on completion.
  ------- Miscellaneous options -------
  -v version    Print version.
  -h help       Print usage and this help message.
`
const TEMPLATE_DIR string = "./assets/templates/"

func main() {
	so := log.New(os.Stdout, "", 0)
	se := log.New(os.Stderr, "", 0)

	p, err := cli.NewArgParser([]cli.Flag{
		{Name: "version", Markers: []string{"-v", "--version"}, HasValue: false},
		{Name: "help", Markers: []string{"-h", "--help"}, HasValue: false},
		{Name: "template", Markers: []string{"-t", "--template"}, HasValue: true},
		{Name: "open", Markers: []string{"-O", "--open"}, HasValue: false},
		{Name: "outputdir", Markers: []string{"-o", "--outputdir"}, HasValue: true},
		{Name: "example", Markers: []string{"-e", "--example"}, HasValue: false},
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

	_, useExampleResume := args.Flags["example"]

	if !useExampleResume && len(args.Positionals) < 1 {
		se.Fatal("missing positional argument resumepath")
	}

	var r resume.Resume

	if useExampleResume {
		r = resume.Example()
	} else {
		r, err = getResume(args.Positionals[0])
		if err != nil {
			so.Fatal(err)
		}
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

	var outputDir string
	if outputDir, exists = args.Flags["outputdir"]; !exists {
		outputDir = os.TempDir()
	}
	outputName := "resume-" + uuid.New().String() + ".html"
	outputPath := outputDir + outputName

	err = os.WriteFile(outputPath, []byte(output), 0777)
	if err != nil {
		se.Fatal(err.Error())
	}

	if args.Flags["open"] == "true" {
		cli.OpenUrl(outputPath)
	}

	so.Printf("Resume saved at %s\n", outputPath)
}

func getResume(resumePath string) (resume.Resume, error) {
	acceptedExt := []string{"json", "yaml", "yml"}

	resumePathParts := strings.Split(resumePath, ".")
	if len(resumePathParts) < 1 {
		return resume.Resume{}, errors.New("cannot identify file extension")
	}

	ext := resumePathParts[len(resumePathParts)-1]
	if hasValidExt := utils.Contains(acceptedExt, ext); !hasValidExt {
		return resume.Resume{}, fmt.Errorf("invalid file extension %s", ext)
	}

	b, err := os.ReadFile(resumePath)
	if err != nil {
		return resume.Resume{}, err
	}

	var r resume.Resume
	if ext == "json" {
		if err = json.Unmarshal(b, &r); err != nil {
			return resume.Resume{}, err
		}
	} else if ext == "yaml" || ext == "yml" {
		if err = yaml.Unmarshal(b, &r); err != nil {
			return resume.Resume{}, err
		}
	}

	return r, nil
}
