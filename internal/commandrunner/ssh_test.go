package commandrunner

import (
	"fmt"
	"strings"
	"testing"
)

// startSSHRunner starts a test SSH server on 127.0.0.1:9090 that understand echo command
// and returns a new SSHCommandRunner configured to connect to the server
func startSSHRunner(t *testing.T) (*SSHCommandRunner, func()) {
	stopServer := startTestSSHServer(t)

	config := SSHConfig{
		Host:       "127.0.0.1",
		Port:       9090,
		User:       "test",
		PrivateKey: generateClientPrivateKey(t),
	}

	runner, err := NewSSHCommandRunner(config)
	if err != nil {
		t.Fatalf("failed to create SSH command runner: %v", err)
	}
	return runner, stopServer
}

func TestSSHCommandRunner_RunCommand(t *testing.T) {
	runner, stopServer := startSSHRunner(t)
	defer stopServer()

	scanner, err := runner.RunCommand(ExecCommand{Command: "echo Hello, World!"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	for scanner.Scan() {
		if scanner.Text() != "Hello, World!" {
			t.Errorf("RunCommand() = %q expected %q", scanner.Text(), "Hello, World!")
		}
	}
}

func TestSSHCommandRunner_RunCommand_UnknownCommand(t *testing.T) {
	runner, stopServer := startSSHRunner(t)
	defer stopServer()

	_, err := runner.RunCommand(ExecCommand{Command: "uknown-command"})
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
}

func TestSSHCommandRunner_RunCommand_LargeOutput(t *testing.T) {
	runner, stopServer := startSSHRunner(t)
	defer stopServer()

	largeText := strings.Repeat("A", 10000)
	scanner, err := runner.RunCommand(ExecCommand{Command: fmt.Sprintf("echo %s", largeText)})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	scanner.Scan()
	output := scanner.Text()
	if output != largeText {
		t.Errorf("expected large output, got %q", output)
	}
}

func TestSSHCommandRunner_RunCommandAsync(t *testing.T) {
	runner, stopServer := startSSHRunner(t)
	defer stopServer()

	t.Run("EchoCommand", func(t *testing.T) {
		output, outputError, err := runner.RunCommandAsync(ExecCommand{Command: "echo Hello, World!"})
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
		output, outputError, err := runner.RunCommandAsync(ExecCommand{Command: "invalid_command"})
		if err != nil {
			t.Fatalf("RunCommandAsync() failed, error = %v", err)
		}

		got, err := readFromChannels(t, output, outputError)
		if err == nil {
			t.Fatalf("RunCommandAsync() expected to fail, error is nil")
		}

		if len(got) > 0 && !strings.Contains(got, "unknown command") {
			t.Fatalf("RunCommandAsync() error = %v, expected to contain %v", err, "unknown command")
		}
	})
}

func TestSSHCommandRunner_RunCommand_Concurrent(t *testing.T) {
	runner, stopServer := startSSHRunner(t)
	defer stopServer()

	commands := []string{"echo Hello", "echo World", "echo Concurrent"}
	expectedOutputs := []string{"Hello", "World", "Concurrent"}
	results := make(chan string, len(commands))
	errors := make(chan error, len(commands))

	for _, cmd := range commands {
		go func(command string) {
			scanner, err := runner.RunCommand(ExecCommand{Command: command})
			if err != nil {
				errors <- err
				return
			}

			if scanner.Scan() {
				results <- scanner.Text()
			} else {
				results <- ""
			}
		}(cmd)
	}

	for range commands {
		select {
		case err := <-errors:
			t.Fatalf("expected no error, got %v", err)
		case result := <-results:
			if !contains(expectedOutputs, result) {
				t.Errorf("unexpected output: %q", result)
			}
		}
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
