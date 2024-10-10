package commandrunner

import (
	"fmt"
	"testing"
)

func TestApplyCommandRootElevation_ShouldAddSudoWithStdinInput(t *testing.T) {
	runner := MockCommandRunner{}

	command := "apt install"
	expected := "sudo -S apt install"

	applyCommandRootElevation(&command, runner.AskPassPath)

	if command != expected {
		t.Errorf("expected %q, got %q", expected, command)
	}
}

func TestApplyCommandRootElevation_ShouldAddSudoWithDashAskPass(t *testing.T) {
	runner := MockCommandRunner{
		AskPassPath: "/path/to/askpass",
	}

	command := "apt install"
	expected := fmt.Sprintf("SUDO_ASKPASS=%s sudo -A %s", runner.AskPassPath, command)

	applyCommandRootElevation(&command, runner.AskPassPath)

	if command != expected {
		t.Errorf("expected %q, got %q", expected, command)
	}
}
