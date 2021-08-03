package mongo

import (
	"context"
	"sync"
	"ws/model"

	"go.mongodb.org/mongo-driver/mongo"
	mstore "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var clientError error
var mongoOnce sync.Once

const (
	CONNECTIONSTRING = "mongodb://localhost:27017"
	DbName           = "images"
	DbCollection     = "details"
)

type Mongo struct {
	Client *mongo.Client
}

//GetClient - Return mongodb connection
func NewMongo() (*Mongo, error) {
	mdb := Mongo{}
	//Do once!
	mongoOnce.Do(func() {
		// Set client options
		co := options.Client().ApplyURI(CONNECTIONSTRING)
		// Connection stage
		c, err := mstore.Connect(context.TODO(), co)
		if err != nil {
			clientError = err
		}
		// Verification stage
		err = c.Ping(context.TODO(), nil)
		if err != nil {
			clientError = err
		}
		client = c
	})

	mdb.Client = client
	return &mdb, clientError
}

// StoreFileMetaData store the image metadata into db
func (m *Mongo) StoreFileMetaData(ctx context.Context, metadata model.Metadata) error {

	collection := m.Client.Database(DbName).Collection(DbCollection)

	_, err := collection.InsertOne(ctx, metadata)
	if err != nil {
		return err
	}

	return nil
}
