package apt

import (
	"fmt"
	"sahand.dev/chisme/internal/commandrunner"
	"sahand.dev/chisme/internal/persistence/models"
)

// Apt is a struct that represents the apt packagemanager manager
type Apt struct {
	CLI           string
	CommandRunner commandrunner.CommandRunner
}

// GetPackages lists all packages and returns them as a slice of Package structs
func (a *Apt) GetPackages() ([]*models.Package, error) {
	command := fmt.Sprintf("%s list", a.CLI)

	scanner, err := a.CommandRunner.RunCommand(commandrunner.ExecCommand{Command: command})
	if err != nil {
		return nil, fmt.Errorf("failed to execute command: %s, err: %w", command, err)
	}

	packages, err := parseOutputCommand(scanner, parseLineToPackage)
	if err != nil {
		return nil, fmt.Errorf("failed to parse output: %w", err)
	}

	return packages, nil

}

// GetUpgradablePackages lists all upgradeable packages and returns them as a slice of Package structs
func (a *Apt) GetUpgradablePackages() ([]*models.Package, error) {
	command := fmt.Sprintf("%s list --upgradable", a.CLI)

	scanner, err := a.CommandRunner.RunCommand(commandrunner.ExecCommand{Command: command})
	if err != nil {
		return nil, fmt.Errorf("failed to execute command: %s, err: %w", command, err)
	}

	packages, err := parseOutputCommand(scanner, parseLineToPackage)
	if err != nil {
		return nil, fmt.Errorf("failed to parse output: %w", err)
	}

	return packages, nil
}
