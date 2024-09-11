package generator

import (
	"bytes"
	"html/template"
)

func (r *resumeData) GenerateHtml(tmplStr string) (string, error) {
	tmpl, err := template.New("tmpl").Parse(tmplStr)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)

	if err = tmpl.Execute(buf, r); err != nil {
		return "", err
	}

	return buf.String(), nil
}
