package main

import (
	"embed"
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

//go:embed templates/*
var views embed.FS

func (t *Template) Init() {
	t.templates = template.New("")
}

func (t *Template) Add(glob string) {
	t.templates = template.Must(t.templates.ParseFS(views, glob))
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
