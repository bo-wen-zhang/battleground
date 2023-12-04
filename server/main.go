// https://stackoverflow.com/questions/38804313/build-docker-image-from-go-code
package main

import (
	"flag"
	"fmt"

	"battleground-server/server"
)

func main() {
	var cfg server.Config

	flag.StringVar(&cfg.Mode, "mode", "development", "Mode (development|testing)")
	flag.Parse()

	o, err := server.NewOrchestrator(cfg)
	if err != nil {
		return
	}

	imageName := "battleground-engine"
	dockerFileName := "Dockerfile"
	contextDirSrc := "../engine"
	buildOutput, err := o.BuildImage(dockerFileName, contextDirSrc, imageName)
	if err != nil {
		return
	}
	fmt.Print(buildOutput)

	err = o.CreateWorker(imageName)
	if err != nil {
		return
	}

	err = o.ContainerLogs()
	if err != nil {
		return
	}

	o.RemoveWorker(o.WorkerIDs[0])
}
