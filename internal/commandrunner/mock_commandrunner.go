package commandrunner

import (
	"bufio"
	"strings"
)

// MockCommandRunner is a mock implementation of the CommandRunner interface
type MockCommandRunner struct {
	Output      string
	Err         []error
	AskPassPath string
}

// RunCommand mocks the execution of a command and returns predefined output and error
func (m *MockCommandRunner) RunCommand(command ExecCommand) (*bufio.Scanner, error) {
	if m.Err != nil && len(m.Err) > 0 {
		return nil, m.Err[0]
	}
	return bufio.NewScanner(strings.NewReader(m.Output)), nil
}

// RunCommandAsync mocks the execution of a command asynchronously and returns predefined output and error
func (m *MockCommandRunner) RunCommandAsync(command ExecCommand) (<-chan string, <-chan error, error) {
	output := make(chan string)
	outputErrors := make(chan error)

	go func() {
		defer close(output)
		for _, line := range strings.Split(m.Output, "\n") {
			if len(line) > 0 {
				output <- line
			}
		}
	}()

	go func() {
		defer close(outputErrors)
		for _, err := range m.Err {
			outputErrors <- err
		}
	}()

	return output, outputErrors, nil
}
