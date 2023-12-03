package main

import (
	job_handler "battleground-engine/job-handler"
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

func writeProgramToFile(program, filename string) error {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("File does not exist or cannot be created")
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	fmt.Fprintf(w, "%v\n", program)

	w.Flush()
	return nil
}

type myJobServer struct {
	job_handler.UnimplementedJobServer
}

func (s *myJobServer) Create(context.Context, *job_handler.CreateRequest) (*job_handler.CreateResponse, error) {
	return &job_handler.CreateResponse{
		Stdout: "",
		Stderr: "",
		Error:  "",
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
}
