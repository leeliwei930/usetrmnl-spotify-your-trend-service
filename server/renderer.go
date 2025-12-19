package server

import (
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
)

type WebTemplate struct {
	templates map[string]*template.Template
}

func NewWebTemplate() *WebTemplate {
	base := template.New("")

	// 1. Parse layouts/partials into base
	err := filepath.Walk("views/layouts", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".html") {
			rel, err := filepath.Rel("views", path)
			if err != nil {
				return err
			}
			name := filepath.ToSlash(rel)
			b, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			_, err = base.New(name).Parse(string(b))
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}

	registry := make(map[string]*template.Template)

	// 2. Parse pages (cloning base)
	err = filepath.Walk("views", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		// Skip layouts directory from being registered as pages
		rel, _ := filepath.Rel("views", path)
		if strings.HasPrefix(rel, "layouts") {
			return nil
		}

		if strings.HasSuffix(path, ".html") {
			name := filepath.ToSlash(rel)

			// Clone base
			tmpl, err := base.Clone()
			if err != nil {
				return err
			}

			b, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			// Parse the page into the clone
			_, err = tmpl.New(name).Parse(string(b))
			if err != nil {
				return err
			}

			registry[name] = tmpl
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	return &WebTemplate{
		templates: registry,
	}
}

func (t *WebTemplate) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, ok := t.templates[name]
	if !ok {
		return echo.NewHTTPError(404, "Template not found: "+name)
	}
	return tmpl.ExecuteTemplate(w, name, data)
}
