package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const timeout = 10 * time.Second

// NewClient established connection to a mongoDb instance using provided URI and auth credentials.
func NewClient(uri, username, password string) (*mongo.Client, error) {
	opts := options.Client().ApplyURI(uri)
	if username != "" && password != "" {
		opts.SetAuth(options.Credential{
			Username: username, Password: password,
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}


