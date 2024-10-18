package main

import (
	"os"

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

	// connPool, err := pgxpool.NewWithConfig(context.Background(), poolConfig(conf))
	// if err != nil {
	// 	log.Fatal().
	// 		Err(err).
	// 		Msg("cannot validate db connection string")
	// }
	// store := data.NewStore(connPool)

	// runChiServer(conf, store)
}

// runChiServer run http server with Chi framework
// func runChiServer(conf util.Config, store data.Store) {
// 	server := api.CreateNewServer(conf, store)

// 	server.ConfigureCORS()
// 	server.AddMiddleware()
// 	server.MountHandlers()
// 	log.Info().
// 		Msgf("start listening on the port %s\n", server.Config.HTTPAddressString)
// 	err := http.ListenAndServe(server.Config.HTTPAddressString, server.Router)
// 	if err != nil {
// 		log.Fatal().
// 			Err(err).
// 			Str("method", "main").
// 			Msg("can not start server")
// 	}
// }
