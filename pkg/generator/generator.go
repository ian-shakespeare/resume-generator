package generator

import (
	"bytes"
	"html/template"
	"resumegenerator/pkg/resume"
)

func GenerateHtml(r *resume.Resume, tmpl string) (string, error) {
	t, err := template.New("tmpl").Parse(tmpl)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)

	if err = t.Execute(buf, r); err != nil {
		return "", err
	}

	return buf.String(), nil
}
