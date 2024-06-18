package commandrunner

import "bufio"

// ExecCommand is a struct that holds the command to be executed and
// a flag to indicate if the command should be run with elevated privileges
type ExecCommand struct {
	Command  string
	Elevated bool
}

// CommandRunner is an interface for running system commands
type CommandRunner interface {
	RunCommand(command ExecCommand) (*bufio.Scanner, error)
	RunCommandAsync(command ExecCommand) (<-chan string, <-chan error, error)
}
