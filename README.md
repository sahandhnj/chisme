# Chisme
Chisme is a Go package that provides functionality to run bash commands asynchronously and manage packages using different package managers like `apt`.

## Usage

### CLI Commands

#### 1. List Upgradable Packages
This command lists all upgradable packages using the specified package manager.

```sh
go run cmd/cli/cli.go --package_manager=apt --command=list_upgradable
```

#### 2. List Installed Packages
This command lists all installed packages using the specified package manager.
```sh
go run cmd/cli/cli.go --package_manager=apt --command=list_installed
```

#### 2. Install a Package (Simulation)
This command simulates the installation of a specified package using the specified package manager.
```sh
go run cmd/cli/cli.go --package_manager=apt --command=install PACKAGENAME
```