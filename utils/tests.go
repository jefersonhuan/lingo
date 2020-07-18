package utils

import (
	"fmt"
	"github.com/viant/endly"
	"github.com/viant/endly/system/docker"
	"github.com/viant/toolbox"
)

var endlyManager = endly.New()
var endlyContext = endlyManager.NewContext(toolbox.NewContext())

func StartMongo() error {
	_, err := endlyManager.Run(endlyContext, &docker.RunRequest{
		Image: "mongo:latest",
		Ports: map[string]string{
			"27017": "27017",
		},
		Name: "mongo",
	})
	return err
}

func StopMongo() {
	_, err := endlyManager.Run(endlyContext, &docker.StopRequest{
		Name: "mongo",
	})

	if err != nil {
		fmt.Println("An error occurred while stopping MongoDB", err)
	}
}
