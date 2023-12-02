package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"battleground-server/internal/util"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

type Orchestrator struct {
	Images    []string //names of images
	WorkerIDs []string //ids of worker containers
	client    *client.Client
	ctx       context.Context
}

func NewOrchestrator() (*Orchestrator, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err, " :unable to init client")
		return &Orchestrator{}, err
	}
	return &Orchestrator{
		client: cli,
		ctx:    context.Background(),
	}, nil
}

func (o *Orchestrator) BuildImage(dockerFileName, contextDirSrc, imageName string) (string, error) {
	dockerFileTarReader, err := util.CreateTarReader(contextDirSrc)
	if err != nil {
		log.Print(err, " :unable to create tar reader")
		return "", err
	}

	imageBuildResponse, err := o.client.ImageBuild(
		o.ctx,
		dockerFileTarReader,
		types.ImageBuildOptions{
			Context:    dockerFileTarReader,
			Dockerfile: dockerFileName,
			Tags:       []string{imageName}})
	if err != nil {
		log.Print(err, " :unable to build docker image")
		return "", err
	}
	defer imageBuildResponse.Body.Close()
	_, err = io.Copy(os.Stdout, imageBuildResponse.Body)
	if err != nil {
		log.Print(err, " :unable to read image build response")
		return "", err
	}

	defer imageBuildResponse.Body.Close()

	buildOutput, err := io.ReadAll(imageBuildResponse.Body)
	if err != nil {
		log.Print(err, ": unable to read build output")
		return "", err
	}
	return string(buildOutput), nil
}

func (o *Orchestrator) CreateWorker(imageName string) error {
	containerCreateResponse, err := o.client.ContainerCreate(o.ctx, &container.Config{
		Image: imageName,
		Tty:   false,
	}, nil, nil, nil, "")
	if err != nil {
		log.Print("Error creating container:", err)
		return err
	}

	err = o.client.ContainerStart(o.ctx, containerCreateResponse.ID, types.ContainerStartOptions{})

	if err != nil {
		log.Print("Error starting container:", err)
		return err
	}
	o.WorkerIDs = append(o.WorkerIDs, containerCreateResponse.ID)

	log.Printf("Docker container %s is running...\n", containerCreateResponse.ID)
	return nil
}

func (o *Orchestrator) RemoveWorker(workerID string) error {
	err := o.client.ContainerRemove(o.ctx, workerID, types.ContainerRemoveOptions{
		//RemoveLinks:   true, learn what links are
		RemoveVolumes: true,
		Force:         true})
	if err != nil {
		log.Printf("Error removing container:", err, ": unable to remove worker %s")
		return err
	}

	log.Printf("Worker %s successfully removed.\n", workerID)
	return nil
}

func (o *Orchestrator) ContainerLogs() error {
	statusCh, errCh := o.client.ContainerWait(o.ctx, o.WorkerIDs[0], container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			fmt.Println("Error waiting for container:", err)
			return err
		}
	case <-statusCh:
	}

	out, err := o.client.ContainerLogs(o.ctx, o.WorkerIDs[0], types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		fmt.Println("Error getting container logs:", err)
		return err
	}

	_, err = stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	if err != nil {
		fmt.Println("Error copying container logs to terminal:", err)
		return err
	}
	return nil
}
