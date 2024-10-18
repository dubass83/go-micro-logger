package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dubass83/go-micro-logger/cmd/api"
	"github.com/dubass83/go-micro-logger/util"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// var interaptSignals = []os.Signal{
// 	os.Interrupt,
// 	syscall.SIGTERM,
// 	syscall.SIGINT,
// }

func main() {
	conf, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("cannot load configuration")
	}
	if conf.Enviroment == "devel" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		// log.Debug().Msgf("config values: %+v", conf)
	}

	// ctx, stop := signal.NotifyContext(context.Background(), interaptSignals...)
	// defer stop()

	runChiServer(conf)
}

// runChiServer run http server with Chi framework
func runChiServer(conf util.Config) {
	server := api.CreateNewServer(conf)

	server.ConfigureCORS()
	server.AddMiddleware()
	server.MountHandlers()
	log.Info().
		Msgf("start listening on the port %s\n", server.Config.WebPort)
	HTTPAddressString := fmt.Sprintf(":%s", server.Config.WebPort)
	err := http.ListenAndServe(HTTPAddressString, server.Router)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("method", "main").
			Msg("can not start server")
	}
}
