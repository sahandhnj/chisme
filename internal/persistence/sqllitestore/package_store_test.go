package sqllitestore

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"sahand.dev/chisme/internal/persistence/models"
	"testing"
	"time"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	if err := SetupDatabase(db); err != nil {
		t.Fatalf("failed to setup test database: %v", err)
	}

	return db
}

func TestSQLitePackageStore_SaveAndGet(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewSQLitePackageStore(db)

	pkg := &models.Package{
		Name:             "test-package",
		InstalledVersion: "1.0.0",
		Version:          "1.0.1",
		Installed:        true,
		LastUpdated:      time.Now(),
	}

	id, err := store.Save(pkg)
	if err != nil {
		t.Fatalf("failed to save package: %v", err)
	}

	retrievedPkg, err := store.Get(id)
	if err != nil {
		t.Fatalf("failed to get package: %v", err)
	}

	if retrievedPkg.Name != pkg.Name || retrievedPkg.InstalledVersion != pkg.InstalledVersion || retrievedPkg.Version != pkg.Version || retrievedPkg.Installed != pkg.Installed {
		t.Errorf("retrieved package does not match saved package: got %+v, want %+v", retrievedPkg, pkg)
	}
}

func TestSQLitePackageStore_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewSQLitePackageStore(db)

	pkg := &models.Package{
		Name:             "test-package",
		InstalledVersion: "1.0.0",
		Version:          "1.0.1",
		Installed:        true,
		LastUpdated:      time.Now(),
	}

	id, err := store.Save(pkg)
	if err != nil {
		t.Fatalf("failed to save package: %v", err)
	}

	pkg.ID = id
	pkg.Version = "1.0.2"
	if err := store.Update(pkg); err != nil {
		t.Fatalf("failed to update package: %v", err)
	}

	retrievedPkg, err := store.Get(id)
	if err != nil {
		t.Fatalf("failed to get package: %v", err)
	}

	if retrievedPkg.Version != "1.0.2" {
		t.Errorf("retrieved package version does not match updated version: got %s, want %s", retrievedPkg.Version, "1.0.2")
	}
}

func TestSQLitePackageStore_GetByName(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewSQLitePackageStore(db)

	pkg := &models.Package{
		Name:             "test-package",
		InstalledVersion: "1.0.0",
		Version:          "1.0.1",
		Installed:        true,
		LastUpdated:      time.Now(),
	}

	if _, err := store.Save(pkg); err != nil {
		t.Fatalf("failed to save package: %v", err)
	}

	retrievedPkg, err := store.GetByName(pkg.Name)
	if err != nil {
		t.Fatalf("failed to get package by name: %v", err)
	}

	if retrievedPkg.Name != pkg.Name || retrievedPkg.InstalledVersion != pkg.InstalledVersion || retrievedPkg.Version != pkg.Version || retrievedPkg.Installed != pkg.Installed {
		t.Errorf("retrieved package does not match saved package: got %+v, want %+v", retrievedPkg, pkg)
	}
}

func TestSQLitePackageStore_GetAll(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewSQLitePackageStore(db)

	// Save multiple packages
	packages := []*models.Package{
		{
			Name:             "test-package-1",
			InstalledVersion: "",
			Version:          "1.0.1",
			Installed:        true,
			LastUpdated:      time.Now(),
		},
		{
			Name:             "test-package-2",
			InstalledVersion: "",
			Version:          "2.0.1",
			Installed:        false,
			LastUpdated:      time.Now(),
		},
		{
			Name:             "test-package-3",
			InstalledVersion: "3.0.0",
			Version:          "3.0.1",
			Installed:        true,
			LastUpdated:      time.Now(),
		},
	}

	for _, pkg := range packages {
		_, err := store.Save(pkg)
		if err != nil {
			t.Fatalf("failed to save package: %v", err)
		}
	}

	// Retrieve all packages
	retrievedPackages, err := store.GetAll()
	if err != nil {
		t.Fatalf("failed to get all packages: %v", err)
	}

	if len(retrievedPackages) != len(packages) {
		t.Fatalf("expected to retrive %d packages, got %d", len(packages), len(retrievedPackages))
	}

	for i, pkg := range packages {
		retrievedPkg := retrievedPackages[i]
		if !retrievedPkg.DeepEqual(pkg) {
			t.Errorf("retrieved package does not match saved package: got %+v, want %+v", retrievedPkg, pkg)
		}
	}
}

func TestSQLitePackageStore_SaveOrUpdatePackage(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewSQLitePackageStore(db)

	pkg := &models.Package{
		Name:             "test-package",
		InstalledVersion: "1.0.0",
		Version:          "1.0.1",
		Installed:        true,
		LastUpdated:      time.Now(),
	}

	// Test save
	err := store.SaveOrUpdatePackage(pkg)
	if err != nil {
		t.Fatalf("failed to save package: %v", err)
	}

	retrievedPkg, err := store.GetByName(pkg.Name)
	if err != nil {
		t.Fatalf("failed to get package by name: %v", err)
	}

	if !retrievedPkg.Equals(pkg) {
		t.Errorf("retrieved package does not match saved package: got %+v, want %+v", retrievedPkg, pkg)
	}

	// Test update
	pkg.Version = "1.0.2"
	err = store.SaveOrUpdatePackage(pkg)
	if err != nil {
		t.Fatalf("failed to update package: %v", err)
	}

	retrievedPkg, err = store.GetByName(pkg.Name)
	if err != nil {
		t.Fatalf("failed to get package by name: %v", err)
	}

	if retrievedPkg.Version != "1.0.2" {
		t.Errorf("retrieved package version does not match updated version: got %s, want %s", retrievedPkg.Version, "1.0.2")
	}
}
