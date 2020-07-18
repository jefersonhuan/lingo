package main

import (
	"mongo-transfer/cmd"
	"mongo-transfer/store"
)

func main() {
	store.InitDatabase()

	cmd.Execute()
}
