package models

import "testing"

func TestPackageEquals(t *testing.T) {
	tests := []struct {
		name     string
		pkg1     *Package
		pkg2     *Package
		expected bool
	}{
		{
			name:     "equal packages",
			pkg1:     &Package{Name: "libc6", CurrVersion: "2.27-3ubuntu1.1", Version: "2.27-3ubuntu1.2", Installed: true},
			pkg2:     &Package{Name: "libc6", CurrVersion: "2.27-3ubuntu1.1", Version: "2.27-3ubuntu1.2", Installed: true},
			expected: true,
		},
		{
			name:     "different names",
			pkg1:     &Package{Name: "libc6", CurrVersion: "2.27-3ubuntu1.1", Version: "2.27-3ubuntu1.2", Installed: true},
			pkg2:     &Package{Name: "libc6-dev", CurrVersion: "2.27-3ubuntu1.1", Version: "2.27-3ubuntu1.2", Installed: true},
			expected: false,
		},
		{
			name:     "different current versions",
			pkg1:     &Package{Name: "libc6", CurrVersion: "2.27-3ubuntu1.1", Version: "2.27-3ubuntu1.2", Installed: true},
			pkg2:     &Package{Name: "libc6", CurrVersion: "2.27-3ubuntu1.0", Version: "2.27-3ubuntu1.2", Installed: true},
			expected: false,
		},
		{
			name:     "different new versions",
			pkg1:     &Package{Name: "libc6", CurrVersion: "2.27-3ubuntu1.1", Version: "2.27-3ubuntu1.2", Installed: true},
			pkg2:     &Package{Name: "libc6", CurrVersion: "2.27-3ubuntu1.1", Version: "2.27-3ubuntu1.3", Installed: true},
			expected: false,
		},
		{
			name:     "different installed status",
			pkg1:     &Package{Name: "libc6", CurrVersion: "2.27-3ubuntu1.1", Version: "2.27-3ubuntu1.2", Installed: true},
			pkg2:     &Package{Name: "libc6", CurrVersion: "2.27-3ubuntu1.1", Version: "2.27-3ubuntu1.2", Installed: false},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pkg1.Equals(tt.pkg2)
			if result != tt.expected {
				t.Errorf("Equals() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestPackageString(t *testing.T) {
	tests := []struct {
		name     string
		pkg      *Package
		expected string
	}{
		{
			name: "all fields set",
			pkg: &Package{
				Name:        "libc6",
				CurrVersion: "2.27-3ubuntu1.1",
				Version:     "2.27-3ubuntu1.2",
				Installed:   true,
			},
			expected: `Package{Name: "libc6", CurrVersion: "2.27-3ubuntu1.1", NewVersion: "2.27-3ubuntu1.2", Installed: true}`,
		},
		{
			name: "not installed",
			pkg: &Package{
				Name:        "libc6",
				CurrVersion: "2.27-3ubuntu1.1",
				Version:     "2.27-3ubuntu1.2",
				Installed:   false,
			},
			expected: `Package{Name: "libc6", CurrVersion: "2.27-3ubuntu1.1", NewVersion: "2.27-3ubuntu1.2", Installed: false}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pkg.String()
			if result != tt.expected {
				t.Errorf("String() = %v, want %v", result, tt.expected)
			}
		})
	}
}
