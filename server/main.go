// https://stackoverflow.com/questions/38804313/build-docker-image-from-go-code
package main

import (
	"flag"
	"fmt"
	"os"

	"battleground-server/server"

	"github.com/rs/zerolog"
)

func main() {
	var cfg server.Config

	flag.StringVar(&cfg.Mode, "mode", "development", "Mode (development|testing)")
	flag.Parse()

	logFile, err := os.OpenFile(
		"battleground_logs.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	logger := zerolog.New(logFile).With().Timestamp().Logger()

	o, err := server.NewOrchestrator(cfg, logger)
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
