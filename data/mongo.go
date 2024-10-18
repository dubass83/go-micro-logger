package data

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/dubass83/go-micro-logger/util"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func New(conf util.Config) (*mongo.Client, error) {
	// Create a Client to a MongoDB server and use Ping to verify that the
	// server is running.

	clientOpts := options.Client().ApplyURI(conf.MongoURL)
	client, err := mongo.Connect(clientOpts)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Fatal().Err(err)
		}
	}()

	// Call Ping to verify that the deployment is up and the Client was
	// configured successfully.
	for i := range 10 {
		if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
			log.Warn().Err(err).Msgf("%d try to connect to mongodb from 10", i+1)
			time.Sleep(time.Second * 2)
			continue
		}
		return client, nil
	}
	return nil, err
}
