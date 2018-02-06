package main

import (
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"snippetbox.org/pkg/models"
	"bytes"
	"github.com/justinas/nosurf"
)

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

type HTMLData struct {
	CSRFToken string
	Flash string
	Form interface{}
	LoggedIn bool
	Path string
	Snippet *models.Snippet
	Snippets []*models.Snippet
}

func (app *App) RenderHTML(w http.ResponseWriter, r *http.Request, page string, data *HTMLData) {
	if data == nil {
		data = &HTMLData{}
	}

	data.Path = r.URL.Path

	// Always add the CSRF token to the data for our templates.
	data.CSRFToken = nosurf.Token(r)

	var err error
	data.LoggedIn, err = app.LoggedIn(r)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	files := []string{
		filepath.Join(app.HTMLDir, "base.html"),
		filepath.Join(app.HTMLDir, page),
	}

	funcs := template.FuncMap{
		"humanDate": humanDate,
	}

	ts, err := template.New("").Funcs(funcs).ParseFiles(files...)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	buf := new(bytes.Buffer)
	err = ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	buf.WriteTo(w)
}
