package context

import (
	"github.com/ahmadwaleed/choreui/app/config"
	"github.com/ahmadwaleed/choreui/app/i18n"
	"github.com/ahmadwaleed/choreui/app/models"
	"github.com/labstack/echo/v4"
)

// AppContext is the new context in the request / response cycle
// We can use the db store, cache and central configuration
type AppContext struct {
	echo.Context
	UserStore models.UserModel
	Config    *config.AppConfig
	Loc       i18n.I18ner
}
