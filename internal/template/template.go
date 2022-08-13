package template

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"path/filepath"

	model "alvintanoto.id/blog/internal/model/database"
	"alvintanoto.id/blog/pkg/forms"
	"github.com/labstack/echo/v4"
)

type Template struct {
	Templates map[string]*template.Template
}

type TemplateData struct {
	Form              *forms.Form
	Flash             string
	FlashError        string
	AuthenticatedUser *model.User
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	fmt.Println("Rendering ", name)

	tmpl, ok := t.Templates[name]
	if !ok {
		return errors.New("template not found")
	}

	if data != nil {
		user := c.Get("user")
		td, ok := data.(*TemplateData)

		if user != nil {
			td.AuthenticatedUser = user.(*model.User)
		}

		if !ok {
			return errors.New("failed to add data")
		}

		fmt.Println("data", user)
		fmt.Println("templateData", td)

		return tmpl.ExecuteTemplate(w, "base", td)
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
