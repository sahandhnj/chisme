package apt

import (
	"fmt"
	"log"
	"sahand.dev/chisme/internal/commandrunner"
	"sahand.dev/chisme/internal/persistence/models"
)

// UpdatePackageSimulation simulates updating a package and returns a channel to read the output and stderr combined
func (a *Apt) UpdatePackageSimulation(pkg *models.Package) (<-chan string, error) {
	command := fmt.Sprintf("%s install --only-upgrade --simulate %s", a.CLI, pkg.Name)

	output, outputErrors, err := a.CommandRunner.RunCommandAsync(commandrunner.ExecCommand{Command: command})
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
// on stderr and in case of error it will terminate the process and return the error
func (a *Apt) UpdatePackage(pkg *models.Package, output chan<- string) error {
	command := fmt.Sprintf("%s install --only-upgrade --simulate %s", a.CLI, pkg.Name)

	return a.exec(command, output)
}

// UpdateAllPackages updates all packages and returns a channel to read the output and also listens
// on stderr and in case of error it will terminate the process and return the error
func (a *Apt) UpdateAllPackages(output chan<- string) error {
	command := fmt.Sprintf("%s upgrade -y", a.CLI)

	return a.exec(command, output)
}

func (a *Apt) Refresh(output chan<- string) error {
	command := fmt.Sprintf("%s update", a.CLI)

	return a.exec(command, output)
}

// exec runs the specified command and manages stdout and stderr
func (a *Apt) exec(command string, output chan<- string) error {
	stdOutput, outputErrors, err := a.CommandRunner.RunCommandAsync(commandrunner.ExecCommand{Command: command})
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
