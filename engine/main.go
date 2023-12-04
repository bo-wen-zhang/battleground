package main

import (
	job_handler "battleground-engine/job-handler"
	"context"
	"fmt"
	"log"
	"net"

	"battleground-engine/worker"

	"google.golang.org/grpc"
)

type myJobServer struct {
	job_handler.UnimplementedJobServer
	worker worker.Worker
}

func NewJobServer() *myJobServer {
	return &myJobServer{
		worker: *worker.NewWorker(1),
	}
}

func (s *myJobServer) Create(ctx context.Context, req *job_handler.CreateRequest) (*job_handler.CreateResponse, error) {

	err := s.worker.WriteProgramToFile([]byte(req.SourceCode), "test.py")
	if err != nil {
		log.Fatal("Error writing program to file.")
	}
	input, err := s.worker.WriteSolutionInput(req.Input)
	if err != nil {
		log.Fatal("Error writing input.")
	}
	res, err := s.worker.ExecuteSolution(input, "test.py")
	if err != nil {
		log.Println("Error executing program:", res)
	}

	return &job_handler.CreateResponse{
		Stdout: res.Stdout.String(),
		Stderr: res.Stderr.String(),
		Error:  fmt.Sprint(err),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8089")
	if err != nil {
		log.Fatal("Cannot create listener:", err)
	}
	serverRegistrar := grpc.NewServer()
	service := &myJobServer{}

	job_handler.RegisterJobServer(serverRegistrar, service)
	err = serverRegistrar.Serve(lis)
	if err != nil {
		log.Fatal("Cannot serve:", err)
	}
	log.Println("Serving...")
}
