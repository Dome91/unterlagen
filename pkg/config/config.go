package config

import (
	"github.com/gorilla/securecookie"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"os"
	"sync"
)

type Config struct {
	Development  bool
	E2E          bool
	Port         string
	CookieSecret []byte
}

var loadConfig func() Config

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	loadConfig = sync.OnceValue(func() Config {
		viper.AutomaticEnv()
		viper.SetEnvPrefix("unterlagen")
		viper.SetDefault("development", false)
		viper.SetDefault("e2e", false)
		viper.SetDefault("port", "8080")
		viper.SetDefault("cookie_secret", securecookie.GenerateRandomKey(32))
		development := viper.GetBool("development")
		e2e := viper.GetBool("e2e")
		port := viper.GetString("port")
		cookieSecret := viper.GetString("cookie_secret")
		config := Config{
			Development:  development,
			E2E:          e2e,
			Port:         port,
			CookieSecret: []byte(cookieSecret),
		}

		if development {
			log.Info().Msg("Started in development mode")
			err := os.MkdirAll(".ws", 0755)
			if err != nil {
				panic(err)
			}
			log.Info().Msg("Created development workspace")
		}

		return config
	})
}

func Get() Config {
	return loadConfig()
}
func Overwrite(key string, value string) {
	viper.Set(key, value)
}
