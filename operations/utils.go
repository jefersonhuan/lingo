package operations

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"mongo-transfer/models"
	"time"
)

func TestConnection(server models.Server) (err error) {
	fmt.Println("Connecting to", server.ToURI())

	client, err := mongo.NewClient(options.Client().ApplyURI(server.ToURI()))
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return
	}

	return client.Ping(ctx, readpref.Primary())
}
