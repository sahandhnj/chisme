package apt

import (
	"errors"
	"sahand.dev/chisme/internal/command_runner"
	"sahand.dev/chisme/internal/package_manager/models"
	"strings"
	"testing"
)

func TestApt_GetPackages(t *testing.T) {
	mockOutput := `
libc6/now 2.27-3ubuntu1.2 amd64 [installed]
libc6-dev/now 2.27-3ubuntu1.2 amd64 [installed]
0ad-data-common/noble,noble 0.0.26-1 all
`
	mockRunner := &command_runner.MockCommandRunner{Output: mockOutput}
	apt := &Apt{
		CLI:           "apt",
		CommandRunner: mockRunner,
	}

	packages, err := apt.GetPackages()
	if err != nil {
		t.Fatalf("GetPackages() failed: %v", err)
	}

	expected := []*models.Package{
		{Name: "libc6", CurrVersion: "2.27-3ubuntu1.2", Version: "2.27-3ubuntu1.2", Installed: true},
		{Name: "libc6-dev", CurrVersion: "2.27-3ubuntu1.2", Version: "2.27-3ubuntu1.2", Installed: true},
		{Name: "0ad-data-common", CurrVersion: "", Version: "0.0.26-1", Installed: false},
	}

	for i, _ := range packages {
		if !packages[i].Equals(expected[i]) {
			t.Errorf("Package at index %d is = %v, expected = %v", i, packages[i], expected[i])
		}
	}
}

func TestApt_GetUpgradablePackages(t *testing.T) {
	mockOutput := `
libc6/now 2.27-3ubuntu1.2 amd64 [upgradable from: 2.27-3ubuntu1.1]
libc6-dev/now 2.27-3ubuntu1.2 amd64 [upgradable from: 2.27-3ubuntu1.1]
`
	mockRunner := &command_runner.MockCommandRunner{Output: mockOutput}
	apt := &Apt{
		CLI:           "apt",
		CommandRunner: mockRunner,
	}

	packages, err := apt.GetUpgradablePackages()
	if err != nil {
		t.Fatalf("GetUpgradablePackages() failed: %v", err)
	}

	expected := []*models.Package{
		{Name: "libc6", CurrVersion: "2.27-3ubuntu1.1", Version: "2.27-3ubuntu1.2", Installed: true},
		{Name: "libc6-dev", CurrVersion: "2.27-3ubuntu1.1", Version: "2.27-3ubuntu1.2", Installed: true},
	}

	for i, _ := range packages {
		if !packages[i].Equals(expected[i]) {
			t.Errorf("Package at index %d is = %v, expected = %v", i, packages[i], expected[i])
		}
	}
}

func TestApt_CommandRunnerError(t *testing.T) {
	mockRunner := &command_runner.MockCommandRunner{Err: errors.New("command failed")}
	apt := &Apt{
		CLI:           "apt",
		CommandRunner: mockRunner,
	}

	_, err := apt.GetPackages()
	if err == nil || !strings.Contains(err.Error(), "command failed") {
		t.Fatalf("GetPackages() error = %v, want %v", err, "command failed")
	}
}

func TestApt_ParseOutputError(t *testing.T) {
	mockOutput := "invalid output"
	mockRunner := &command_runner.MockCommandRunner{Output: mockOutput}
	apt := &Apt{
		CLI:           "apt",
		CommandRunner: mockRunner,
	}

	_, err := apt.GetPackages()
	if err == nil || !strings.Contains(err.Error(), "unexpected number of fields in line") {
		t.Fatalf("GetPackages() error = %v, want %v", err, "unexpected number of fields in line")
	}
}

func TestGetUpgradablePackages_ReturnsListOfPackages(t *testing.T) {
	apt := &Apt{CommandRunner: &command_runner.BashCommandRunner{}, CLI: "apt"}
	_, err := apt.GetUpgradablePackages()
	if err != nil {
		t.Fatalf("GetUpgradablePackages() failed: %v", err)
	}
}

func TestPackageVersionMismatch(t *testing.T) {
	apt := &Apt{CommandRunner: &command_runner.BashCommandRunner{}, CLI: "apt"}
	packages, err := apt.GetUpgradablePackages()
	if err != nil {
		t.Fatalf("ListPackages() failed: %v", err)
	}
	if len(packages) > 0 {
		for _, pkg := range packages {
			if pkg.CurrVersion == pkg.Version {
				t.Fatalf("Package %s has the same current and new version: %s", pkg.Name, pkg.CurrVersion)
			}
		}
	}
}
