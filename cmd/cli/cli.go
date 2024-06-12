package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"sahand.dev/chisme/internal/commandrunner"
	"sahand.dev/chisme/internal/packagemanager"
	"sahand.dev/chisme/internal/packagemanager/apt"
	"sahand.dev/chisme/internal/persistence"
	"sahand.dev/chisme/internal/persistence/sqllitestore"
	"time"
)

func main() {
	db, err := sql.Open("sqlite3", "./chisme.db")
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	// Initialize the database schema
	if err = sqllitestore.SetupDatabase(db); err != nil {
		log.Fatalf("failed to setup database: %v", err)
	}

	var packageStore persistence.PackageStore
	packageStore = sqllitestore.NewSQLitePackageStore(db)

	// Define command-line arguments
	packageManager := flag.String("packagemanager", "apt", "The packagemanager manager to use (e.g., apt, yum)")
	command := flag.String("command", "list_upgradable", "The command to run (e.g., list_upgradable, update, remove, install)")

	flag.Parse()
	args := flag.Args()

	// Initialize the appropriate packagemanager manager
	var pkgManager packagemanager.PackageManger
	switch *packageManager {
	case "apt":
		pkgManager = &apt.Apt{
			CommandRunner: &commandrunner.BashCommandRunner{},
			CLI:           "apt",
		}
	default:
		_, _ = fmt.Fprintf(os.Stderr, "Unsupported packagemanager manager: %s\n", *packageManager)
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
			err := packageStore.SaveOrUpdatePackage(pkg)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Error saving package: %s\n", err.Error())
			}
			fmt.Printf("Package: %s, Installed Version: %s, New Version: %s\n", pkg.Name, pkg.InstalledVersion, pkg.Version)
		}
	case "list_installed":
		packages, err := pkgManager.GetPackages()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error listing installed packages: %s\n", err.Error())
			os.Exit(1)
		}
		for _, pkg := range packages {
			if pkg.Installed {
				err := packageStore.SaveOrUpdatePackage(pkg)
				if err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "Error saving package: %s\n", err.Error())
				}
				fmt.Printf("Package: %s, Installed Version: %s, New Version: %s\n", pkg.Name, pkg.InstalledVersion, pkg.Version)
			}
		}
	case "update":
		packageName := args[0]
		pkg, err := packageStore.GetByName(packageName)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error getting package: %s\n", err.Error())
			os.Exit(1)
		}
		outputStream, err := pkgManager.UpdatePackageSimulation(pkg)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error simulating package installation: %s\n", err.Error())
			os.Exit(1)
		}
		for line := range outputStream {
			fmt.Println(line)
		}

		pkg.Installed = true
		pkg.InstalledVersion = pkg.Version
		err = packageStore.Update(pkg)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error updating package: %s\n", err.Error())
			os.Exit(1)
		}
		err = packageStore.UpdateLastUpdate(pkg, time.Now())
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error updating last update: %s\n", err.Error())
			os.Exit(1)

		}
	default:
		_, _ = fmt.Fprintf(os.Stderr, "Unsupported command: %s\n", *command)
		os.Exit(1)
	}
}
