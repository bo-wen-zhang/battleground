package worker_test

import (
	"battleground-engine/worker"
	"os"
	"strings"
	"testing"
)

func TestHelloWorld(t *testing.T) {
	w := worker.NewWorker(1)
	input, err := w.WriteSolutionInput("")
	if err != nil {
		t.Fatal("Error writing empty input.")
	}
	res, err := w.ExecuteSolution(input, "../test_scripts/hello_world.py")
	if err != nil {
		t.Fatal("Error executing program.")
	}

	expected := "Hello World"
	actual := strings.TrimSpace(res.Stdout.String())

	if actual != expected {
		t.Errorf("Expected %v got %v instead", expected, actual)
	}
}

func TestCreateFile(t *testing.T) {
	w := worker.NewWorker(1)
	input, err := w.WriteSolutionInput("")
	if err != nil {
		t.Fatal("Error writing empty input.")
	}
	res, err := w.ExecuteSolution(input, "../test_scripts/create_file.py")
	if err != nil {
		t.Log(res.Stderr.String())
		return
	}
	t.Error("Error program executed successfully")

	if _, err := os.Stat("../test_scripts/user_created.txt"); err == nil {
		t.Error("Error file was created")
		err = os.Remove("user_created.txt")
		if err != nil {
			t.Log(err)
			t.Fatal()
		}
	}
}

func TestWriteInput(t *testing.T) {
	w := worker.NewWorker(1)
	input, err := w.WriteSolutionInput("Joe\nBloggs")
	if err != nil {
		t.Fatal("Error writing input.")
	}
	res, err := w.ExecuteSolution(input, "../test_scripts/print_input.py")
	if err != nil {
		t.Fatal("Error executing program:", res)
	}
	expected := "Hello Joe Bloggs"
	actual := strings.TrimSpace(res.Stdout.String())

	if actual != expected {
		t.Errorf("Expected %v got %v instead", expected, actual)
	}

}
