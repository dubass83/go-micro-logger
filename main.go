package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"time"

	"github.com/dubass83/go-micro-logger/cmd/api"
	"github.com/dubass83/go-micro-logger/cmd/gapi"
	"github.com/dubass83/go-micro-logger/data"
	"github.com/dubass83/go-micro-logger/pb"
	"github.com/dubass83/go-micro-logger/util"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"google.golang.org/grpc"
)

// var interaptSignals = []os.Signal{
// 	os.Interrupt,
// 	syscall.SIGTERM,
// 	syscall.SIGINT,
// }

func main() {
	// load config from file and env variables
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
	// create mongo client
	clientOpts := options.Client().ApplyURI(conf.MongoURL)
	client, err := mongo.Connect(clientOpts)
	if err != nil {
		log.Fatal().Err(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatal().Err(err)
		}
	}()

	logStorage, err := data.NewMongologStorage(client)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create log storage")
	}

	// create new RPC service
	rpcService := api.CreateNewRPCService(conf, logStorage)
	err = rpc.Register(rpcService)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to register RPC service")
	}

	go runRPCService(rpcService)

	// create new gRPC service
	gRPCService := gapi.CreateNewLogServer(conf, logStorage)

	go runGRPCService(gRPCService)

	// ctx, stop := signal.NotifyContext(context.Background(), interaptSignals...)
	// defer stop()

	server := api.CreateNewServer(conf, logStorage)

	runChiServer(server)
}

// runChiServer run http server with Chi framework
func runChiServer(server *api.Server) {

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

// runRPCService starts an RPC service on the specified port from the RPCService configuration.
// It listens for incoming TCP connections and serves them using the rpc.ServeConn method.
// If there is an error while listening on the port or accepting a connection, it logs the error.
//
// Parameters:
//   - rpcs: A pointer to an api.RPCService which contains the configuration for the RPC service.
//
// Returns:
//   - error: An error if the service fails to start or encounters an issue while running.
func runRPCService(rpcs *api.RPCService) {
	log.Info().Msgf("starting RPC service on port %s", rpcs.Config.RPCPort)
	listen, err := net.Listen("tcp", ":"+rpcs.Config.RPCPort)
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to listen on port %s", rpcs.Config.RPCPort)
		return
	}
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Error().Err(err).Msg("failed to accept connection")
			continue
		}
		go rpc.ServeConn(conn)
	}
}

// runGRPCService starts a gRPC service using the provided LogServer configuration.
// It listens on the port specified in the LogServer configuration, registers the
// LogService server, and serves incoming gRPC requests.
//
// Parameters:
//
//	gRPC - A pointer to a LogServer instance containing the gRPC configuration.
//
// The function logs the start of the gRPC service, and in case of any errors during
// listening or serving, it logs the error and terminates the process.
func runGRPCService(gRPC *gapi.LogServer) {
	log.Info().Msgf("starting gRPC service on port %s", gRPC.Config.GRPCPort)
	listen, err := net.Listen("tcp", ":"+gRPC.Config.GRPCPort)
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to listen on port %s", gRPC.Config.GRPCPort)
		return
	}
	defer listen.Close()

	s := grpc.NewServer()
	pb.RegisterLogServiceServer(s, gRPC)
	log.Info().Msg("gRPC server started and registered LogService")
	if err := s.Serve(listen); err != nil {
		log.Fatal().Err(err).Msg("failed to serve gRPC")
		return
	}
}
