package template

import (
	"errors"
	"html/template"
	"io"
	"path/filepath"

	"alvintanoto.id/blog/pkg/forms"
	"github.com/labstack/echo/v4"
)

type Template struct {
	Templates map[string]*template.Template
}

type TemplateData struct {
	Form       *forms.Form
	Flash      string
	FlashError string
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, ok := t.Templates[name]
	if !ok {
		return errors.New("template not found")
	}
	return tmpl.ExecuteTemplate(w, "base", data)
}

func NewTemplateCache(dir string) map[string]*template.Template {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.html"))
	if err != nil {
		return nil
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return nil
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.html"))
		if err != nil {
			return nil
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.html"))
		if err != nil {
			return nil
		}

		cache[name] = ts
	}

	return cache
}
