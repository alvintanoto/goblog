package template

import (
	"errors"
	"html/template"
	"io"
	"path/filepath"

	model "alvintanoto.id/blog/internal/model/database"
	"alvintanoto.id/blog/pkg/forms"
	"alvintanoto.id/blog/pkg/log"
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
	log.Get().InfoLog.Println("Rendering ", name)

	tmpl, ok := t.Templates[name]
	if !ok {
		log.Get().ErrorLog.Println("Render error, template not found:", name)
		return errors.New("template not found")
	}

	if data != nil {
		td, ok := data.(*TemplateData)
		if !ok {
			log.Get().ErrorLog.Println("Render error, failed to parse data:", name)
			return errors.New("failed parse data")
		}

		td = addTemplateData(td, c)

		return tmpl.ExecuteTemplate(w, "base", td)
	}

	return tmpl.ExecuteTemplate(w, "base", data)
}

func addTemplateData(td *TemplateData, c echo.Context) *TemplateData {
	user := c.Get("user")

	if user != nil {
		td.AuthenticatedUser = user.(*model.User)
	}

	return td
}

func NewTemplateCache(dir string) map[string]*template.Template {
	log.Get().InfoLog.Println("Start caching template...")

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

	log.Get().InfoLog.Println("Caching template done...")
	return cache
}
