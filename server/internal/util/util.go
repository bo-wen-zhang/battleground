package util

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// Creates an io.Reader by tarballing a given directory (used for docker build context)
func CreateTarReader(sourceDir string) (io.Reader, error) {
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
}
