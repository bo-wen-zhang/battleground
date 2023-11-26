package main

import (
	"battleground-engine/worker"
	"bufio"
	"fmt"
	"os"
)

func writeProgramToFile(program, filename string) error {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("File does not exist or cannot be created")
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	fmt.Fprintf(w, "%v\n", program)

	w.Flush()
	return nil
}

func main() {

	w := worker.NewWorker(1)
	fmt.Println("Hello")
	input, err := w.WriteSolutionInput("")
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Println("Hello")
	res, err := w.ExecuteSolution(input, "test_scripts/helloworld.py")
	if err != nil {
		fmt.Println(err)
		fmt.Println("Stderr:", res.Stderr.String())
		return
	}
	fmt.Println("Hello")
	fmt.Println("Stdout:", res.Stdout.String())
}
