package worker

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"
)

type codeResult struct {
	Stdout *bytes.Buffer
	Stderr *bytes.Buffer
	Usage  *syscall.Rusage
}

type Worker struct {
	timeout float32 `default:"1"`
}

func NewWorker(timeout float32) *Worker {
	return &Worker{
		timeout: timeout,
	}
}

//judge: input, program, output

func (w *Worker) WriteProgramToFile(program []byte, filename string) error {
	err := os.WriteFile(filename, program, 0644)
	if err != nil {
		log.Print("File does not exist or cannot be created")
		return err
	}

	// writer := bufio.NewWriter(file)
	// fmt.Fprintf(writer, "%v\n", program)

	// writer.Flush()
	return nil
}

func (w *Worker) WriteSolutionInput(s string) (*os.File, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
	defer cancel()

	read, write, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	echo := exec.CommandContext(ctx, "echo", "-e", s)
	echo.Stdout = write
	if err := echo.Start(); err != nil {
		fmt.Println("Error with echo:", err)
		read.Close()
		return nil, err
	}
	defer echo.Wait()
	write.Close()
	return read, nil
}

func (w *Worker) ExecuteSolution(solutionInput *os.File, solutionPath string) (*codeResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
	defer cancel()
	//originally I closed the os.File handle in the WriteSolutionInput method
	defer solutionInput.Close()

	res := &codeResult{
		Stdout: new(bytes.Buffer),
		Stderr: new(bytes.Buffer),
	}

	cmd := exec.CommandContext(ctx, "python3", solutionPath)
	//cmd.Dir = "./temp/"
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	uid := 1001
	gid := 1001
	cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(uid), Gid: uint32(gid), NoSetGroups: true}
	cmd.Stdin = solutionInput
	cmd.Stdout = res.Stdout
	cmd.Stderr = res.Stderr

	//spin up goroutine to kill children processes when context timeout
	go func() {
		<-ctx.Done()

		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}()
	if err := cmd.Start(); err != nil {
		log.Println("Error with start:", err)
		return res, err
	}
	if err := cmd.Wait(); err != nil {
		log.Println("Error with wait:", err)
		return res, err
	}

	return res, nil
}
