package transfer

import (
	"context"
	"fmt"
	"github.com/vbauerster/mpb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"lingo/database"
	"lingo/utils"
	"sync"
	"time"
)

const barTitleWidth = 45
const paginationSize = 16 * 1000 * 1000 // kb

var wg, buffering sync.WaitGroup
var failures []error

func (transfer *Transfer) clone() (err error) {
	source := transfer.Source
	target := transfer.Target

	p := mpb.New(mpb.WithWidth(64))

	fmt.Println("\nCloning databases from", source.ID)

	err = source.LoadAll()
	if err != nil {
		return
	}

	for _, db := range source.Databases {
		dbName := db.Specification.Name

		buffers := make([]CollectionBuffer, len(db.Collections))
		wg.Add(len(db.Collections))

		go func(db database.Database) {
			for index, coll := range db.Collections {

				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				result := source.Client.Database(dbName).RunCommand(ctx, bson.M{"collStats": coll})

				var stats bson.M
				err = result.Decode(&stats)

				if err == nil {
					getStat(&buffers[index].size, stats["size"])
					getStat(&buffers[index].avgObjSize, stats["avgObjSize"])
				}

				sourceCollection := source.Client.Database(dbName).Collection(coll)
				buffers[index].handler = target.Client.Database(dbName).Collection(coll)

				if err = stepCloning(sourceCollection, &buffers[index], p); err != nil {
					mes := fmt.Errorf("an error occurred while cloning collection %s: %w", coll, err)
					pushError(mes)

					continue
				}

				cancel()
			}
		}(db)
	}

	wg.Wait()
	buffering.Wait()

	if len(failures) != 0 {
		fmt.Println(utils.ColorfulString("yellow", "\nThe following errors occurred:"))

		for _, failure := range failures {
			fmt.Println(utils.ColorfulString("red", failure.Error()))
		}
	}

	return nil
}

func stepCloning(source *mongo.Collection, buffer *CollectionBuffer, p *mpb.Progress) (err error) {
	var limit int64 = 4000
	var nPages, page int64

	if err = fetchPageCount(source, *buffer, &nPages, &limit); err != nil {
		wg.Done()
		return
	}

	bar := startBarForCollection(source.Database().Name()+"."+source.Name(), nPages, p)

	buffer.docs = make([][]bson.M, nPages)
	buffering.Add(int(nPages))

	go func() {
		for page = 0; page < nPages; page++ {
			ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
			opts := options.Find().SetLimit(limit).SetSkip(page * limit)

			cursor, err := source.Find(ctx, bson.D{}, opts)
			if err != nil {
				buffering.Done()
				pushError(err)
				continue
			}

			if err := cursor.All(context.TODO(), &buffer.docs[page]); err != nil {
				buffering.Done()
				pushError(err)
				continue
			}

			buffer.flush(int(page))
			bar.Increment()

			cancel()
		}
	}()

	wg.Done()

	return
}

func fetchPageCount(coll *mongo.Collection, buffer CollectionBuffer, nPages, limit *int64) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	totalDocs, err := coll.CountDocuments(ctx, bson.D{})
	if totalDocs == 0 {
		return fmt.Errorf("%s has no readable documents", coll.Name())
	}

	if totalDocs > *limit && buffer.size != 0 && buffer.avgObjSize != 0 {
		*limit = paginationSize / int64(buffer.avgObjSize)
	}

	if err != nil {
		return err
	} else if totalDocs < *limit {
		*nPages = 1
	} else {
		*nPages = totalDocs/(*limit) + 1
	}

	return ctx.Err()
}
