package main

import (
	"context"
	"html/template"
	"net/http"
	"fmt"
	"mime/multipart"

	"github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

type App struct {
	htmx         *htmx.HTMX
	appTemplates *Template
}

type Page struct {
	Title   string
	Boosted bool
	Chart   template.HTML
}

type SettingsGlobal struct {
    Name string
	File *multipart.FileHeader
	Dropdown string
}

var settingsGlobal SettingsGlobal

func (a *App) Index(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)

	page := Page{Title: "Index", Boosted: h.HxBoosted}

	if page.Boosted == true {
		return c.Render(http.StatusOK, "index", &page)
	}

	return c.Render(http.StatusOK, "index.html", &page)
}

func (a *App) About(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)

	page := Page{Title: "About", Boosted: h.HxBoosted}

	if page.Boosted == true {
		return c.Render(http.StatusOK, "about", &page)
	}

	return c.Render(http.StatusOK, "about.html", &page)
}

func (a *App) Contact(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)

	page := Page{Title: "Contact", Boosted: h.HxBoosted}

	if page.Boosted == true {
		return c.Render(http.StatusOK, "contact", &page)
	}

	return c.Render(http.StatusOK, "contact.html", &page)
}

func (a *App) Settings(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)

	page := Page{Title: "Settings", Boosted: h.HxBoosted}

	if page.Boosted == true {
		return c.Render(http.StatusOK, "settings", &page)
	}

	return c.Render(http.StatusOK, "settings.html", &page)
}

func (a *App) Test(c echo.Context) error {
	return c.Render(http.StatusOK, "test", Page{Title: "Test"})
}

func (a *App) Submit(c echo.Context) (err error) {
    name := c.FormValue("name")
    email := c.FormValue("email")
	fmt.Println("Name: ", name);
	fmt.Println("Email: ", email);
	return c.String(http.StatusOK, "Submitted!")
}

func (a *App) Chart(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)

	chart := template.HTML(CreateLineChart())
	page := Page{Title: "Chart", Boosted: h.HxBoosted, Chart: chart}

  if page.Boosted == true {
    return c.Render(http.StatusOK, "chart", page)
  }
	return c.Render(http.StatusOK, "chart.html", page)
}
func (a *App) setSettings(c echo.Context) (err error) {
    err = c.Request().ParseMultipartForm(10 << 20) // 10 MB
    if err != nil {
        fmt.Println("Error:", err)
    }
    name := c.FormValue("name")
    dropdown := c.FormValue("dropdown")
    fileHeader, err := c.FormFile("file")
    if err != nil {
        fmt.Println("Error:", err)
    }

    file, err := fileHeader.Open()
    if err != nil {
        fmt.Println("Error:", err)
    }
    defer file.Close()

    settingsGlobal = SettingsGlobal{
        Name:     name,
        File:     fileHeader,
        Dropdown: dropdown,
    }
    fmt.Println(settingsGlobal.Name)
    fmt.Println(settingsGlobal.Dropdown)
    fmt.Println(settingsGlobal.File)
	return c.String(http.StatusOK, "Submitted!")
}

func main() {
	e := echo.New()
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(HtmxMiddleware)
	
	app := &App{
		appTemplates: new(Template),
	}

	app.appTemplates.Init()
	app.appTemplates.Add("templates/*.html")
	app.appTemplates.Add("templates/*/*.html")

	e.Renderer = app.appTemplates

	e.GET("/", app.Index)
	e.GET("/about", app.About)
	e.GET("/contact", app.Contact)
	e.GET("/settings", app.Settings)
	e.GET("/test", app.Test)
	e.GET("/chart", app.Chart)

	e.POST("/submit", app.Submit)
	e.POST("/setSettings", app.setSettings)
	e.Static("/", "dist")

	e.Logger.Fatal(e.Start(":3000"))
}

func HtmxMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		hxh := htmx.HxRequestHeader{
			HxBoosted:               htmx.HxStrToBool(c.Request().Header.Get("HX-Boosted")),
			HxCurrentURL:            c.Request().Header.Get("HX-Current-URL"),
			HxHistoryRestoreRequest: htmx.HxStrToBool(c.Request().Header.Get("HX-History-Restore-Request")),
			HxPrompt:                c.Request().Header.Get("HX-Prompt"),
			HxRequest:               htmx.HxStrToBool(c.Request().Header.Get("HX-Request")),
			HxTarget:                c.Request().Header.Get("HX-Target"),
			HxTriggerName:           c.Request().Header.Get("HX-Trigger-Name"),
			HxTrigger:               c.Request().Header.Get("HX-Trigger"),
		}

		ctx = context.WithValue(ctx, htmx.ContextRequestHeader, hxh)

		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}
