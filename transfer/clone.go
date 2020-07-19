package transfer

import (
	"context"
	"fmt"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"mongo-transfer/utils"
	"strings"
	"sync"
)

const barTitleWidth = 45

var wg sync.WaitGroup

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
		for _, coll := range db.Collections {
			sourceCollection := source.Client.Database(db.Specification.Name).Collection(coll)
			targetCollection := target.Client.Database(db.Specification.Name).Collection(coll)

			totalDocs, err := sourceCollection.CountDocuments(context.TODO(), bson.D{})
			if totalDocs == 0 {
				continue
			}

			var limit int64 = 3000
			var nPages int64

			if err != nil {
				return err
			} else if totalDocs < limit {
				nPages = 1
			} else {
				nPages = totalDocs/limit + 1
			}

			wg.Add(int(nPages))
			buffer := make(chan []interface{}, nPages)

			bar := startBarForCollection(db.Specification.Name+"."+coll, nPages, p)

			go startBuffer(targetCollection, buffer, bar)

			if err = stepCloning(sourceCollection, buffer, nPages, limit); err != nil {
				mes := fmt.Errorf("an error occurred while cloning collection %s: %w", coll, err)
				fmt.Println(utils.ColorfulString("red", mes.Error()))

				continue
			}

			wg.Wait()
			close(buffer)

			bar.Completed()
		}
	}

	return
}

func stepCloning(source *mongo.Collection, buffer chan []interface{}, nPages, limit int64) (err error) {
	var page int64

	for page = 0; page < nPages; page++ {
		opts := options.Find().SetLimit(limit).SetSkip(page * limit)
		cursor, err := source.Find(context.TODO(), bson.D{}, opts)
		if err != nil {
			break
		}

		go storeQueryResults(cursor, buffer)
	}

	return
}

func storeQueryResults(cursor *mongo.Cursor, buffer chan []interface{}) (error, []interface{}) {
	var results []bson.M

	if err := cursor.All(context.TODO(), &results); err != nil {
		fmt.Println(err)
	}

	docs := make([]interface{}, len(results))

	for index, result := range results {
		docs[index] = result
	}

	go func() {
		buffer <- docs
	}()

	return nil, docs
}

func startBuffer(target *mongo.Collection, docs chan []interface{}, bar *mpb.Bar) {
	for data := range docs {
		saveDocs(target, data)
		bar.Increment()

		wg.Done()
	}
}

func saveDocs(target *mongo.Collection, docs []interface{}) {
	opts := options.InsertMany().SetOrdered(true)
	if _, err := target.InsertMany(context.TODO(), docs, opts); err != nil {
		//fmt.Println(err)
	}
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
