package core

import (
	"github.com/ahmadwaleed/choreui/app/database"
	"github.com/ahmadwaleed/choreui/app/i18n"
	"github.com/labstack/echo/v4"
)

// AppContext is the new context in the request / response cycle
// We can use the db store, cache and central configuration
type AppContext struct {
	echo.Context
	App   *Application
	Loc   i18n.I18ner
	Store *database.Store
}
