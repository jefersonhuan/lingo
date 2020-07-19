package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"mongo-transfer/store"
	"mongo-transfer/utils"
	"strconv"
	"time"
)

type Server struct {
	ID        string
	URI       string
	CreatedAt time.Time

	Databases []Database

	Ctx    context.Context
	Client *mongo.Client
}

type Database struct {
	Specification mongo.DatabaseSpecification
	Collections   []string
}

const serversCollection = "servers"

func (database Database) String() (content string) {
	content = utils.ColorfulString("blue", "Database: ")
	content += database.Specification.Name + "\n"
	content += "SizeOnDisk: " + strconv.FormatInt(database.Specification.SizeOnDisk, 10) + " kb\n"

	if len(database.Collections) != 0 {
		content += utils.ColorfulString("cyan", "Collections:\n")

		for _, col := range database.Collections {
			content += "- " + col + "\n"
		}
	}

	return
}

func (server *Server) BuildURI(username, password, host string, port int) {
	prefix := "mongodb://"
	suffix := host + ":" + strconv.Itoa(port)

	var auth string

	if username != "" && password != "" {
		auth = username + ":" + password + "@"
	}

	server.URI = prefix + auth + suffix
}

func (server *Server) Fetch() (err error) {
	if err = store.Conn.Read(serversCollection, server.ID, &server); err != nil {
		return
	}

	return server.Connect()
}

func (server Server) Save() error {
	return store.Conn.Write(serversCollection, server.ID, server)
}

func (server Server) Delete() error {
	return store.Conn.Delete(serversCollection, server.ID)
}
