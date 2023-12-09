// https://stackoverflow.com/questions/38804313/build-docker-image-from-go-code
package main

import (
	"flag"
	"fmt"
	"os"

	"battleground-server/internal/manager"

	"github.com/rs/zerolog"
)

type config struct {
	port int
	env  string
}

type application struct {
	config  config
	logger  zerolog.Logger
	manager manager.Manager
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "mode", "development", "Mode (development|testing)")
	flag.Parse()

	logFile, err := os.OpenFile(
		"/home/bo/Documents/battleground/server/logs/battleground_logs.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	logger := zerolog.New(logFile).With().Timestamp().Logger()

	o, err := manager.NewManager(logger)
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
