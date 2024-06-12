package persistence

import (
	"sahand.dev/chisme/internal/persistence/models"
	"time"
)

type PackageStore interface {
	Save(pkg *models.Package) (int, error)
	Update(pkg *models.Package) error
	UpdateLastUpdate(pkg *models.Package, t time.Time) error
	Get(id int) (*models.Package, error)
	GetAll() ([]*models.Package, error)
	GetByName(name string) (*models.Package, error)
	SaveOrUpdatePackage(pkg *models.Package) error
}
