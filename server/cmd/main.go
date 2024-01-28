// https://stackoverflow.com/questions/38804313/build-docker-image-from-go-code
package main

import (
	"flag"
	"os"
	"sync"

	"battleground-server/internal/manager"

	"github.com/docker/docker/client"
	"github.com/rs/zerolog"
)

type config struct {
	port int
	env  string
}

const (
	imageName      = "battleground-engine"
	dockerFileName = "Dockerfile"
	contextDirSrc  = "../engine"
)

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

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logger.Fatal().Err(err).Msg("Unable to init client")
	}

	man, err := manager.NewManager(imageName, cli, logger, 1, []string{"8089"})
	if err != nil {
		return
	}

	man.BuildImage(dockerFileName, contextDirSrc)

	defer func() {
		if r := recover(); r != nil {
			man.RemoveAllContainers()
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		man.MaintainEngines()

	}()

	wg.Wait()

	man.RemoveAllContainers()
}
