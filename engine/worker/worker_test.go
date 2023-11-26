package worker_test

import (
	"battleground-engine/worker"
	"strings"
	"testing"
)

func TestHelloWorld(t *testing.T) {
	w := worker.NewWorker(1)
	input, err := w.WriteSolutionInput("")
	if err != nil {
		t.Fatal("Error writing empty input.")
	}
	res, err := w.ExecuteSolution(input, "../test_scripts/helloworld.py")
	if err != nil {
		t.Fatal("Error executing program.")
	}

	expected := "Hello World"
	actual := strings.TrimSpace(res.Stdout.String())

	if actual != expected {
		t.Errorf("Expected %v got %v instead", expected, actual)
	}

}
