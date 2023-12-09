package main

import (
	job_handler "battleground-engine/job-handler"
	"context"
	"fmt"
	"net"
	"os"

	"battleground-engine/worker"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type myJobServer struct {
	job_handler.UnimplementedJobServer
	worker worker.Worker
	logger zerolog.Logger
}

func NewJobServer(logger zerolog.Logger) *myJobServer {
	return &myJobServer{
		worker: *worker.NewWorker(1),
		logger: logger,
	}
}

func (s *myJobServer) Create(ctx context.Context, req *job_handler.CreateRequest) (*job_handler.CreateResponse, error) {
	s.logger.Info().Msg("Received requested")
	err := s.worker.WriteProgramToFile([]byte(req.SourceCode), "test.py")
	if err != nil {
		s.logger.Fatal().Err(err).Msg("Error writing program to file")
	}
	input, err := s.worker.WriteSolutionInput(req.Input)
	if err != nil {
		s.logger.Fatal().Err(err).Msg("Error writing input")
	}
	res, err := s.worker.ExecuteSolution(input, "test.py")
	if err != nil {
		s.logger.Warn().Err(err).Msg("Error executing")
	}

	return &job_handler.CreateResponse{
		Stdout: res.Stdout.String(),
		Stderr: res.Stderr.String(),
		Error:  fmt.Sprint(err),
	}, nil
}

func main() {

	logFile, err := os.OpenFile(
		"logs/engine_logs.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	logger := zerolog.New(logFile).With().Timestamp().Logger()

	lis, err := net.Listen("tcp", ":8089")
	if err != nil {
		logger.Fatal().Err(err).Msg("Cannot create listener")
	}
	serverRegistrar := grpc.NewServer()
	service := &myJobServer{
		logger: logger,
	}

	job_handler.RegisterJobServer(serverRegistrar, service)
	err = serverRegistrar.Serve(lis)
	if err != nil {
		logger.Fatal().Err(err).Msg("Cannot serve")
	}
}
