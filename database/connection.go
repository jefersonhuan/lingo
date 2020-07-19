package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

func (server *Server) Connect() (err error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(server.URI).SetMaxPoolSize(0))
	if err != nil {
		return
	}

	ctx, _ := context.WithTimeout(context.TODO(), 60*time.Second)

	err = client.Connect(ctx)
	if err != nil {
		return
	}

	server.Client = client

	return
}

func (server *Server) Ping() (err error) {
	err = server.Connect()
	if err != nil {
		return
	}

	defer server.Disconnect()

	return server.Client.Ping(context.TODO(), readpref.Primary())
}

func (server *Server) Disconnect() {
	err := server.Client.Disconnect(context.TODO())
	if err != nil {
		fmt.Printf("an error occurred while disconnecting the server: %v\n", err)
	}
}
