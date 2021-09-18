package view

import (
	"html/template"

	"github.com/ahmadwaleed/choreui/app/core/session"
	"github.com/ahmadwaleed/choreui/app/database"
	"github.com/ahmadwaleed/choreui/app/i18n"
	"github.com/ahmadwaleed/choreui/app/utils"
	"github.com/labstack/echo/v4"
)

func I18nPlugin() template.FuncMap {
	return template.FuncMap{
		"Lang": i18n.Get,
	}
}

func RoutePlugin(c echo.Context) template.FuncMap {
	return template.FuncMap{
		"route": func(name string) string {
			return utils.Route(c, name)
		},
	}
}

func SessionPlugin(c echo.Context) template.FuncMap {
	sess := session.NewSessionStore(c)
	return template.FuncMap{
		"Auth": func() bool {
			return sess.GetBool("Auth")
		},
		"User": func() database.User {
			return sess.Get("User").(database.User)
		},
		"Old": func(name string) string {
			data, err := c.FormParams()
			if err != nil {
				c.Logger().Error(err)
			}
			if val, ok := data[name]; ok {
				if len(val) == 0 {
					return ""
				}
				return val[0]
			}

			return ""
		},
	}
}
