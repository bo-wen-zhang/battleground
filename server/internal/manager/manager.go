package manager

import (
	pb "battleground-server/engine_service"
	"bufio"
	"context"
	"fmt"
	"sync"
	"time"

	"battleground-server/internal/util"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/timeout"
	q "github.com/oleiade/lane/v2"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Engine struct {
	containerID string
	port        string
	Stub        pb.EngineServiceClient
}

func (man *Manager) BuildEngine(port string) (string, error) {
	containerID, err := man.CreateEngineContainer(port)
	if err != nil {
		return "", err
	}

	conn, err := man.EstablishEngineConn(containerID, port)
	if err != nil {
		man.RemoveEngineContainer(containerID)
		return "", err
	}

	stub := pb.NewEngineServiceClient(conn)

	engine := &Engine{
		containerID: containerID,
		port:        port,
		Stub:        stub,
	}

	man.Engines[containerID] = engine
	return containerID, nil
}

func (man *Manager) EstablishEngineConn(containerID, port string) (*grpc.ClientConn, error) {

	timeoutValue := 2 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeoutValue)
	defer cancel()
	serverAddress := "127.0.0.1:" + port
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),                // No tls for now
		grpc.WithUnaryInterceptor(timeout.UnaryClientInterceptor(timeoutValue)), // Set request deadline
		grpc.WithBlock(), // Make dial context a blocking call
	}

	conn, err := grpc.DialContext(ctx, serverAddress, opts...)
	if err != nil {
		return nil, fmt.Errorf("EstablishEngineConn: Failed to dial engine id %s port %s: %w", containerID, port, err)
	}
	return conn, nil
}

// A manager manages the engines that the server connects to.
type Manager struct {
	ImageName           string //names of images
	dockerClient        *client.Client
	ctx                 context.Context
	logger              zerolog.Logger
	expectedEngineCount int                //indicates the number of active engines that should be managed
	Engines             map[string]*Engine //map of containerid to engines that are actively running
	availablePorts      *q.Deque[string]   //ports which can be used by engines for grpc
}

// Creates a new Manager
// Returns an error if docker client cannot be created
func NewManager(imageName string, cli *client.Client, logger zerolog.Logger, expectedEngineCount int, ports []string) (*Manager, error) {

	return &Manager{
		ImageName:           imageName,
		dockerClient:        cli,
		ctx:                 context.Background(),
		logger:              logger,
		Engines:             map[string]*Engine{},
		availablePorts:      q.NewDeque[string](ports...),
		expectedEngineCount: expectedEngineCount,
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

func (man *Manager) RemoveAllContainers() error {
	for _, engine := range man.Engines {
		err := man.dockerClient.ContainerRemove(man.ctx, engine.containerID, types.ContainerRemoveOptions{
			RemoveVolumes: true,
			Force:         true})
		if err != nil {
			man.logger.Error().Err(err).Msgf("Failed to remove container %s", engine.containerID)
			return err
		}
	}
	man.logger.Info().Msg("Removed all containers")
	return nil
}

// Currently not used.
func (man *Manager) EngineLogs(containerID string) error {
	//containerID := man.Engines[0].

	out, err := man.dockerClient.ContainerLogs(man.ctx, containerID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	})
	if err != nil {
		fmt.Println("Error getting container logs:", err)
		return err
	}
	defer out.Close()
	scanner := bufio.NewScanner(out)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	return nil
}

// Blocking until an engine is stopped or removed, removes engine from manager's list.
// This function also handles closing the manager's gRPC connection to the engine.
// Lastly, it should signal the manager to build a new engine.
func (man *Manager) WaitEngineShutdown(containerID string, buildNeeded chan<- string) (int64, error) {
	defer func() {
		port := man.Engines[containerID].port
		delete(man.Engines, containerID)
		buildNeeded <- port
	}()
	var status container.WaitResponse
	statusCh, errCh := man.dockerClient.ContainerWait(man.ctx, containerID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return 0, fmt.Errorf("HandleEngineShutdown: Error receiving engine status code: %w", err)
		}
	case status = <-statusCh:
	}

	return status.StatusCode, nil
}

// function to monitor the engines, and build a new one each time an engine shuts down
func (man *Manager) MaintainEngines() {
	buildNeeded := make(chan string, man.expectedEngineCount)
	built := make(chan string, man.expectedEngineCount)
	var wg sync.WaitGroup

	for i := 0; i < man.expectedEngineCount; i++ {
		port, ok := man.availablePorts.Pop()
		if !ok {
			panic(nil)
		}
		buildNeeded <- port
	}

	wg.Add(1)
	go func() {
		//Waits for signal that an engine is down
		for {
			id := <-built
			go man.WaitEngineShutdown(id, buildNeeded)
		}
	}()

	go func() {
		defer wg.Done()
		//Waits for signal to build a new engine
		for {
			port := <-buildNeeded
			//build engine
			id, err := man.BuildEngine(port)
			if err != nil {
				panic(err)
			}
			built <- id
			//for testing purposes
			res, err := man.Engines[id].Stub.GetProgramResult(context.Background(), &pb.Program{
				UserId:     69,
				SourceCode: "print(\"Helloo World\")",
			})
			if err != nil {
				man.logger.Panic().Err(err).Msg("Error getting job response")
			}
			fmt.Println(res)
			//
		}

	}()

	wg.Wait()
}
