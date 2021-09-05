package core

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"path/filepath"
	"strings"

	"github.com/ahmadwaleed/choreui/app/i18n"
	"github.com/labstack/echo/v4"
)

var mainTmpl = `{{define "main" }} {{ template "base" . }} {{ end }}`

type templateRenderer struct {
	templates map[string]*template.Template
}

// NewTemplateRenderer creates a new setup to render layout based go templates
func newTemplateRenderer(layoutsDir, templatesDir string) *templateRenderer {
	r := &templateRenderer{}
	r.templates = make(map[string]*template.Template)
	r.Load(layoutsDir, templatesDir)
	return r
}

func (t *templateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmplname := strings.Split(name, ".")
	if len(tmplname) != 2 {
		return fmt.Errorf("could not parse given template")
	}

	layout, include := tmplname[0], fmt.Sprintf("%s.tmpl", tmplname[1])

	tmpl, ok := t.templates[include]
	if !ok {
		c.Logger().Fatalf("could not found template: %s", include)
		return fmt.Errorf("could not found template: %s", include)
	}

	return tmpl.ExecuteTemplate(w, layout, data)
}

func (t *templateRenderer) Load(layoutsDir, templatesDir string) {
	layouts, err := filepath.Glob(layoutsDir)
	if err != nil {
		log.Fatal(err)
	}

	includes, err := glob(templatesDir)
	if err != nil {
		log.Fatal(err)
	}

	funcMap := template.FuncMap{
		"Lang": i18n.Get,
	}

	mainTemplate := template.New("main")
	mainTemplate.Funcs(funcMap)

	mainTemplate, err = mainTemplate.Parse(mainTmpl)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range includes {
		fileName := filepath.Base(file)
		files := append(layouts, file)
		t.templates[fileName], err = mainTemplate.Clone()

		if err != nil {
			log.Fatal(err)
		}

		t.templates[fileName] = template.Must(t.templates[fileName].ParseFiles(files...))
	}
}

func glob(dir string) ([]string, error) {
	root, pattern := filepath.Split(dir)

	var files []string
	err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			fls, err := filepath.Glob(fmt.Sprintf("%s/%s", path, pattern))
			if err != nil {
				return err
			}

			files = append(files, fls...)
		}

		return nil
	})

	return files, err
}
