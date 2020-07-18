package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"mongo-transfer/store"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	ID        string
	URI       string
	CreatedAt time.Time

	Databases []Database

	ctx    context.Context
	client *mongo.Client
}

type Database struct {
	Name        string
	Collections []string
}

const serversCollection = "servers"

func (database Database) String() (content string) {
	divider := strings.Repeat("-", 40)

	content = "Database: " + database.Name + "\n"

	if len(database.Collections) != 0 {
		content += "\nCollections:\n"

		for index, col := range database.Collections {
			content += col + "\t"

			if len(database.Collections)-1 == index {
				content += "\n" + divider + "\n"
			}
		}
	}

	return
}

func (server *Server) FromURI(username, password, host string, port int) {
	prefix := "mongodb://"
	suffix := host + ":" + strconv.Itoa(port)

	var auth string

	if username != "" && password != "" {
		auth = username + ":" + password + "@"
	}

	server.URI = prefix + auth + suffix
}

func (server *Server) Load(ID string) (err error) {
	if err = store.Conn.Read(serversCollection, ID, &server); err != nil {
		return
	}

	if err = server.Connect(); err != nil {
		return
	}

	return
}

func (server Server) Save() error {
	if err := store.Conn.Write(serversCollection, server.ID, server); err != nil {
		return err
	}

	return nil
}

func (server Server) Delete() error {
	if err := store.Conn.Delete(serversCollection, server.ID); err != nil {
		return err
	}

	return nil
}
