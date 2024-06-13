package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"sahand.dev/chisme/internal/commandrunner"
	"sahand.dev/chisme/internal/packagemanager"
	"sahand.dev/chisme/internal/packagemanager/apt"
	"sahand.dev/chisme/internal/persistence"
	"sahand.dev/chisme/internal/persistence/sqllitestore"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Load environment variables from .env file if it exists
	loadEnv()

	// Initialize the database
	db := initDatabase("./chisme.db")
	defer db.Close()

	// Initialize package store
	packageStore := sqllitestore.NewSQLitePackageStore(db)

	// Parse command-line arguments
	packageManagerArg, command, host, args := parseFlags()

	commandRunner := initCommandRunner(host)

	pkgManager := initPackageManager(packageManagerArg, commandRunner)

	executeCommand(command, pkgManager, packageStore, args)
}

func loadEnv() {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}
}

func initDatabase(path string) *sql.DB {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	if err := sqllitestore.SetupDatabase(db); err != nil {
		log.Fatalf("Failed to setup database: %v", err)
	}
	return db
}

func parseFlags() (string, string, string, []string) {
	packageManager := flag.String("package_manager", "apt", "The package manager to use (e.g., apt, yum)")
	command := flag.String("command", "list_upgradable", "The command to run (e.g., list_upgradable, update, remove, install)")
	host := flag.String("host", "local", "Host to execute command (e.g., localhost, ssh)")
	flag.Parse()
	return *packageManager, *command, *host, flag.Args()
}

func initCommandRunner(host string) commandrunner.CommandRunner {
	switch host {
	case "local":
		return &commandrunner.BashCommandRunner{}
	case "ssh":
		privateKey, err := os.ReadFile(getEnv("SSH_PRIVATE_KEY_PATH", ""))
		if err != nil {
			log.Fatalf("Error reading private key: %s", err)
		}
		sshConfig := commandrunner.SSHConfig{
			Host:               getEnv("SSH_HOST", "host"),
			Port:               getEnvAsInt("SSH_PORT", 22),
			User:               getEnv("SSH_USER", "u"),
			PrivateKey:         privateKey,
			PrivateKeyPassword: getEnv("SSH_PRIVATE_KEY_PASSWORD", "p"),
		}
		commandRunner, err := commandrunner.NewSSHCommandRunner(sshConfig)
		if err != nil {
			log.Fatalf("Error creating SSH command runner: %s", err)
		}
		return commandRunner
	default:
		log.Fatalf("Unsupported host: %s", host)
		return nil
	}
}

func initPackageManager(packageManager string, commandRunner commandrunner.CommandRunner) packagemanager.PackageManger {
	switch packageManager {
	case "apt":
		return &apt.Apt{
			CommandRunner: commandRunner,
			CLI:           "apt",
		}
	default:
		log.Fatalf("Unsupported package manager: %s", packageManager)
		return nil
	}
}

func executeCommand(command string, pkgManager packagemanager.PackageManger, packageStore persistence.PackageStore, args []string) {
	switch command {
	case "list_upgradable":
		listUpgradablePackages(pkgManager, packageStore)
	case "list_installed":
		listInstalledPackages(pkgManager, packageStore)
	case "update":
		updatePackage(pkgManager, packageStore, args)
	case "update_all":
		updateAllPackages(pkgManager, packageStore)
	default:
		log.Fatalf("Unsupported command: %s", command)
	}
}

func listUpgradablePackages(pkgManager packagemanager.PackageManger, packageStore persistence.PackageStore) {
	packages, err := pkgManager.GetUpgradablePackages()
	if err != nil {
		log.Fatalf("Error listing upgradable packages: %s", err)
	}
	for _, pkg := range packages {
		if err := packageStore.SaveOrUpdatePackage(pkg); err != nil {
			log.Printf("Error saving package: %s", err)
		}
		fmt.Printf("Package: %s, Installed Version: %s, New Version: %s\n", pkg.Name, pkg.InstalledVersion, pkg.Version)
	}
}

func listInstalledPackages(pkgManager packagemanager.PackageManger, packageStore persistence.PackageStore) {
	packages, err := pkgManager.GetPackages()
	if err != nil {
		log.Fatalf("Error listing installed packages: %s", err)
	}
	for _, pkg := range packages {
		if pkg.Installed {
			if err := packageStore.SaveOrUpdatePackage(pkg); err != nil {
				log.Printf("Error saving package: %s", err)
			}
			fmt.Printf("Package: %s, Installed Version: %s, New Version: %s\n", pkg.Name, pkg.InstalledVersion, pkg.Version)
		}
	}
}

func updatePackage(pkgManager packagemanager.PackageManger, packageStore persistence.PackageStore, args []string) {
	if len(args) == 0 {
		log.Fatalf("Package name required for update command")
	}
	packageName := args[0]
	pkg, err := packageStore.GetByName(packageName)
	if err != nil {
		log.Fatalf("Error getting package: %s", err)
	}
	outputStream := make(chan string)
	go func() {
		for line := range outputStream {
			fmt.Println(line)
		}
	}()
	if err := pkgManager.UpdatePackage(pkg, outputStream); err != nil {
		log.Fatalf("Error updating package: %s", err)
	}
	pkg.Installed = true
	pkg.InstalledVersion = pkg.Version
	if err := packageStore.Update(pkg); err != nil {
		log.Fatalf("Error updating package: %s", err)
	}
	if err := packageStore.UpdateLastUpdate(pkg, time.Now()); err != nil {
		log.Fatalf("Error updating last update: %s", err)
	}
}

func updateAllPackages(pkgManager packagemanager.PackageManger, packageStore persistence.PackageStore) {
	outputStream := make(chan string)
	go func() {
		for line := range outputStream {
			fmt.Println(line)
		}
	}()
	if err := pkgManager.UpdateAllPackages(outputStream); err != nil {
		log.Fatalf("Error updating package: %s", err)
	}

	packages, err := pkgManager.GetUpgradablePackages()
	if err != nil {
		log.Fatalf("Error listing upgradable packages: %s", err)
	}
	for _, aptPackage := range packages {
		pkg, _ := packageStore.GetByName(aptPackage.Name)
		if pkg == nil {
			pkg.LastUpdated = time.Now()
			if err := packageStore.SaveOrUpdatePackage(pkg); err != nil {
				log.Printf("Error saving package: %s", err)
			}
		} else if pkg.InstalledVersion != aptPackage.InstalledVersion {
			pkg.InstalledVersion = aptPackage.InstalledVersion
			if err := packageStore.Update(pkg); err != nil {
				log.Printf("Error updating package: %s", err)
			}
			if err := packageStore.UpdateLastUpdate(pkg, time.Now()); err != nil {
				log.Printf("Error updating package: %s", err)
			}
		}
		fmt.Printf("Package: %s, Installed Version: %s, New Version: %s\n", pkg.Name, pkg.InstalledVersion, pkg.Version)
	}
}

// Utility functions for environment variables
func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getEnvAsInt(key string, fallback int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return fallback
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return fallback
	}
	return value
}
