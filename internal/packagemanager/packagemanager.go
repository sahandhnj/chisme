package packagemanager

import (
	"sahand.dev/chisme/internal/persistence/models"
)

// PackageManger is representing the abstraction of a packagemanager manager
type PackageManger interface {
	GetPackages() ([]*models.Package, error)
	GetUpgradablePackages() ([]*models.Package, error)

	Refresh(output chan<- string) error

	UpdatePackageSimulation(pkg *models.Package) (<-chan string, error)
	UpdatePackage(pkg *models.Package, output chan<- string) error
	UpdateAllPackages(output chan<- string) error
}
