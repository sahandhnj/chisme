package models

import (
	"fmt"
	"time"
)

// Package represents a package_manager that can be upgraded
type Package struct {
	ID               int       `json:"id"`
	Name             string    `json:"name"`
	InstalledVersion string    `json:"installed_version"`
	Version          string    `json:"version"`
	Installed        bool      `json:"installed"`
	LastUpdated      time.Time `json:"last_updated"`
}

// Equals compares two Package instances for equality
func (p *Package) Equals(other *Package) bool {
	if p == other {
		return true
	}
	if other == nil {
		return false
	}
	return p.Name == other.Name &&
		p.InstalledVersion == other.InstalledVersion &&
		p.Version == other.Version &&
		p.Installed == other.Installed
}

// DeepEqual compares two Package instances for deep equality ( including LastUpdated and ID)
func (p *Package) DeepEqual(other *Package) bool {
	return p.Equals(other) && p.LastUpdated.Equal(other.LastUpdated) && p.ID == other.ID
}

func (p *Package) String() string {
	return fmt.Sprintf("Package{Name: %q, InstalledVersion: %q, NewVersion: %q, Installed: %t}",
		p.Name, p.InstalledVersion, p.Version, p.Installed)
}
