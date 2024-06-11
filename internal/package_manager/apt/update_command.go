package apt

import (
	"fmt"
	"log"
	"sahand.dev/chisme/internal/package_manager/models"
)

func (a *Apt) UpdatePackageSimulation(pkg *models.Package) (<-chan string, error) {
	command := fmt.Sprintf("%s install --only-upgrade --simulate %s", a.CLI, pkg.Name)

	output, outputErrors, err := a.CommandRunner.RunCommandAsync(command)
	if err != nil {
		return nil, fmt.Errorf("failed to execute command: %s, err: %w", command, err)
	}

	go func() {
		for err := range outputErrors {
			log.Println(err)
		}
	}()

	return output, nil
}
