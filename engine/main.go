package main

import (
	pb "battleground-engine/engine_service"
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"battleground-engine/worker"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type engineServer struct {
	pb.UnimplementedEngineServiceServer
	worker worker.Worker
	logger zerolog.Logger
}

func NewJobServer(logger zerolog.Logger) *engineServer {
	return &engineServer{
		worker: *worker.NewWorker(1),
		logger: logger,
	}
}

func (s *engineServer) GetProgramResult(ctx context.Context, req *pb.Program) (*pb.Result, error) {
	s.logger.Info().Msg("Received requested")
	err := s.worker.WriteProgramToFile([]byte(req.SourceCode), "test.py")
	if err != nil {
		s.logger.Fatal().Err(err).Msg("Error writing program to file")
	}
	var givenInput string
	if req.Input != nil {
		givenInput = *req.Input
	}
	input, err := s.worker.WriteSolutionInput(givenInput)
	if err != nil {
		s.logger.Fatal().Err(err).Msg("Error writing input")
	}
	res, err := s.worker.ExecuteSolution(input, "test.py")
	if err != nil {
		s.logger.Warn().Err(err).Msg("Error executing")
	}

	return &pb.Result{
		StandardOutput: res.Stdout.String(),
		StandardError:  res.Stderr.String(),
		ElapsedTime:    "",
		MemoryUsage:    "",
		EngineError:    fmt.Sprint(err),
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
	logger.Info().Msg("This engine has started running.")

	lis, err := net.Listen("tcp4", ":8089")
	if err != nil {
		logger.Fatal().Err(err).Msg("Cannot create listener")
	}
	serverRegistrar := grpc.NewServer()
	service := &engineServer{
		logger: logger,
	}

	go func() {
		pb.RegisterEngineServiceServer(serverRegistrar, service)
		err = serverRegistrar.Serve(lis)
		if err != nil {
			logger.Fatal().Err(err).Msg("Cannot serve")
		}
	}()

	for {
		time.Sleep(2 * time.Second)
		fmt.Println("lets gooooo")
	}
}
