package generator

import (
	"bytes"
	"fmt"
	"html/template"
	"resumegenerator/pkg/resume"
	"strconv"
	"time"
)

func year(t *time.Time) string {
	if t == nil {
		return ""
	}
	return strconv.Itoa(t.Year())
}

func monthYear(t *time.Time) string {
	if t == nil {
		return ""
	}
	return fmt.Sprintf("%s %d", t.Month().String(), t.Year())
}

func GenerateHtml(r *resume.Resume, tmpl string) (string, error) {
	t, err := template.New("tmpl").Funcs(template.FuncMap{
		"year":       year,
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
