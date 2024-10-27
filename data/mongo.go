package data

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/dubass83/go-micro-logger/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type Mongo struct {
	Client   *mongo.Client
	LogEntry LogEntry
}

func MongoConnect(conf util.Config) (*mongo.Client, error) {
	// Create a Client to a MongoDB server and use Ping to verify that the
	// server is running.

	clientOpts := options.Client().ApplyURI(conf.MongoURL)
	client, err := mongo.Connect(clientOpts)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatal().Err(err)
		}
	}()

	// Call Ping to verify that the deployment is up and the Client was
	// configured successfully.
	for i := range 10 {
		if err = client.Ping(ctx, readpref.Primary()); err != nil {
			log.Warn().Err(err).Msgf("%d try to connect to mongodb from 10", i+1)
			time.Sleep(time.Second * 2)
			continue
		}
		return client, nil
	}
	return nil, err
}

func (m *Mongo) Insert(entry LogEntry) error {
	collection := m.Client.Database("logs").Collection("logs")

	_, err := collection.InsertOne(context.TODO(), LogEntry{
		Name:      entry.Name,
		Data:      entry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Error().Err(err).Msg("Error inserting into logs")
		return err
	}

	return nil
}

func (m *Mongo) All() ([]*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := m.Client.Database("logs").Collection("logs")

	opts := options.Find()
	opts.SetSort(bson.D{{"created_at", -1}})

	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Error().Err(err).Msg("Finding all docs error")
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []*LogEntry

	for cursor.Next(ctx) {
		var item LogEntry
		err := cursor.Decode(&item)
		if err != nil {
			log.Error().Err(err).Msg("Error decoding item into LogEntry struct")
			return nil, err
		}
		logs = append(logs, &item)
	}
	return logs, nil
}

func (m *Mongo) GetOne(id string) (*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := m.Client.Database("logs").Collection("logs")

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Error().Err(err).Msg("can not get object id from string")
		return nil, err
	}
	var entry LogEntry
	err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&entry)
	if err != nil {
		log.Error().Err(err).Msg("failed to find log by ID")
		return nil, err
	}
	return &entry, nil
}

func (m *Mongo) DropCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := m.Client.Database("logs").Collection("logs")

	err := collection.Drop(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed drop the collection")
		return err
	}
	return nil
}

func (m *Mongo) Update() (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := m.Client.Database("logs").Collection("logs")

	docID, err := primitive.ObjectIDFromHex(m.LogEntry.ID)
	if err != nil {
		log.Error().Err(err).Msg("can not get object id from string")
	}

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": docID},
		bson.D{
			{"$set", bson.D{
				{"name", m.LogEntry.Name},
				{"data", m.LogEntry.Data},
				{"updated_at", time.Now()},
			}},
		},
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to update log entry from the struct")
		return nil, err
	}

	return result, nil
}
