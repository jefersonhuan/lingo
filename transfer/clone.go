package transfer

import (
	"context"
	"fmt"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"mongo-transfer/database"
	"mongo-transfer/utils"
	"strings"
	"sync"
	"time"
)

const barTitleWidth = 45

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
				sourceCollection := source.Client.Database(dbName).Collection(coll)

				buffers[index].handler = target.Client.Database(dbName).Collection(coll)

				if err = stepCloning(sourceCollection, &buffers[index], p); err != nil {
					mes := fmt.Errorf("an error occurred while cloning collection %s: %w", coll, err)
					pushError(mes)

					continue
				}
			}
		}(db)
	}

	wg.Wait()

	fmt.Println(utils.ColorfulString("cyan", "Finishing sync"))

	buffering.Wait()

	if len(failures) != 0 {
		fmt.Println(utils.ColorfulString("yellow", "\nThe following errors occurred:"))

		for _, failure := range failures {
			fmt.Println(utils.ColorfulString("red", failure.Error()))
		}
	}

	return
}

func stepCloning(source *mongo.Collection, buffer *CollectionBuffer, p *mpb.Progress) (err error) {
	var limit int64 = 4000
	var nPages, page int64

	if err = fetchPageCount(source, &nPages, limit); err != nil {
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
				pushError(err)
				cancel()
			}

			if err := cursor.All(context.TODO(), &buffer.docs[page]); err != nil {
				pushError(err)
				cancel()
			}

			buffer.flush(int(page))
			bar.Increment()

			cancel()
		}
	}()

	wg.Done()

	return
}

func fetchPageCount(coll *mongo.Collection, nPages *int64, limit int64) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	totalDocs, err := coll.CountDocuments(ctx, bson.D{})

	if err != nil {
		return err
	} else if totalDocs < limit {
		*nPages = 1
	} else {
		*nPages = totalDocs/limit + 1
	}

	return nil
}

func startBarForCollection(name string, total int64, p *mpb.Progress) *mpb.Bar {
	if len(name) < barTitleWidth {
		name += strings.Repeat(" ", barTitleWidth-len(name))
	}

	return p.AddBar(total,
		mpb.PrependDecorators(
			decor.Name(name),
			decor.Percentage(decor.WCSyncSpace),
		),
		mpb.AppendDecorators(
			decor.OnComplete(
				decor.EwmaETA(decor.ET_STYLE_GO, 60), "finished",
			),
		),
	)
}

func pushError(err error) {
	failures = append(failures, err)
}
