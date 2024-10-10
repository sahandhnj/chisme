package commandrunner

import (
	"strings"
	"testing"
)

func TestBashCommandRunner_RunCommand_Cat(t *testing.T) {
	cmdRunner := BashCommandRunner{}

	scanner, err := cmdRunner.RunCommand(ExecCommand{Command: "echo 'Hello, World!'"})
	if err != nil {
		t.Fatalf("RunCommand() failed, error = %v", err)
	}

	for scanner.Scan() {
		if scanner.Text() != "Hello, World!" {
			t.Errorf("RunCommand() = %q expected %q", scanner.Text(), "Hello, World!")
		}
	}
}

func TestBashCommandRunner_RunCommand_NoExistingCommand(t *testing.T) {
	cmdRunner := BashCommandRunner{}

	_, err := cmdRunner.RunCommand(ExecCommand{Command: "non-existing-command"})
	if err == nil {
		t.Fatalf("RunCommand() expected an error but got nil")
	}

	if !strings.Contains(err.Error(), "execution failed") {
		t.Fatalf("RunCommand() error = %v, expected to contain %v", err, "execution failed")
	}
}

func TestBashCommandRunner_RunCommandAsync(t *testing.T) {
	cmdRunner := BashCommandRunner{}

	t.Run("EchoCommand", func(t *testing.T) {
		output, outputError, err := cmdRunner.RunCommandAsync(ExecCommand{Command: "echo 'Hello, World!'"})
		if err != nil {
			t.Fatalf("RunCommandAsync() failed, error = %v", err)
		}

		got, err := readFromChannels(t, output, outputError)
		if err != nil {
			t.Fatalf("RunCommandAsync() failed, error = %v", err)
		}

		expected := "Hello, World!"
		if got != expected {
			t.Errorf("RunCommandAsync() = %q, want %q", got, expected)
		}
	})

	t.Run("CatCommand", func(t *testing.T) {
		output, outputError, err := cmdRunner.RunCommandAsync(ExecCommand{Command: "echo 'Hello, World!' | cat"})
		if err != nil {
			t.Fatalf("RunCommandAsync() failed, error = %v", err)
		}

		got, err := readFromChannels(t, output, outputError)
		if err != nil {
			t.Fatalf("RunCommandAsync() failed, error = %v", err)
		}

		expected := "Hello, World!"
		if got != expected {
			t.Errorf("RunCommandAsync() = %q, want %q", got, expected)
		}
	})

	t.Run("CommandWithError", func(t *testing.T) {
		output, outputError, err := cmdRunner.RunCommandAsync(ExecCommand{Command: "invalid_command"})
		if err != nil {
			t.Fatalf("RunCommandAsync() failed, error = %v", err)
		}

		got, err := readFromChannels(t, output, outputError)
		if err == nil {
			t.Fatalf("RunCommandAsync() expected to fail, error is nil")
		}

		if len(got) > 0 && !strings.Contains(got, "command not found") {
			t.Fatalf("RunCommandAsync() error = %v, expected to contain %v", err, "command not found")
		}
	})
}
