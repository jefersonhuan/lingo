package database

import (
	"go.mongodb.org/mongo-driver/bson"
	"mongo-transfer/utils"
)

func (server *Server) LoadAll() error {
	return utils.StepsFunctions(server.LoadDatabases, server.LoadCollections)
}

func (server *Server) LoadDatabases() (err error) {
	result, err := server.Client.ListDatabases(server.Ctx, bson.D{})
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
		cols, err := server.Client.Database(database.Specification.Name).ListCollectionNames(server.Ctx, bson.D{})
		if err != nil {
			continue
		}

		server.Databases[index].Collections = cols
	}

	return nil
}
