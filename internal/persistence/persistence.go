package persistence

import (
	"sahand.dev/chisme/internal/persistence/models"
	"time"
)

// PackageStore is an interface that represents the persistence layer for packages
type PackageStore interface {
	Save(pkg *models.Package) (int, error)
	Update(pkg *models.Package) error
	UpdateLastUpdate(pkg *models.Package, t time.Time) error
	Get(id int) (*models.Package, error)
	GetAll() ([]*models.Package, error)
	GetByName(name string) (*models.Package, error)
	SaveOrUpdatePackage(pkg *models.Package) error
}
