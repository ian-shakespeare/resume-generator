package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"resumegenerator/internal/cli"
	"resumegenerator/pkg/resume"
	"resumegenerator/pkg/utils"
	"strings"
	"time"
)

const CLI_VERSION string = "1.1.0"
const CLI_HELP = `Resume Generator command line tool

%s
  resumegen/%s

%s
  $ resumegen [COMMAND]

%s
  version   Print CLI version information
  help      Print this message
  create    Create a new resume data file
  generate  Generate a resume from a data file

`
const TEMPLATE_DIR string = "/assets/templates/"

var LINE_LEADER = cli.NewYellowString("> ").String()

func main() {
	sout := log.New(os.Stdout, "", 0)
	serr := log.New(os.Stderr, "", 0)

	/*
		ex, err := os.Executable()
		if err != nil {
			serr.Fatal(err.Error())
		}

		dir := filepath.Dir(ex)
	*/

	if len(os.Args) < 2 {
		printHelp(sout)
		return
	}

	command := os.Args[1]
	switch command {
	case "help":
		printHelp(sout)
		return
	case "version":
		printVersion(sout)
		return
	case "create":
		err := createResumeData(os.Stdin, os.Stdout)
		if err != nil {
			serr.Print(err.Error())
		}
		return
	}
}

func printHelp(out *log.Logger) {
	out.Printf(
		CLI_HELP,
		cli.NewBrightWhiteString("VERSION").Bold().String(),
		CLI_VERSION,
		cli.NewBrightWhiteString("USAGE").Bold().String(),
		cli.NewBrightWhiteString("COMMANDS").Bold().String(),
	)
}

func printVersion(out *log.Logger) {
	out.Println(CLI_VERSION)
}

func createResumeData(r io.Reader, w io.Writer) error {
	prompts := []string{
		"Name", "Job Title", "Email",
		"Phone Number", "Prelude", "Location (optional)",
		"LinkedIn ID (optional)", "GitHub Username (optional)", "Facebook Username (optional)",
		"Instagram Username (optional)", "Twitter Handle (optional)", "Portfolio Link (optional)",
	}

	responses := []string{}
	for _, prompt := range prompts {
		responses = append(responses, strings.Trim(promptUntilSuccess(r, w, formatPrompt(prompt)), " "))
	}

	skills, err := cli.PromptList(
		r,
		w,
		formatListPrompt("Skills"),
		"",
		LINE_LEADER,
		cli.GetAnsiCode(cli.CYAN),
	)
	if err != nil {
		return err
	}

	educations := promptEducations(r, w)

	res := resume.Resume{
		Name:        responses[0],
		Title:       responses[1],
		Email:       responses[2],
		PhoneNumber: responses[3],
		Prelude:     responses[4],
		Location:    responses[5],
		LinkedIn:    responses[6],
		Github:      responses[7],
		Facebook:    responses[8],
		Instagram:   responses[9],
		Twitter:     responses[10],
		Portfolio:   responses[11],
		Skills: utils.Map(skills, func(s string, _ int) string {
			return strings.Trim(s, " ")
		}),
		Educations: educations,
	}

	fmt.Println(res)
	return nil
}

func formatPrompt(prompt string) string {
	return LINE_LEADER + prompt + ": " + cli.GetAnsiCode(cli.CYAN)
}

func formatListPrompt(prompt string) string {
	return prompt + " (enter to finish):\n" + LINE_LEADER + cli.GetAnsiCode(cli.CYAN)
}

func promptUntilSuccess(r io.Reader, w io.Writer, prompt string) string {
	response, err := cli.Prompt(r, w, prompt)
	for err != nil {
		response, err = cli.Prompt(r, w, prompt)
	}
	return response
}

func promptDateUntilSuccess(r io.Reader, w io.Writer, prompt string) time.Time {
	response, err := cli.PromptDate(r, w, prompt)
	for err != nil {
		w.Write([]byte("Invalid date, please enter in the format YYYY-MM-DD\n"))
		response, err = cli.PromptDate(r, w, prompt)
	}
	return response
}

func promptEducations(r io.Reader, w io.Writer) []resume.Education {
	var educations []resume.Education
	for {
		degree := promptUntilSuccess(r, w, formatPrompt("Degree (enter to finish educations)"))
		if degree == "" {
			break
		}

		field := promptUntilSuccess(r, w, formatPrompt("Field of Study"))
		institution := promptUntilSuccess(r, w, formatPrompt("Institution"))

		began := promptDateUntilSuccess(r, w, formatPrompt("Date Began (YYYY-MM-DD)"))
		current := promptUntilSuccess(r, w, formatPrompt("Currently Studying (y/n)"))
		var finished *time.Time
		if current != "y" {
			val := promptDateUntilSuccess(r, w, formatPrompt("Date Finished (YYYY-MM-DD"))
			finished = &val
		}

		location := promptUntilSuccess(r, w, formatPrompt("Location (optional)"))
		gpa := promptUntilSuccess(r, w, formatPrompt("GPA (optional)"))

		educations = append(educations,
			resume.Education{
				Degree:       strings.Trim(degree, " "),
				FieldOfStudy: strings.Trim(field, " "),
				Institution:  strings.Trim(institution, " "),
				Began:        began,
				Current:      current == "y",
				Finished:     finished,
				Location:     strings.Trim(location, " "),
				GPA:          strings.Trim(gpa, " "),
			},
		)
	}

	return educations
}
