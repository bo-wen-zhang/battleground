package manager

import (
	"bufio"
	"context"
	"fmt"

	"battleground-server/internal/util"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/rs/zerolog"
)

type Manager struct {
	Images    []string //names of images
	WorkerIDs []string //ids of worker containers
	client    *client.Client
	ctx       context.Context
	logger    zerolog.Logger
}

func NewManager(logger zerolog.Logger) (*Manager, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logger.Fatal().Err(err).Msg("Unable to init client")
		return &Manager{}, err
	}
	return &Manager{
		client: cli,
		ctx:    context.Background(),
		logger: logger,
	}, nil
}

// Builds image from dockerfile and build context
func (man *Manager) BuildImage(dockerFileName, contextDirSrc, imageName string) (string, error) {
	dockerFileTarReader, err := util.CreateTarReader(contextDirSrc)
	if err != nil {
		man.logger.Fatal().Err(err).Msg("Unable to tarball build context")
		return "", err
	}

	imageBuildResponse, err := man.client.ImageBuild(
		man.ctx,
		dockerFileTarReader,
		types.ImageBuildOptions{
			Context:    dockerFileTarReader,
			Dockerfile: dockerFileName,
			Tags:       []string{imageName}})
	if err != nil {
		man.logger.Fatal().Err(err).Msg("Unable to build image")
		return "", err
	}
	defer imageBuildResponse.Body.Close()
	// _, err = io.Copy(os.Stdout, imageBuildResponse.Body)
	// if err != nil {
	// 	o.logger.Fatal().Err(err).Msg("Unable to read build response from stdout")
	// 	return "", err
	// }

	buildScanner := bufio.NewScanner(imageBuildResponse.Body)
	for buildScanner.Scan() {
		line := buildScanner.Text()
		// if strings.Contains(line, "error") {
		// 	fmt.Println(line)
		// }
		fmt.Println(line)
	}

	// buildOutput, err := io.ReadAll(imageBuildResponse.Body)
	// if err != nil {
	// 	man.logger.Fatal().Err(err).Msg("Unable to read build response into []byte")
	// 	return "", err
	// }
	man.logger.Info().Msg("Image built.")
	return "", nil
	//return string(buildOutput), nil
}

// Creates a worker container
func (man *Manager) CreateWorker(imageName string) error {

	entryPoint := []string{"./battleground-engine"}

	containerCreateResponse, err := man.client.ContainerCreate(man.ctx, &container.Config{
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
				Target: "go/logs",
			},
		},
		// Binds: []string{
		// 	"/home/bo/Documents/battleground/server/engine_logs.log:/engine_logs.log",
		// },
	}, nil, nil, "")
	if err != nil {
		man.logger.Fatal().Err(err).Msg("Unable to create container")
		return err
	}

	err = man.client.ContainerStart(man.ctx, containerCreateResponse.ID, types.ContainerStartOptions{})
	if err != nil {
		man.logger.Fatal().Err(err).Msg("Unable to start container")
		return err
	}
	man.WorkerIDs = append(man.WorkerIDs, containerCreateResponse.ID)

	man.logger.Info().Msgf("Started container %s", containerCreateResponse.ID)
	return nil
}

func (man *Manager) RemoveWorker(workerID string) error {
	err := man.client.ContainerRemove(man.ctx, workerID, types.ContainerRemoveOptions{
		//RemoveLinks:   true, learn what links are
		RemoveVolumes: true,
		Force:         true})
	if err != nil {
		man.logger.Error().Err(err).Msgf("Unable to remove container %s", workerID)
		return err
	}

	man.logger.Info().Msgf("Removed container %s", workerID)
	return nil
}

func (man *Manager) ContainerLogs() error {

	// waiter, err := man.client.ContainerAttach(man.ctx, man.WorkerIDs[0], types.ContainerAttachOptions{
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
		out, err := man.client.ContainerLogs(man.ctx, man.WorkerIDs[0], types.ContainerLogsOptions{
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

	statusCh, errCh := man.client.ContainerWait(man.ctx, man.WorkerIDs[0], container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			fmt.Println("Error waiting for container:", err)
			return err
		}
	case status := <-statusCh:
		man.logger.Info().Msgf("Container status code %#+v", status.StatusCode)
	}

	// _, err = stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	// if err != nil {
	// 	fmt.Println("Error copying container logs to terminal:", err)
	// 	return err
	// }
	return nil
}
