// https://stackoverflow.com/questions/38804313/build-docker-image-from-go-code
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sync"

	pb "battleground-server/engine_service"
	"battleground-server/internal/manager"

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

	man, err := manager.NewManager(imageName, logger)
	if err != nil {
		return
	}

	man.BuildImage(dockerFileName, contextDirSrc)
	hostPort := "8089"

	err = man.BuildEngine(hostPort)
	if err != nil {
		logger.Error().Err(err).Msg("Could not build engine.")
		return
	}

	defer func() {
		if r := recover(); r != nil {
			man.RemoveAllContainers()
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		man.ContainerLogs()
	}()

	res, err := man.Engines[0].Stub.GetProgramResult(context.Background(), &pb.Program{
		UserId:     69,
		SourceCode: "print(\"Hello World\")",
	})
	if err != nil {
		logger.Panic().Err(err).Msg("Error getting job response")
	}
	fmt.Println(res)

	wg.Wait()

	man.RemoveAllContainers()
}
