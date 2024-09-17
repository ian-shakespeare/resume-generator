package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"resumegenerator/internal/cli"
	"resumegenerator/pkg/generator"
	"resumegenerator/pkg/resume"
	"resumegenerator/pkg/utils"
	"strings"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

const CLI_VERSION string = "1.0.0"
const CLI_HELP string = `usage: resumegen [resumepath] [-efhv] [-t template] [-o outputdir]
  ------- Generator options -------
  -e example    Use example data.
  -f format     Set output format <html/pdf>.
  -o outputdir  Set output directory.
  -t template   Path to template.
  ------- Miscellaneous options -------
  -v version    Print version.
  -h help       Print usage and this help message.
`
const TEMPLATE_DIR string = "/assets/templates/"

func main() {
	so := log.New(os.Stdout, "", 0)
	se := log.New(os.Stderr, "", 0)

	ex, err := os.Executable()
	if err != nil {
		se.Fatal(err.Error())
	}

	dir := filepath.Dir(ex)

	p, err := cli.NewArgParser([]cli.Flag{
		{Name: "version", Markers: []string{"-v", "--version"}, HasValue: false},
		{Name: "help", Markers: []string{"-h", "--help"}, HasValue: false},
		{Name: "template", Markers: []string{"-t", "--template"}, HasValue: true},
		{Name: "outputdir", Markers: []string{"-o", "--outputdir"}, HasValue: true},
		{Name: "example", Markers: []string{"-e", "--example"}, HasValue: false},
		{Name: "format", Markers: []string{"-f", "--format"}, HasValue: true},
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

	r, err := getResume(&args)
	if err != nil {
		se.Fatal(err.Error())
	}

	selectedTmpl, exists := args.Flags["template"]
	if !exists {
		selectedTmpl = "default"
	}

	tmpl, err := os.ReadFile(dir + TEMPLATE_DIR + selectedTmpl + ".html")
	if err != nil {
		se.Fatal(err.Error())
	}

	var b []byte
	if args.Flags["format"] == "html" {
		b, err = generator.GenerateHtml(&r, tmpl)
	} else {
		b, err = generator.GeneratePdf(&r, tmpl)
	}
	if err != nil {
		se.Fatal(err.Error())
	}

	o := getOutputPath(&args)

	err = os.WriteFile(o, []byte(b), 0777)
	if err != nil {
		se.Fatal(err.Error())
	}

	so.Printf("Resume saved at %s\n", o)
}

func getOutputPath(args *cli.Arguments) string {
	dir, exists := args.Flags["outputdir"]
	if !exists {
		dir = os.TempDir()
	}
	if dir == "." {
		dir = ""
	}

	ext := ".pdf"
	if args.Flags["format"] == "html" {
		ext = ".html"
	}

	name := "resume-" + uuid.New().String()

	return dir + name + ext
}

func getResume(args *cli.Arguments) (resume.Resume, error) {
	var path string
	if args.Flags["example"] == "true" || len(args.Positionals) < 1 {
		return resume.Example(), nil
	} else {
		path = args.Positionals[0]
	}

	acceptedExt := []string{"json", "yaml", "yml"}

	resumePathParts := strings.Split(path, ".")
	if len(resumePathParts) < 1 {
		return resume.Resume{}, errors.New("cannot identify file extension")
	}

	ext := resumePathParts[len(resumePathParts)-1]
	if hasValidExt := utils.Contains(acceptedExt, ext); !hasValidExt {
		return resume.Resume{}, fmt.Errorf("invalid file extension %s", ext)
	}

	b, err := os.ReadFile(path)
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
