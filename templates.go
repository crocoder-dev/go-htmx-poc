package main

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

func NewTemplate(glob string) *Template {
	return &Template{
		templates: template.Must(template.ParseGlob(glob)),
	}
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl := template.Must(t.templates.Clone())
	tmpl = template.Must(tmpl.Parse("templates/" + name))
	return t.templates.ExecuteTemplate(w, name, data)
}
