package command_runner

import "bufio"

// CommandRunner is an interface for running system commands
type CommandRunner interface {
	RunCommand(command string) (*bufio.Scanner, error)
	RunCommandAsync(command string) (<-chan string, <-chan error, error)
}
