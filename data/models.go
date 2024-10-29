package data

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type LogEntry struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type LogStorage interface {
	Insert(entry LogEntry) error
	All() ([]*LogEntry, error)
	GetOne(id string) (*LogEntry, error)
	DropCollection() error
	Update() (*mongo.UpdateResult, error)
}

func NewMongologStorage(client *mongo.Client) (LogStorage, error) {

	logStorage := &Mongo{
		Client:   client,
		LogEntry: LogEntry{},
	}
	return logStorage, nil
}
