package transfer

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"lingo/database"
	"time"
)

type Transfer struct {
	Source *database.Server
	Target *database.Server

	StartedAt  time.Time
	FinishedAt time.Time
}

type CollectionBuffer struct {
	handler *mongo.Collection
	docs    [][]bson.M

	size, avgObjSize float64
}

func (buffer *CollectionBuffer) flush(page int) {
	var docs = make([]interface{}, len(buffer.docs[page]))

	for index, doc := range buffer.docs[page] {
		docs[index] = doc
	}

	opts := options.InsertMany().SetOrdered(true)
	if _, err := buffer.handler.InsertMany(context.TODO(), docs, opts); err != nil {
		pushError(err)
	}

	buffering.Done()
}
