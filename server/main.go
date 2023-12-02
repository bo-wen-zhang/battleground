// https://stackoverflow.com/questions/38804313/build-docker-image-from-go-code
package main

import (
	"fmt"

	"battleground-server/server"
)

func main() {

	o, err := server.NewOrchestrator()
	if err != nil {
		return
	}

	dockerFileName := "Dockerfile"
	contextDirSrc := "../engine"
	buildOutput, err := o.BuildImage(dockerFileName, contextDirSrc, "battleground-engine")
	if err != nil {
		return
	}
	fmt.Print(buildOutput)

	err = o.CreateWorker("battleground-engine")
	if err != nil {
		return
	}

	err = o.ContainerLogs()
	if err != nil {
		return
	}

	o.RemoveWorker(o.WorkerIDs[0])
}
