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
	client, err := mongo.NewClient(options.Client().ApplyURI(server.URI))
	if err != nil {
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	err = client.Connect(ctx)
	if err != nil {
		return
	}

	server.client = client
	server.ctx = ctx

	return
}

func (server Server) Ping() (err error) {
	fmt.Println("Connecting to", server.URI)

	err = server.Connect()
	if err != nil {
		return
	}

	defer server.Disconnect()

	return server.client.Ping(server.ctx, readpref.Primary())
}

func (server *Server) Disconnect() {
	err := server.client.Disconnect(server.ctx)
	if err != nil {
		fmt.Printf("an error occurred while disconnecting the server: %v\n", err)
	}
}
