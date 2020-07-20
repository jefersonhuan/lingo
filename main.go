package main

import (
	"lingo/cmd"
	"lingo/store"
)

func main() {
	store.InitDatabase()

	cmd.Execute()
}
