package server

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"battleground-server/internal/util"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/rs/zerolog"
)

type Config struct {
	//mode=testing to launch containers to run unit tests
	Mode string
}

type Orchestrator struct {
	Images    []string //names of images
	WorkerIDs []string //ids of worker containers
	client    *client.Client
	ctx       context.Context
	config    Config
	logger    zerolog.Logger
}

func NewOrchestrator(cfg Config, logger zerolog.Logger) (*Orchestrator, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logger.Fatal().Err(err).Msg("Unable to init client")
		return &Orchestrator{}, err
	}
	return &Orchestrator{
		client: cli,
		ctx:    context.Background(),
		config: cfg,
		logger: logger,
	}, nil
}

func (o *Orchestrator) BuildImage(dockerFileName, contextDirSrc, imageName string) (string, error) {
	dockerFileTarReader, err := util.CreateTarReader(contextDirSrc)
	if err != nil {
		o.logger.Fatal().Err(err).Msg("Unable to tarball build context")
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
		o.logger.Fatal().Err(err).Msg("Unable to build image")
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
	o.logger.Info().Msg("Image built.")
	return string(buildOutput), nil
}

func (o *Orchestrator) CreateWorker(imageName string) error {

	var entryPoint strslice.StrSlice

	if o.config.Mode == "testing" {
		entryPoint = []string{"go", "test", "-v", "./worker/worker_test.go"}
	} else if o.config.Mode == "development" {
		//entryPoint = []string{"go", "run", "main.go"}
		entryPoint = []string{"./battleground-engine"}
	} else {
		log.Fatal("Mode (testing|development) configured incorrectly.")
	}

	containerCreateResponse, err := o.client.ContainerCreate(o.ctx, &container.Config{
		Image:      imageName,
		Entrypoint: entryPoint,
		ExposedPorts: nat.PortSet{
			"8089/tcp": struct{}{},
		},
		AttachStderr: true,
		AttachStdin:  false,
		Tty:          true,
		AttachStdout: true,
		OpenStdin:    false,
	}, &container.HostConfig{
		PortBindings: nat.PortMap{
			"8089/tcp": []nat.PortBinding{
				{
					HostIP:   "127.0.0.1",
					HostPort: "8089",
				},
			}},
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: "/home/bo/Documents/battleground/server/logs",
				Target: "/go/logs",
			},
		},
		// Binds: []string{
		// 	"/home/bo/Documents/battleground/server/engine_logs.log:/engine_logs.log",
		// },
	}, nil, nil, "")
	if err != nil {
		o.logger.Fatal().Err(err).Msg("Unable to create container")
		return err
	}

	err = o.client.ContainerStart(o.ctx, containerCreateResponse.ID, types.ContainerStartOptions{})
	if err != nil {
		o.logger.Fatal().Err(err).Msg("Unable to start container")
		return err
	}
	o.WorkerIDs = append(o.WorkerIDs, containerCreateResponse.ID)

	o.logger.Info().Msgf("Started container %s", containerCreateResponse.ID)
	return nil
}

func (o *Orchestrator) RemoveWorker(workerID string) error {
	err := o.client.ContainerRemove(o.ctx, workerID, types.ContainerRemoveOptions{
		//RemoveLinks:   true, learn what links are
		RemoveVolumes: true,
		Force:         true})
	if err != nil {
		o.logger.Error().Err(err).Msgf("Unable to remove container %s", workerID)
		return err
	}

	o.logger.Info().Msgf("Removed container %s", workerID)
	return nil
}

func (o *Orchestrator) ContainerLogs() error {

	// waiter, err := o.client.ContainerAttach(o.ctx, o.WorkerIDs[0], types.ContainerAttachOptions{
	// 	Stderr: true,
	// 	Stdout: true,
	// 	Stdin:  true,
	// 	Stream: true,
	// })
	// if err != nil {
	// 	fmt.Println("Error attaching:", err)
	// 	return err
	// }

	// go io.Copy(os.Stdout, waiter.Reader)
	// go io.Copy(os.Stderr, waiter.Reader)
	// go io.Copy(waiter.Conn, os.Stdin)

	go func() {
		out, err := o.client.ContainerLogs(o.ctx, o.WorkerIDs[0], types.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Follow:     true,
		})
		if err != nil {
			fmt.Println("Error getting container logs:", err)
			return
		}
		defer out.Close()
		scanner := bufio.NewScanner(out)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	statusCh, errCh := o.client.ContainerWait(o.ctx, o.WorkerIDs[0], container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			fmt.Println("Error waiting for container:", err)
			return err
		}
	case status := <-statusCh:
		o.logger.Info().Msgf("Container status code %#+v", status.StatusCode)
	}

	// _, err = stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	// if err != nil {
	// 	fmt.Println("Error copying container logs to terminal:", err)
	// 	return err
	// }
	return nil
}
