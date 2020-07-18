package database

import (
	"go.mongodb.org/mongo-driver/bson"
	"mongo-transfer/utils"
)

func (server *Server) LoadAll() error {
	return utils.StepsFunctions(server.LoadDatabases, server.LoadCollections)
}

func (server *Server) LoadDatabases() (err error) {
	result, err := server.client.ListDatabaseNames(server.ctx, bson.D{})
	if err != nil {
		return
	}

	for _, database := range result {
		server.Databases = append(server.Databases, Database{
			Name: database,
		})
	}

	return
}

func (server *Server) LoadCollections() error {
	for index, database := range server.Databases {
		cols, err := server.client.Database(database.Name).ListCollectionNames(server.ctx, bson.D{})
		if err != nil {
			continue
		}

		server.Databases[index].Collections = cols
	}

	return nil
}
