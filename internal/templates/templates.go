package templates

import (
	"html/template"
	"io"
)

type Templates struct {
	templates *template.Template
}

func New() *Templates {
	t := template.Must(template.ParseGlob("web/views/*.html"))
	return &Templates{
		templates: t,
	}
}

func (t *Templates) Render(w io.Writer, name string, data interface{}) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
