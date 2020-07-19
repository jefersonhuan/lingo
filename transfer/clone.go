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
)

const barTitleWidth = 45

var wg, buffering sync.WaitGroup

func (transfer *Transfer) clone() (err error) {
	source := transfer.Source
	target := transfer.Target

	var failures []error

	p := mpb.New(mpb.WithWidth(64))

	fmt.Println("\nCloning databases from", source.ID)

	err = source.LoadAll()
	if err != nil {
		return
	}

	for _, db := range source.Databases {
		buffers := make([]CollectionBuffer, len(db.Collections))
		wg.Add(len(db.Collections))

		go func(db database.Database) {
			for index, coll := range db.Collections {
				sourceCollection := source.Client.Database(db.Specification.Name).Collection(coll)

				buffers[index].handler = target.Client.Database(db.Specification.Name).Collection(coll)
				buffers[index].mutex = &sync.Mutex{}

				if err = stepCloning(sourceCollection, &buffers[index], &failures, p); err != nil {
					mes := fmt.Errorf("an error occurred while cloning collection %s: %w", coll, err)
					failures = append(failures, mes)

					continue
				}
			}
		}(db)
	}

	wg.Wait()

	fmt.Println(utils.ColorfulString("cyan", "Finishing sync"))

	buffering.Wait()

	if len(failures) != 0 {
		fmt.Println(utils.ColorfulString("yellow", "The following errors occurred:"))

		for _, failure := range failures {
			fmt.Println(utils.ColorfulString("red", failure.Error()))
		}
	}

	return
}

func stepCloning(source *mongo.Collection, buffer *CollectionBuffer, failures *[]error, p *mpb.Progress) (err error) {
	totalDocs, err := source.CountDocuments(context.TODO(), bson.D{})

	var limit int64 = 2500
	var nPages int64

	if err != nil {
		fmt.Println(err)
	} else if totalDocs < limit {
		nPages = 1
	} else {
		nPages = totalDocs/limit + 1
	}

	bar := startBarForCollection(source.Database().Name()+"."+source.Name(), nPages, p)

	var page int64

	buffer.docs = make([][]bson.M, nPages)
	buffering.Add(int(nPages))

	for page = 0; page < nPages; page++ {
		opts := options.Find().SetLimit(limit).SetSkip(page * limit)
		cursor, err := source.Find(context.TODO(), bson.D{}, opts)
		if err != nil {
			break
		}

		if err := cursor.All(context.TODO(), &buffer.docs[page]); err != nil {
			fmt.Println(err)
		}

		go func(page int) {
			buffer.flush(page, &buffering, failures)
		}(int(page))

		bar.Increment()
	}

	wg.Done()

	bar.Completed()

	return
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
