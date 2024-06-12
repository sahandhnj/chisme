package sqllite_store

import (
	"database/sql"
	"errors"
	"fmt"
	"sahand.dev/chisme/internal/persistence/models"
	"time"
)

// SQLitePackageStore is a struct that represents a SQLite implementation of the PackageStore interface
type SQLitePackageStore struct {
	db *sql.DB
}

// NewSQLitePackageStore is a function that returns a new SQLitePackageStore
func NewSQLitePackageStore(db *sql.DB) *SQLitePackageStore {
	return &SQLitePackageStore{db: db}
}

// Save is a method that saves a package to the SQLite database
func (s *SQLitePackageStore) Save(pkg *models.Package) (int, error) {
	result, err := s.db.Exec("INSERT INTO packages (name, installed_version, version, installed, last_updated) VALUES (?, ?, ?, ?, ?)", pkg.Name, pkg.InstalledVersion, pkg.Version, pkg.Installed, pkg.LastUpdated)
	if err != nil {
		return 0, fmt.Errorf("error saving package: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting last insert id: %w", err)
	}

	return int(id), nil
}

// Update is a method that updates a package in the SQLite database
func (s *SQLitePackageStore) Update(pkg *models.Package) error {
	_, err := s.db.Exec("UPDATE packages SET installed_version = ?, version = ?, installed = ? WHERE name = ?", pkg.InstalledVersion, pkg.Version, pkg.Installed, pkg.Name)
	if err != nil {
		return fmt.Errorf("error updating package: %w", err)
	}

	return nil
}

// UpdateLastUpdate is a method that updates the last update time of a package in the SQLite database
func (s *SQLitePackageStore) UpdateLastUpdate(pkg *models.Package, t time.Time) error {
	_, err := s.db.Exec("UPDATE packages SET last_updated = ? WHERE name = ?", t, pkg.Name)
	if err != nil {
		return fmt.Errorf("error updating last_updated of package: %w", err)
	}

	return nil
}

// Get is a method that retrieves a package from the SQLite database by its ID
func (s *SQLitePackageStore) Get(id int) (*models.Package, error) {
	row := s.db.QueryRow("SELECT name, installed_version, version, installed, last_updated FROM packages WHERE id = ?", id)

	var pkg models.Package
	err := row.Scan(&pkg.Name, &pkg.InstalledVersion, &pkg.Version, &pkg.Installed, &pkg.LastUpdated)
	if err != nil {
		return nil, fmt.Errorf("error getting package: %w", err)
	}

	return &pkg, nil
}

// GetByName is a method that retrieves a package from the SQLite database by its name
func (s *SQLitePackageStore) GetByName(name string) (*models.Package, error) {
	row := s.db.QueryRow("SELECT name, installed_version, version, installed, last_updated FROM packages WHERE name = ?", name)

	var pkg models.Package
	err := row.Scan(&pkg.Name, &pkg.InstalledVersion, &pkg.Version, &pkg.Installed, &pkg.LastUpdated)
	if err != nil {
		return nil, fmt.Errorf("error getting package by name: %w", err)
	}

	return &pkg, nil
}

// GetAll is a method that retrieves all packages from the SQLite database
func (s *SQLitePackageStore) GetAll() ([]*models.Package, error) {
	rows, err := s.db.Query("SELECT name, installed_version, version, installed, last_updated FROM packages")
	if err != nil {
		return nil, fmt.Errorf("error getting packages: %w", err)
	}
	defer rows.Close()

	var packages []*models.Package
	for rows.Next() {
		var pkg models.Package
		err = rows.Scan(&pkg.Name, &pkg.InstalledVersion, &pkg.Version, &pkg.Installed, &pkg.LastUpdated)
		if err != nil {
			return nil, fmt.Errorf("error scanning package: %w", err)
		}

		packages = append(packages, &pkg)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return packages, err
}

// SaveOrUpdatePackage inserts a new package or updates an existing package in the SQLite database
func (s *SQLitePackageStore) SaveOrUpdatePackage(pkg *models.Package) error {
	existingPkg, err := s.GetByName(pkg.Name)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("error checking package existence: %w", err)
	}

	if existingPkg == nil {
		_, err := s.Save(pkg)
		if err != nil {
			return fmt.Errorf("error saving package: %w", err)
		}
	} else {
		err := s.Update(pkg)
		if err != nil {
			return fmt.Errorf("error updating package: %w", err)
		}
	}

	return nil
}
