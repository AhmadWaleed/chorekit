package config

import (
	"log"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	Address          string `env:"ADDRESS" envDefault:":8080"`
	AssetsBuildDir   string `env:"ASSETS_BUILD_DIR"`
	TemplateDir      string `env:"TPL_DIR"`
	LayoutDir        string `env:"LAYOUT_DIR"`
	ConnectionString string `env:"CONNECTION_STRING,required"`
	IsProduction     bool   `env:"PRODUCTION"`
	GrayLogAddr      string `env:"GRAYLOG_ADDR"`
	RequestLogger    bool   `env:"REQUEST_LOGGER"`
	LocaleDir        string `env:"LOCALE_DIR" envDefault:"locales"`
	Lang             string `env:"LANG" envDefault:"en_US"`
	LangDomain       string `env:"LANG_DOMAIN" envDefault:"default"`
	AppKey           string `env:"APP_KEY,required"`
}

func NewConfig(files ...string) (*AppConfig, error) {
	err := godotenv.Load(files...)

	if err != nil {
		log.Printf("could not load app config: %q, %v\n", files, err)
	}

	cfg := AppConfig{}
	err = env.Parse(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
