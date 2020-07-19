package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"mongo-transfer/utils"
)

func (server *Server) LoadAll() error {
	return utils.StepsFunctions(server.LoadDatabases, server.LoadCollections)
}

func (server *Server) LoadDatabases() (err error) {
	result, err := server.Client.ListDatabases(context.TODO(), bson.D{})
	if err != nil {
		return
	}

	server.Databases = make([]Database, len(result.Databases))

	for index, database := range result.Databases {
		server.Databases[index] = Database{
			Specification: database,
		}
	}

	return
}

func (server *Server) LoadCollections() error {
	for index, database := range server.Databases {
		cols, err := server.Client.Database(database.Specification.Name).ListCollectionNames(context.TODO(), bson.D{})
		if err != nil {
			continue
		}

		server.Databases[index].Collections = cols
	}

	return nil
}
