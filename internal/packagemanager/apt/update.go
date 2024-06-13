package apt

import (
	"fmt"
	"log"
	"sahand.dev/chisme/internal/persistence/models"
)

// UpdatePackageSimulation simulates updating a package and returns a channel to read the output and stderr combined
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

// UpdatePackage updates a package and returns a channel to read the output and also listens
// on stderr and in case of err it will terminate the process and return the error
func (a *Apt) UpdatePackage(pkg *models.Package, output chan<- string) error {
	command := fmt.Sprintf("%s install --only-upgrade --simulate %s", a.CLI, pkg.Name)

	stdOutput, outputErrors, err := a.CommandRunner.RunCommandAsync(command)
	if err != nil {
		return fmt.Errorf("failed to execute command: %s, err: %w", command, err)
	}

	done := make(chan struct{})
	errorChan := make(chan error, 1)

	go func() {
		defer close(done)
		for line := range stdOutput {
			output <- line
		}
		close(output)
	}()

	go func() {
		for err := range outputErrors {
			if err != nil {
				errorChan <- err
				return
			}
		}
		errorChan <- nil
	}()

	select {
	case <-done:
		return <-errorChan
	case err := <-errorChan:
		return err
	}
}
