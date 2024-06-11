package command_runner

import (
	"bufio"
	"strings"
)

// MockCommandRunner is a mock implementation of the CommandRunner interface
type MockCommandRunner struct {
	Output string
	Err    error
}

// RunCommand mocks the execution of a command and returns predefined output and error
func (m *MockCommandRunner) RunCommand(command string) (*bufio.Scanner, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return bufio.NewScanner(strings.NewReader(m.Output)), nil
}

func (m *MockCommandRunner) RunCommandAsync(command string) (<-chan string, <-chan error, error) {
	if m.Err != nil {
		return nil, nil, m.Err
	}

	output := make(chan string)
	go func() {
		defer close(output)
		for _, line := range strings.Split(m.Output, "\n") {
			output <- line
		}
	}()

	return output, nil, nil
}
