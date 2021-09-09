package core

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ahmadwaleed/choreui/app/config"
	"github.com/ahmadwaleed/choreui/app/i18n"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type Template struct {
	Extension string
	Folder    string
	LayoutDir string
	Data      map[string]interface{}
}

func NewTemplate(config *config.AppConfig) *Template {
	return &Template{
		Extension: config.TemplateExt,
		Folder:    config.TemplateFolder,
		LayoutDir: config.TemplateLayoutDir,
	}
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmplname := strings.Split(name, ".")
	if len(tmplname) != 2 {
		return fmt.Errorf("could not parse given template name")
	}

	root, tmpl := tmplname[0], tmplname[1]

	layouts, err := filepath.Glob(
		t.Folder + string(os.PathSeparator) + t.LayoutDir + string(os.PathSeparator) + "*." + t.Extension,
	)
	if err != nil {
		return fmt.Errorf("could not load template layouts: %v", err)
	}

	var includes []string
	includes = append(includes, layouts...)
	includes = append(includes, t.Folder+string(os.PathSeparator)+tmpl+"."+t.Extension)

	template, err := template.New(name).Funcs(template.FuncMap{
		"Lang": i18n.Get,
	}).ParseFiles(includes...)
	if err != nil {
		fmt.Println(includes)
		return fmt.Errorf("could not parse template files: %v", err)
	}

	if err := t.LoadSessionData(c); err != nil {
		return fmt.Errorf("could not load session data: %v", err)
	}

	return template.ExecuteTemplate(w, root+"."+t.Extension, data)
}

func (t *Template) LoadSessionData(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return err
	}

	for key, val := range sess.Values {
		if k, ok := key.(string); ok {
			t.Data[k] = val
		}
	}

	return nil
}
