package main

import (
	"flag"
	"net/http"
	"os"

	_ "net/http/pprof"

	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// PORT is a port which the server will listen on
var PORT string

func init() {
	flag.StringVar(&PORT, "p", "3228", "server port")
}

func main() {
	// TODO: change to normal log
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// parse all flags set in `init`
	flag.Parse()

	log.Info().Str("port", PORT).Msg("starting server")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "HEAD", "POST", "PUT", "DELETE", "PATCH"},
	})
	h := c.Handler(handler{})

	http.ListenAndServe(":"+PORT, h)
}
