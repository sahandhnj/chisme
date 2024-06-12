package apt

import (
	"sahand.dev/chisme/internal/commandrunner"
	"sahand.dev/chisme/internal/persistence/models"
	"strings"
	"testing"
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
