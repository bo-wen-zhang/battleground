// https://stackoverflow.com/questions/38804313/build-docker-image-from-go-code
package main

import (
	"archive/tar"
	"bytes"
	"context"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Fatal(err, " :unable to init client")
	}

	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()

	dockerFile := "Dockerfile"
	dockerFileReader, err := os.Open("../engine/Dockerfile")
	if err != nil {
		log.Fatal(err, " :unable to open Dockerfile")
	}
	readDockerFile, err := io.ReadAll(dockerFileReader)
	if err != nil {
		log.Fatal(err, " :unable to read dockerfile")
	}

	tarHeader := &tar.Header{
		Name: dockerFile,
		Size: int64(len(readDockerFile)),
	}
	err = tw.WriteHeader(tarHeader)
	if err != nil {
		log.Fatal(err, " :unable to write tar header")
	}
	_, err = tw.Write(readDockerFile)
	if err != nil {
		log.Fatal(err, " :unable to write tar body")
	}
	dockerFileTarReader := bytes.NewReader(buf.Bytes())

	imageBuildResponse, err := cli.ImageBuild(
		ctx,
		dockerFileTarReader,
		types.ImageBuildOptions{
			Context:    dockerFileTarReader,
			Dockerfile: dockerFile,
			Remove:     true})
	if err != nil {
		log.Fatal(err, " :unable to build docker image")
	}
	defer imageBuildResponse.Body.Close()
	_, err = io.Copy(os.Stdout, imageBuildResponse.Body)
	if err != nil {
		log.Fatal(err, " :unable to read image build response")
	}
}

// package main

// import (
// 	"context"
// 	"fmt"
// 	"io"
// 	"os"
// 	"archive/tar"

// 	"github.com/docker/docker/api/types"
// 	"github.com/docker/docker/client"
// )

// func main() {
// 	dockerfilePath := "../engine/Dockerfile"
// 	contextPath, err := tar.createTarArchive("../engine/")
// 	if err != nil {
// 		fmt.Println("Error opening context path:", err)
// 		return
// 	}
// 	defer contextPath.Close()

// 	cli, err := client.NewEnvClient()
// 	if err != nil {
// 		fmt.Println("Error creating Docker client:", err)
// 		return
// 	}

// 	imageBuildRes, err := cli.ImageBuild(context.Background(), contextPath, types.ImageBuildOptions{
// 		Dockerfile: dockerfilePath,
// 		//Context:    io.NopCloser(os.Stdin),
// 	})
// 	if err != nil {
// 		fmt.Println("Error building Docker image:", err)
// 		return
// 	}
// 	defer imageBuildRes.Body.Close()

// 	buildOutput, err := io.ReadAll(imageBuildRes.Body)
// 	if err != nil {
// 		fmt.Println("Error reading build output:", err)
// 		return
// 	}
// 	fmt.Println(string(buildOutput))

// 	containerRunRes, err := cli.ContainerCreate(context.Background(), nil, nil, nil, nil, "")
// 	if err != nil {
// 		fmt.Println("Error creating container:", err)
// 		return
// 	}

// 	err = cli.ContainerStart(context.Background(), containerRunRes.ID, types.ContainerStartOptions{})
// 	if err != nil {
// 		fmt.Println("Error starting container:", err)
// 		return
// 	}
// 	fmt.Printf("Docker container %s is running...\n", containerRunRes.ID)
// }
