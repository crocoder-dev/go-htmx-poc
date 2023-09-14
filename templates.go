package main

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Init() {
  t.templates = template.New("")
}

func (t *Template) Add(glob string) {
  t.templates = template.Must(t.templates.ParseGlob(glob))
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
