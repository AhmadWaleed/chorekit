package view

import (
	"html/template"

	"github.com/ahmadwaleed/choreui/app/i18n"
)

func I18nPlugin() template.FuncMap {
	return template.FuncMap{
		"Lang": i18n.Get,
	}
}
