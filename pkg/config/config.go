package config

import (
	"github.com/gorilla/securecookie"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

var Development = os.Getenv("DEVELOPMENT") != ""
var Port = os.Getenv("PORT")
var CookieSecret = []byte(os.Getenv("COOKIE_SECRET"))

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if Port == "" {
		Port = "8080"
	}

	if len(CookieSecret) == 0 {
		CookieSecret = securecookie.GenerateRandomKey(32)
	}

	if Development {
		log.Info().Msg("Started in development mode")
		err := os.MkdirAll(".ws", 0755)
		if err != nil {
			panic(err)
		}
		log.Info().Msg("Created development workspace")
	}

}
