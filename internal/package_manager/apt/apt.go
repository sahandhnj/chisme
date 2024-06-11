package apt

import (
	"fmt"
	"sahand.dev/chisme/internal/command_runner"
	"sahand.dev/chisme/internal/package_manager/models"
)

// Apt is a struct that represents the apt package_manager manager
type Apt struct {
	CLI           string
	CommandRunner command_runner.CommandRunner
}

// GetPackages lists all packages and returns them as a slice of Package structs
func (a *Apt) GetPackages() ([]*models.Package, error) {
	command := fmt.Sprintf("%s list", a.CLI)

	scanner, err := a.CommandRunner.RunCommand(command)
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

	scanner, err := a.CommandRunner.RunCommand(command)
	if err != nil {
		return nil, fmt.Errorf("failed to execute command: %s, err: %w", command, err)
	}

	packages, err := parseOutputCommand(scanner, parseLineToPackage)
	if err != nil {
		return nil, fmt.Errorf("failed to parse output: %w", err)
	}

	return packages, nil
}

//
//// UpdatePackages updates the given packages
//func (a *Apt) UpdatePackages(packages []*models.Package) err {
//	for _, pkg := range packages {
//		command := fmt.Sprintf("%s-get install --only-upgrade -y %s", a.CLI, pkg.Name)
//		if _, err := a.CommandRunner.RunCommand(command); err != nil {
//			return fmt.Errorf("failed to update package_manager %s: %w", pkg.Name, err)
//		}
//	}
//	return nil
//}
//
//// RemovePackage removes the given package_manager
//func (a *Apt) RemovePackage(pkg *models.Package) err {
//	command := fmt.Sprintf("%s-get remove -y %s", a.CLI, pkg.Name)
//	if _, err := a.CommandRunner.RunCommand(command); err != nil {
//		return fmt.Errorf("failed to remove package_manager %s: %w", pkg.Name, err)
//	}
//	return nil
//}
//
//// InstallPackage installs the given package_manager
//func (a *Apt) InstallPackage(pkg *models.Package) err {
//	command := fmt.Sprintf("%s-get install -y %s", a.CLI, pkg.Name)
//	if _, err := a.CommandRunner.RunCommand(command); err != nil {
//		return fmt.Errorf("failed to install package_manager %s: %w", pkg.Name, err)
//	}
//	return nil
//}
