package apt

import (
	"errors"
	"sahand.dev/chisme/internal/commandrunner"
	"sahand.dev/chisme/internal/persistence/models"
	"strings"
	"testing"
	"time"
)

func TestUpdatePackageSimulation(t *testing.T) {
	mockOutput := `Reading package lists...
Building dependency tree...
Reading state information...
The following packages will be upgraded:
   libc6
1 upgraded, 0 newly installed, 0 to remove and 0 not upgraded.
Inst libc6 [2.27-3ubuntu1.1] (2.27-3ubuntu1.2 Ubuntu:18.04/bionic [amd64])
Conf libc6 (2.27-3ubuntu1.2 Ubuntu:18.04/bionic [amd64])`

	mockRunner := &commandrunner.MockCommandRunner{Output: mockOutput}
	apt := &Apt{
		CommandRunner: mockRunner,
		CLI:           "apt",
	}

	pkg := &models.Package{Name: "libc6"}

	output, err := apt.UpdatePackageSimulation(pkg)
	if err != nil {
		t.Fatalf("DryRunUpdatePackage() failed: %v", err)
	}

	outputString := ""
	for line := range output {
		outputString += line + "\n"
	}

	if !strings.Contains(outputString, "libc6 [2.27-3ubuntu1.1] (2.27-3ubuntu1.2") {
		t.Errorf("DryRunUpdatePackage() output = %v, want it to contain package update info", output)
	}
}

func TestUpdatePackage_Success(t *testing.T) {
	mockRunner := &commandrunner.MockCommandRunner{
		Output: "line1\nline2\nline3\n",
		Err:    nil,
	}

	aptManager := &Apt{
		CommandRunner: mockRunner,
		CLI:           "apt",
	}

	output := make(chan string)
	pkg := &models.Package{Name: "example-package"}

	err := aptManager.UpdatePackage(pkg, output)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var lines []string
	for line := range output {
		lines = append(lines, line)
	}

	expectedLines := []string{"line1", "line2", "line3"}
	if len(lines) != len(expectedLines) {
		t.Fatalf("expected %d lines, got %d", len(expectedLines), len(lines))
	}
	for i, line := range lines {
		if line != expectedLines[i] {
			t.Errorf("expected line %d to be %q, got %q", i, expectedLines[i], line)
		}
	}
}

func TestUpdatePackage_Error(t *testing.T) {
	mockRunner := &commandrunner.MockCommandRunner{
		Output: "",
		Err:    []error{errors.New("some error")},
	}

	aptManager := &Apt{
		CommandRunner: mockRunner,
		CLI:           "apt",
	}

	output := make(chan string)
	pkg := &models.Package{Name: "example-package"}

	err := aptManager.UpdatePackage(pkg, output)
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}

	expectedError := "some error"
	if err.Error() != expectedError {
		t.Fatalf("expected error %q, got %q", expectedError, err.Error())
	}
}

func TestUpdatePackage_MixedOutput(t *testing.T) {
	mockRunner := &commandrunner.MockCommandRunner{
		Output: "line1\nline2\n",
		Err:    []error{errors.New("error1"), errors.New("error2")},
	}

	aptManager := &Apt{
		CommandRunner: mockRunner,
		CLI:           "apt",
	}

	output := make(chan string)
	pkg := &models.Package{Name: "example-package"}

	err := aptManager.UpdatePackage(pkg, output)
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}

	expectedError := "error1"
	if err.Error() != expectedError {
		t.Fatalf("expected error %q, got %q", expectedError, err.Error())
	}

	var lines []string
	for line := range output {
		lines = append(lines, line)
	}

	expectedLines := []string{"line1", "line2"}
	if len(lines) != len(expectedLines) {
		t.Fatalf("expected %d lines, got %d", len(expectedLines), len(lines))
	}
	for i, line := range lines {
		if line != expectedLines[i] {
			t.Errorf("expected line %d to be %q, got %q", i, expectedLines[i], line)
		}
	}
}

func TestUpdatePackage_EmptyOutput(t *testing.T) {
	mockRunner := &commandrunner.MockCommandRunner{
		Output: "",
		Err:    nil,
	}

	aptManager := &Apt{
		CommandRunner: mockRunner,
		CLI:           "apt",
	}

	output := make(chan string)
	pkg := &models.Package{Name: "example-package"}

	err := aptManager.UpdatePackage(pkg, output)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var lines []string
	for line := range output {
		lines = append(lines, line)
	}

	if len(lines) != 0 {
		t.Fatalf("expected no lines, got %d", len(lines))
	}
}

func TestUpdatePackage_LongRunning(t *testing.T) {
	mockRunner := &commandrunner.MockCommandRunner{
		Output: "line1\nline2\nline3\n",
		Err:    nil,
	}

	aptManager := &Apt{
		CommandRunner: mockRunner,
		CLI:           "apt",
	}

	output := make(chan string)
	pkg := &models.Package{Name: "example-package"}

	go func() {
		time.Sleep(100 * time.Millisecond)
		close(output)
	}()

	err := aptManager.UpdatePackage(pkg, output)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var lines []string
	for line := range output {
		lines = append(lines, line)
	}

	expectedLines := []string{"line1", "line2", "line3"}
	if len(lines) != len(expectedLines) {
		t.Fatalf("expected %d lines, got %d", len(expectedLines), len(lines))
	}
	for i, line := range lines {
		if line != expectedLines[i] {
			t.Errorf("expected line %d to be %q, got %q", i, expectedLines[i], line)
		}
	}
}
