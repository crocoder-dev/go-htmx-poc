package main

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

type App struct {
	htmx         *htmx.HTMX
	appTemplates *Template
}

type Page struct {
	Title   string
	Boosted bool
	Chart   template.HTML
	Chart2  template.HTML
}

type SettingsGlobal struct {
	Name     string
	File     *multipart.File
	Dropdown string
}

type ApiResponse struct {
	UserId    int    `json:"userId"`
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
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

func (a *App) Fetch(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)

	os.Remove("db/sqlite-database.db")

	fmt.Println("Creating sqlite-database.db...")
	file, err := os.Create("db/sqlite-database.db")
	if err != nil {
		fmt.Println(err.Error())
	}
	file.Close()

	sqliteDatabase, _ := sql.Open("sqlite3", "./db/sqlite-database.db")
	defer sqliteDatabase.Close()
	createTable(sqliteDatabase)

	insertStudent(sqliteDatabase, "0001", "Liana Kim", "Bachelor")
	insertStudent(sqliteDatabase, "0002", "Glen Rangel", "Bachelor")
	insertStudent(sqliteDatabase, "0003", "Martin Martins", "Master")
	insertStudent(sqliteDatabase, "0004", "Alayna Armitage", "PHD")
	insertStudent(sqliteDatabase, "0005", "Marni Benson", "Bachelor")
	insertStudent(sqliteDatabase, "0006", "Derrick Griffiths", "Master")
	insertStudent(sqliteDatabase, "0007", "Leigh Daly", "Bachelor")
	insertStudent(sqliteDatabase, "0008", "Marni Benson", "PHD")
	insertStudent(sqliteDatabase, "0009", "Klay Correa", "Bachelor")

	page := Page{Title: "Fetch", Boosted: h.HxBoosted}

	if page.Boosted == true {
		return c.Render(http.StatusOK, "fetch", &page)
	}

	return c.Render(http.StatusOK, "fetch.html", &page)
}

func (a *App) Test(c echo.Context) error {
	return c.Render(http.StatusOK, "test", Page{Title: "Test"})
}

func (a *App) Submit(c echo.Context) (err error) {
	name := c.FormValue("name")
	email := c.FormValue("email")
	fmt.Println("Name: ", name)
	fmt.Println("Email: ", email)
	return c.String(http.StatusOK, "Submitted!")
}

func (a *App) Chart(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)

	chart := template.HTML(CreateTestChart())
	chat2 := template.HTML(CreateTestChart2())
	page := Page{Title: "Chart", Boosted: h.HxBoosted, Chart: chart, Chart2: chat2}

	if page.Boosted == true {
		return c.Render(http.StatusOK, "chart", page)
	}
	return c.Render(http.StatusOK, "chart.html", page)
}

func (a *App) getData(c echo.Context) error {
	sqliteDatabase, _ := sql.Open("sqlite3", "./db/sqlite-database.db")
	defer sqliteDatabase.Close()

	var data = getStudents(sqliteDatabase)
	return c.HTML(http.StatusOK, data)
}

func (a *App) setSettings(c echo.Context) (err error) {
	err = c.Request().ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	name := c.FormValue("name")
	dropdown := c.FormValue("dropdown")
	fileHeader, err := c.FormFile("file")
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	file, err := fileHeader.Open()
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	defer file.Close()

	settingsGlobal = SettingsGlobal{
		Name:     name,
		File:     &file,
		Dropdown: dropdown,
	}
	fmt.Println(settingsGlobal.Name)
	fmt.Println(settingsGlobal.Dropdown)
	fmt.Println(settingsGlobal.File)
	return c.String(http.StatusOK, fileHeader.Filename)
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
	e.GET("/fetch", app.Fetch)

	e.POST("/submit", app.Submit)
	e.POST("/setSettings", app.setSettings)
	e.GET("/getData", app.getData)
	e.Static("/", "dist")

	e.Logger.Fatal(e.Start(":3000"))
}

func createTable(db *sql.DB) {
	createStudentTableSQL := `CREATE TABLE student (
		"idStudent" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"code" TEXT,
		"name" TEXT,
		"program" TEXT		
	  );` // SQL Statement for Create Table

	log.Println("Create student table...")
	statement, err := db.Prepare(createStudentTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements
	log.Println("student table created")
}

// We are passing db reference connection from main to our method with other parameters
func insertStudent(db *sql.DB, code string, name string, program string) {
	log.Println("Inserting student record ...")
	insertStudentSQL := `INSERT INTO student(code, name, program) VALUES (?, ?, ?)`
	statement, err := db.Prepare(insertStudentSQL) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(code, name, program)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func getStudents(db *sql.DB) string {
	row, err := db.Query("SELECT * FROM student ORDER BY name")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	var optionList string
	for row.Next() {
		var id int
		var code string
		var name string
		var program string
		row.Scan(&id, &code, &name, &program)
		optionList = optionList + "<option>" + name + "</option>"
		log.Println("Student: ", code, " ", name, " ", program)
	}
	return optionList
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
