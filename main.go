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

func init() {
	flag.StringVar(&port, "p", "3228", "server port")
	flag.BoolVar(&prettyPrint, "pp", false, "enable pretty print")

	// parse all flags set in `init`
	flag.Parse()

	if prettyPrint {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.Stamp})
	}
}

func main() {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "HEAD", "POST", "PUT", "DELETE", "PATCH"},
	})
	h := c.Handler(handler{})

	log.Info().Str("port", port).Msg("starting server")
	http.ListenAndServe(":"+port, h)
}
