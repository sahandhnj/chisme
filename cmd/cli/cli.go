package main

import (
	"flag"
	"fmt"
	"os"
	"sahand.dev/chisme/internal/command_runner"
	"sahand.dev/chisme/internal/package_manager"
	"sahand.dev/chisme/internal/package_manager/apt"
	"sahand.dev/chisme/internal/package_manager/models"
)

func main() {
	// Define command-line arguments
	packageManager := flag.String("package_manager", "apt", "The package_manager manager to use (e.g., apt, yum)")
	command := flag.String("command", "list_upgradable", "The command to run (e.g., list_upgradable, update, remove, install)")

	flag.Parse()
	args := flag.Args()

	// Initialize the appropriate package_manager manager
	var pkgManager package_manager.PackageManger
	switch *packageManager {
	case "apt":
		pkgManager = &apt.Apt{
			CommandRunner: &command_runner.BashCommandRunner{},
			CLI:           "apt",
		}
	default:
		_, _ = fmt.Fprintf(os.Stderr, "Unsupported package_manager manager: %s\n", *packageManager)
		os.Exit(1)
	}

	switch *command {
	case "list_upgradable":
		packages, err := pkgManager.GetUpgradablePackages()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error listing upgradable packages: %s\n", err.Error())
			os.Exit(1)
		}
		for _, pkg := range packages {
			fmt.Printf("Package: %s, Current Version: %s, New Version: %s\n", pkg.Name, pkg.CurrVersion, pkg.Version)
		}
	case "list_installed":
		packages, err := pkgManager.GetPackages()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error listing installed packages: %s\n", err.Error())
			os.Exit(1)
		}
		for _, pkg := range packages {
			if pkg.Installed {
				fmt.Printf("Package: %s, Current Version: %s, New Version: %s\n", pkg.Name, pkg.CurrVersion, pkg.Version)
			}
		}
	case "install":
		packageName := args[0]
		outputStream, err := pkgManager.UpdatePackageSimulation(&models.Package{Name: packageName})
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error simulating package installation: %s\n", err.Error())
			os.Exit(1)
		}
		for line := range outputStream {
			fmt.Println(line)
		}

	default:
		_, _ = fmt.Fprintf(os.Stderr, "Unsupported command: %s\n", *command)
		os.Exit(1)
	}
}
