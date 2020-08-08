package main

import (
	"flag"
	"net/http"
	"os"
	"time"

	_ "net/http/pprof"

	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var port string
var prettyPrint bool

func initFlags() {
	flag.StringVar(&port, "p", "3228", "server port")
	flag.BoolVar(&prettyPrint, "pp", false, "enable pretty print")

	flag.Parse()

	envPort := os.Getenv("PORT")
	if envPort != "" {
		port = envPort
		log.Info().Str("port", port).Msg("using env port")
	}

	if prettyPrint {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.Stamp})
	}
}

func main() {
	initFlags()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "HEAD", "POST", "PUT", "DELETE", "PATCH"},
	})
	h := c.Handler(handler{})

	log.Info().Str("port", port).Msg("starting server")

	err := http.ListenAndServe(":"+port, h)
	if err != nil {
		log.Err(err).Send()
		os.Exit(1)
	}
}
