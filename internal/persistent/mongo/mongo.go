package mongo

import (
	"context"
	"sync"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	singletonDBInstance *mongo.Client
	mu                  sync.Mutex
)

type DBConnection struct {
	*mongo.Client
}

func NewMongoDBConnection(uri string) (*DBConnection, error) {
	mu.Lock()
	defer mu.Unlock()

	if singletonDBInstance != nil {
		return &DBConnection{
			Client: singletonDBInstance,
		}, nil
	}

	client, err := mongo.Connect(nil, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	singletonDBInstance = client

	return &DBConnection{
		Client: singletonDBInstance,
	}, nil
}

func (db *DBConnection) Close(ctx context.Context) error {
	mu.Lock()
	defer mu.Unlock()

	if singletonDBInstance != nil {
		if err := singletonDBInstance.Disconnect(ctx); err != nil {
			return err
		}
		singletonDBInstance = nil
	}
	return nil
}

func (db *DBConnection) GetDBInstance() *mongo.Client {
	return singletonDBInstance
}
