package main

import (
	"context"
	"net/http"

	"github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

type State struct {
}

type App struct {
	htmx         *htmx.HTMX
	appTemplates *Template

	state State
}

type Page struct {
	Title string
}

func (a *App) Index(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)

	if h.HxBoosted == true {
		return c.Render(http.StatusOK, "index", Page{Title: "Index Boost"})
	}

	return c.Render(http.StatusOK, "index.html", Page{Title: "Index"})
}

func (a *App) About(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)

	if h.HxBoosted == true {
		return c.Render(http.StatusOK, "about", Page{Title: "About Boost"})
	}

	return c.Render(http.StatusOK, "about.html", Page{Title: "About"})
}

func main() {
	e := echo.New()
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(HtmxMiddleware)

	app := &App{
		appTemplates: new(Template),
		state:        State{},
	}

	app.appTemplates.Add("templates/*.html")

	e.Renderer = app.appTemplates

	e.GET("/", app.Index)
	e.GET("/about", app.About)

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
