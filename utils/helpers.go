package utils

import (
	"fmt"
	"github.com/viant/endly"
	"github.com/viant/endly/system/docker"
	"github.com/viant/toolbox"
)

func StepsFunctions(functions ...func() error) error {
	for _, f := range functions {
		err := f()
		if err != nil {
			return err
		}
	}

	return nil
}

var colors = map[string]string{
	"reset": "\033[0m",

	"red":    "\033[31m",
	"green":  "\033[32m",
	"yellow": "\033[33m",
	"blue":   "\033[34m",
	"purple": "\033[35m",
	"cyan":   "\033[36m",
	"white":  "\033[37m",
}

func ColorfulString(color string, content string) string {
	code, ok := colors[color]

	if !ok {
		return content
	} else {
		return code + string(content) + colors["reset"]
	}
}

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
