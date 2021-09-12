package view

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ahmadwaleed/choreui/app/config"
	"github.com/ahmadwaleed/choreui/app/database"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

var (
	// FlashError is a tailwindcss class
	FlashError = "error"
	// FlashSuccess is a tailwindcss class
	FlashSuccess = "success"
	// FlashInfo is a tailwindcss class
	FlashInfo = "info"
	// FlashWarning is a tailwindcss class
	FlashWarning = "warning"
)

type Flash struct {
	Message string
	Type    string
}

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

	templates, err := template.New(name).Funcs(I18nPlugin()).ParseFiles(includes...)
	if err != nil {
		return fmt.Errorf("could not parse template files: %v", err)
	}

	sess, err := session.Get("session", c)
	if err != nil {
		return err
	}

	if sess.Values["auth"] != nil && sess.Values["auth"].(bool) == true {
		t.Vars["auth"] = true
		if user, ok := sess.Values["user"].(database.User); ok {
			t.Vars["user"] = user
		}
	} else {
		t.Vars["auth"] = true
	}

	flashes := sess.Flashes()
	t.Vars["flashes"] = make([]Flash, len(flashes))
	if len(flashes) > 0 {
		for i, f := range flashes {
			switch f.(type) {
			case Flash:
				t.Vars["flashes"].([]Flash)[i] = f.(Flash)
			default:
				t.Vars["flashes"].([]Flash)[i] = Flash{f.(string), FlashInfo}
			}
		}
	}

	params, err := c.FormParams()
	if err != nil {
		return fmt.Errorf("could not get form params: %v", err)
	}
	formData := make(map[string]string, len(params))
	for key, val := range params {
		formData[key] = val[0]
	}

	t.Vars["data"] = data
	t.Vars["form"] = formData

	return templates.ExecuteTemplate(w, root+"."+t.Extension, t.Vars)
}
