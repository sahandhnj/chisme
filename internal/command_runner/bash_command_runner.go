package command_runner

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
)

// BashCommandRunner implements CommandRunner for bash commands
type BashCommandRunner struct{}

// RunCommand runs a bash command and returns the output as a scanner
func (b *BashCommandRunner) RunCommand(command string) (*bufio.Scanner, error) {
	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("command execution failed: %w", err)
	}

	scanner := bufio.NewScanner(bytes.NewReader(output))
	return scanner, nil
}

// RunCommandAsync runs a bash command asynchronously and returns a channel with the output lines
func (b *BashCommandRunner) RunCommandAsync(command string) (<-chan string, <-chan error, error) {
	output := make(chan string, 10)
	errorsChan := make(chan error)

	cmd := exec.Command("bash", "-c", command)

	// Setup command
	stdOut, stdErr, err := setupCommand(cmd)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to setup command: %w", err)
	}

	// Start command
	err = cmd.Start()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start command: %w", err)
	}

	go func() {
		defer close(output)
		defer close(errorsChan)
		pipeCommandOutputLineByLine(stdOut, stdErr, output, errorsChan)

		if err := cmd.Wait(); err != nil {
			errorsChan <- fmt.Errorf("command finished with error: %w", err)
		}
	}()

	return output, errorsChan, nil
}

// setupCommand sets up the command and returns the stdout and stderr pipes
func setupCommand(cmd *exec.Cmd) (io.ReadCloser, io.ReadCloser, error) {
	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get stdout pipe: %w", err)
	}
	stdErr, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get stderr pipe: %w", err)
	}
	return stdOut, stdErr, err
}

// pipeCommandOutputLineByLine pipes the command output line by line to the output channel
func pipeCommandOutputLineByLine(stdOut io.ReadCloser, stdErr io.ReadCloser, output chan string, errorsChan chan error) {
	scanner := bufio.NewScanner(io.MultiReader(stdOut, stdErr))
	for scanner.Scan() {
		output <- scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		errorsChan <- fmt.Errorf("error reading stderr: %w", err)
	}
}
