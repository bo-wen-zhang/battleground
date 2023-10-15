package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func writeProgramToFile(program, filename string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("File does not exists or cannot be created")
		return err
	}
	if err := os.Truncate(filename, 0); err != nil {
		fmt.Printf("Failed to truncate: %v", err)
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	fmt.Fprintf(w, "%v\n", program)

	w.Flush()
	return nil
}

func main() {

	//program := `while True: input()`
	program := `x = input()
y = input()
print("Hello", x, y)`
	//program := `print("Hello")`
	//program := `while True: x = 0`
	//print("Hello")`
	filename := "test.py"
	err := writeProgramToFile(program, filename)
	if err != nil {
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 1000*time.Millisecond)
	var out bytes.Buffer
	var stderr bytes.Buffer
	r, w, err := os.Pipe()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer r.Close()
	echo := exec.CommandContext(ctx, "echo", "-e", `John\nHarris`)
	echo.Stdout = w
	err = echo.Start()
	if err != nil {
		fmt.Println("Error with echo:", err)
		return
	}
	defer echo.Wait()
	w.Close()
	code := exec.CommandContext(ctx, "python", filename)
	code.Stdin = r
	code.Stdout = &out
	code.Stderr = &stderr
	if err := code.Start(); err != nil {
		fmt.Println("Error with start:", err)
		//return
	}
	if err := code.Wait(); err != nil {
		fmt.Println("Error with wait:", err)
	}
	fmt.Println("Stdout:", out.String())
	fmt.Println("Stderr:", stderr.String())
	fmt.Printf("Usage: %+v", code.ProcessState.SysUsage().(*syscall.Rusage))
}
