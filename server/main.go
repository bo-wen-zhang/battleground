// https://stackoverflow.com/questions/38804313/build-docker-image-from-go-code
package server

import (
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/stdcopy"
)

func main() {

	o, err := NewOrchestrator()
	if err != nil {
		return
	}

	dockerFileName := "Dockerfile"
	contextDirSrc := "../engine"
	buildOutput, err := o.BuildImage(dockerFileName, contextDirSrc, "battleground-engine")
	if err != nil {
		return
	}
	fmt.Print(buildOutput)

	err = o.CreateWorker("battleground-engine")

	statusCh, errCh := o.client.ContainerWait(o.ctx, o.WorkerIDs[0], container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			fmt.Println("Error waiting for container:", err)
			return
		}
	case <-statusCh:
	}

	out, err := o.client.ContainerLogs(o.ctx, o.WorkerIDs[0], types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		fmt.Println("Error getting container logs:", err)
		return
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	o.RemoveWorker(o.WorkerIDs[0])
}
