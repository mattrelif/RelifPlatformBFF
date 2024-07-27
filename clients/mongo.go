package clients

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func NewMongoClient(uri string, connectionTimeout time.Duration) (*mongo.Client, error) {
	var client *mongo.Client

	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()

	opts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, opts)

	if err != nil {
		return nil, err
	}

	return client, nil
}

func DisconnectMongoClient(client *mongo.Client) error {
	return client.Disconnect(context.TODO())
}
