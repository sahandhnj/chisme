package models

import "fmt"

// Package represents a package_manager that can be upgraded
type Package struct {
	Name        string
	CurrVersion string
	Version     string
	Installed   bool
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
		p.CurrVersion == other.CurrVersion &&
		p.Version == other.Version &&
		p.Installed == other.Installed
}

func (p *Package) String() string {
	return fmt.Sprintf("Package{Name: %q, CurrVersion: %q, NewVersion: %q, Installed: %t}",
		p.Name, p.CurrVersion, p.Version, p.Installed)
}
