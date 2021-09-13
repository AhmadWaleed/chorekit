package view

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ahmadwaleed/choreui/app/config"
	"github.com/ahmadwaleed/choreui/app/core/session"
	"github.com/labstack/echo/v4"
)

type Template struct {
	Extension string
	Folder    string
	LayoutDir string
	Vars      map[string]interface{}
}

func NewTemplate(config *config.AppConfig) *Template {
	return &Template{
		Extension: config.TemplateExt,
		Folder:    config.TemplateFolder,
		LayoutDir: config.TemplateLayoutDir,
		Vars:      make(map[string]interface{}),
	}
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	var root = "base"
	if strings.Contains(name, "auth") {
		root = "auth"
	}

	layouts, err := filepath.Glob(
		t.Folder + string(os.PathSeparator) + t.LayoutDir + string(os.PathSeparator) + "*." + t.Extension,
	)

	var includes []string
	includes = append(includes, layouts...)
	includes = append(includes, t.Folder+string(os.PathSeparator)+name+"."+t.Extension)

	templates, err := template.New(name).Funcs(I18nPlugin()).Funcs(SessionPlugin(c)).ParseFiles(includes...)
	if err != nil {
		return fmt.Errorf("could not parse template files: %v", err)
	}

	sess := session.NewSessionStoreFunc(c)
	flashes := sess.Session.Flashes()
	t.Vars["flashes"] = make([]session.Flash, len(flashes))
	if len(flashes) > 0 {
		for i, f := range flashes {
			switch f.(type) {
			case session.Flash:
				t.Vars["flashes"].([]session.Flash)[i] = f.(session.Flash)
			default:
				t.Vars["flashes"].([]session.Flash)[i] = session.Flash{Message: f.(string), Type: session.FlashInfo}
			}
		}
	}

	t.Vars["data"] = data

	return templates.ExecuteTemplate(w, root+"."+t.Extension, t.Vars)
}
