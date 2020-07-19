package transfer

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"mongo-transfer/database"
	"sync"
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
	mutex   *sync.Mutex
	docs    [][]bson.M
}

func (buffer *CollectionBuffer) flush(page int, wg *sync.WaitGroup, failures *[]error) {
	var docs = make([]interface{}, len(buffer.docs[page]))

	for index, doc := range buffer.docs[page] {
		docs[index] = doc
	}

	opts := options.InsertMany().SetOrdered(true)
	if _, err := buffer.handler.InsertMany(context.TODO(), docs, opts); err != nil {
		*failures = append(*failures, err)
	}

	wg.Done()
}
