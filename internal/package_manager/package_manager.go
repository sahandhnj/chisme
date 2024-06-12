package package_manager

import (
	"sahand.dev/chisme/internal/persistence/models"
)

// PackageManger is representing the abstraction of a package_manager manager
type PackageManger interface {
	GetPackages() ([]*models.Package, error)
	GetUpgradablePackages() ([]*models.Package, error)
	UpdatePackageSimulation(pkg *models.Package) (<-chan string, error)
}
