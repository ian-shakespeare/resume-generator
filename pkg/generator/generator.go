package generator

import (
	"bytes"
	"fmt"
	"html/template"
	"resumegenerator/pkg/resume"
	"time"
)

func monthYear(t *time.Time) string {
	if t == nil {
		return ""
	}
	return fmt.Sprintf("%s %d", t.Month().String(), t.Year())
}

func GenerateHtml(r *resume.Resume, tmpl string) (string, error) {
	t, err := template.New("tmpl").Funcs(template.FuncMap{
		"month_year": monthYear,
	}).Parse(tmpl)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)

	if err = t.Execute(buf, r); err != nil {
		return "", err
	}

	return buf.String(), nil
}
