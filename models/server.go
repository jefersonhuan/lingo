package models

import (
	"fmt"
	"mongo-transfer/store"
	"time"
)

type Server struct {
	ID        string
	Host      string
	Port      int
	User      string
	Password  string
	CreatedAt time.Time
}

const serversCollection = "servers"

func (server Server) ToURI() string {
	uri := fmt.Sprintf("mongodb://%s:%d", server.Host, server.Port)

	return uri
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
