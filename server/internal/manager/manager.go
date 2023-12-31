package manager

import (
	pb "battleground-server/engine_service"
	"bufio"
	"context"
	"fmt"
	"time"

	"battleground-server/internal/util"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/timeout"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Engine struct {
	containerID string
	port        string
	clientStub  pb.EngineServiceClient
	ctx         context.Context
}

func NewEngineConn(containerID, port string) (*grpc.ClientConn, error) {

	timeoutValue := 2 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeoutValue)
	defer cancel()
	serverAddress := "127.0.0.1:" + port
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(timeout.UnaryClientInterceptor(timeoutValue)),
	}
	conn, err := grpc.DialContext(ctx, serverAddress, opts...)
	if err != nil { //failed to dial
		return nil, err
	}

	return conn, nil
}

type Manager struct {
	ImageName    string   //names of images
	WorkerIDs    []string //ids of worker containers
	dockerClient *client.Client
	ctx          context.Context
	logger       zerolog.Logger
	engines      map[string]Engine //map of containerID to engine struct
	//studs        map[string]job_handler.JobClient
}

// Creates a new Manager
// Returns an error if docker client cannot be created
func NewManager(imageName string, logger zerolog.Logger) (*Manager, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logger.Fatal().Err(err).Msg("Unable to init client")
		return &Manager{}, err
	}
	return &Manager{
		ImageName:    imageName,
		dockerClient: cli,
		ctx:          context.Background(),
		logger:       logger,
		engines:      map[string]Engine{},
	}, nil
}

// Builds image from dockerfile and build context
// Fatals if it is unable to build an image as it means
// that future containers for the engine cannot be created
func (man *Manager) BuildImage(dockerFileName, contextDirSrc string) {
	dockerFileTarReader, err := util.CreateTarReader(contextDirSrc)
	if err != nil {
		man.logger.Fatal().Err(err).Msg("Unable to tarball build context")
	}

	imageBuildResponse, err := man.dockerClient.ImageBuild(
		man.ctx,
		dockerFileTarReader,
		types.ImageBuildOptions{
			Context:    dockerFileTarReader,
			Dockerfile: dockerFileName,
			Tags:       []string{man.ImageName}})
	if err != nil {
		man.logger.Fatal().Err(err).Msg("Unable to build image")
	}
	defer imageBuildResponse.Body.Close()

	buildScanner := bufio.NewScanner(imageBuildResponse.Body)
	for buildScanner.Scan() {
		line := buildScanner.Text()
		fmt.Println(line)
	}

	man.logger.Info().Msg("Image built.")
}

// Requires a port from host machine for mapping to the container
// Creates a container for the engine
// Returns the docker id of the container and an error if it fails to be created or started by docker
func (man *Manager) CreateEngineContainer(hostPort string) (string, error) {

	entryPoint := []string{"./battleground-engine"}

	containerCreateResponse, err := man.dockerClient.ContainerCreate(man.ctx, &container.Config{
		Image:      man.ImageName,
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
					HostPort: hostPort,
				},
			}},
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: "/home/bo/Documents/battleground/server/logs",
				Target: "/go/logs",
			},
		},
	}, nil, nil, "")
	if err != nil {
		man.logger.Error().Err(err).Msgf("Unable to create container with port %s", hostPort)
		return "", err
	}

	err = man.dockerClient.ContainerStart(man.ctx, containerCreateResponse.ID, types.ContainerStartOptions{})
	if err != nil {
		man.logger.Error().Err(err).Msgf("Unable to start container with port %s", hostPort)
		return "", err
	}

	man.logger.Info().Msgf("Started container %s at port %s ", containerCreateResponse.ID, hostPort)
	return containerCreateResponse.ID, nil
}

func (man *Manager) RemoveEngineContainer(engineID string) error {
	err := man.dockerClient.ContainerRemove(man.ctx, engineID, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true})
	if err != nil {
		man.logger.Error().Err(err).Msgf("Unable to remove container %s", engineID)
		return err
	}

	man.logger.Info().Msgf("Removed container %s", engineID)
	return nil
}

// TODO: Here I should remove engine container from docker if it crashes
func (man *Manager) ContainerLogs() error {

	go func() {
		out, err := man.dockerClient.ContainerLogs(man.ctx, man.WorkerIDs[0], types.ContainerLogsOptions{
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

	statusCh, errCh := man.dockerClient.ContainerWait(man.ctx, man.WorkerIDs[0], container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			fmt.Println("Error waiting for container:", err)
			return err
		}
	case status := <-statusCh:
		man.logger.Info().Msgf("Container status code %#+v", status.StatusCode)
	}

	return nil
}
