package worker

import (
	"bytes"
	"context"
	"fmt"
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

type worker struct {
	timeout float32 `default:"1"`
}

func NewWorker(timeout float32) *worker {
	return &worker{
		timeout: timeout,
	}
}

//judge: input, program, output

func (w *worker) WriteSolutionInput(s string) (*os.File, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
	defer cancel()

	read, write, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	//defer read.Close()
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

func (w *worker) ExecuteSolution(solutionInput *os.File, solutionPath string) (*codeResult, error) {
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
		//fmt.Println("Error with start:", err)
		return res, err
	}
	if err := cmd.Wait(); err != nil {
		//fmt.Println("Error with wait:", err)
		return res, err
	}

	return res, nil
	//fmt.Println("Stdout:", out.String())
	//fmt.Println("Stderr:", stderr.String())
	//fmt.Printf("Usage: %+v", cmd.ProcessState.SysUsage().(*syscall.Rusage))
}
