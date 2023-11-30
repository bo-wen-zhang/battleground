// https://stackoverflow.com/questions/38804313/build-docker-image-from-go-code
package main

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err, " :unable to init client")
	}

	dockerFile := "Dockerfile"
	contextDir := "../engine"

	dockerFileTarReader, err := createTarReader(contextDir)
	if err != nil {
		log.Fatal(err, " :unable to create tar reader")
	}

	imageBuildResponse, err := cli.ImageBuild(
		ctx,
		dockerFileTarReader,
		types.ImageBuildOptions{
			Context:    dockerFileTarReader,
			Dockerfile: dockerFile,
			Tags:       []string{"battleground-engine"}})
	if err != nil {
		log.Fatal(err, " :unable to build docker image")
	}
	defer imageBuildResponse.Body.Close()
	_, err = io.Copy(os.Stdout, imageBuildResponse.Body)
	if err != nil {
		log.Fatal(err, " :unable to read image build response")
	}

	defer imageBuildResponse.Body.Close()

	buildOutput, err := io.ReadAll(imageBuildResponse.Body)
	if err != nil {
		fmt.Println("Error reading build output:", err)
		return
	}
	fmt.Println(string(buildOutput))

	containerCreateResponse, err := cli.ContainerCreate(ctx, &container.Config{
		Image:        "battleground-engine:latest",
		AttachStdout: true,
	}, nil, nil, nil, "")
	if err != nil {
		fmt.Println("Error creating container:", err)
		return
	}

	err = cli.ContainerStart(ctx, containerCreateResponse.ID, types.ContainerStartOptions{})

	if err != nil {
		fmt.Println("Error starting container:", err)
		return
	}

	fmt.Printf("Docker container %s is running...\n", containerCreateResponse.ID)
}

func createTarReader(sourceDir string) (io.Reader, error) {
	r, w := io.Pipe()

	go func() {
		defer w.Close()

		gzipWriter := gzip.NewWriter(w)
		defer gzipWriter.Close()

		tarWriter := tar.NewWriter(gzipWriter)
		defer tarWriter.Close()

		err := filepath.WalkDir(sourceDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			info, err := d.Info()
			if err != nil {
				return err
			}
			header, err := tar.FileInfoHeader(info, d.Name())
			if err != nil {
				return err
			}

			relPath, err := filepath.Rel(sourceDir, path)
			if err != nil {
				return err
			}

			header.Name = relPath

			if err := tarWriter.WriteHeader(header); err != nil {
				return err
			}

			if !d.IsDir() {
				file, err := os.Open(path)
				if err != nil {
					return err
				}
				defer file.Close()

				_, err = io.Copy(tarWriter, file)
				if err != nil {
					return err
				}
			}

			return nil
		})

		if err != nil {
			w.CloseWithError(err)
		}

	}()

	return r, nil

	// go func() {
	// 	defer w.Close()

	// 	gzipWriter := gzip.NewWriter(w)
	// 	defer gzipWriter.Close()

	// 	tarWriter := tar.NewWriter(gzipWriter)
	// 	defer tarWriter.Close()

	// 	err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
	// 		if err != nil {
	// 			return err
	// 		}

	// 		header, err := tar.FileInfoHeader(info, info.Name())
	// 		if err != nil {
	// 			return err
	// 		}

	// 		relPath, err := filepath.Rel(sourceDir, path)
	// 		if err != nil {
	// 			return err
	// 		}
	// 		header.Name = relPath

	// 		if err := tarWriter.WriteHeader(header); err != nil {
	// 			return err
	// 		}

	// 		if !info.IsDir() {
	// 			file, err := os.Open(path)
	// 			if err != nil {
	// 				return err
	// 			}
	// 			defer file.Close()

	// 			_, err = io.Copy(tarWriter, file)
	// 			if err != nil {
	// 				return err
	// 			}
	// 		}

	// 		return nil
	// 	})

	// 	if err != nil {
	// 		w.CloseWithError(err)
	// 	}
	// }()

	// return r, nil
}
