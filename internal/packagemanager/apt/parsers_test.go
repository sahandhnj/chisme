package apt

import (
	"bufio"
	"fmt"
	"sahand.dev/chisme/internal/persistence/models"
	"slices"
	"strings"
	"testing"
)

func TestOutputParser(t *testing.T) {
	tests := []struct {
		name    string
		scanner *bufio.Scanner
		parser  func(line string) (string, error)
		result  []string
		err     bool
	}{
		{
			name:    "scanner with 3 lines and a parser that returns the line as is",
			scanner: bufio.NewScanner(strings.NewReader("line1\nline2\nline3")),
			parser: func(line string) (string, error) {
				return line, nil
			},
			result: []string{"line1", "line2", "line3"},
			err:    false,
		},
		{
			name:    "scanner with 3 lines and a parser that returns the first character of the line",
			scanner: bufio.NewScanner(strings.NewReader("AA\nBB\nCC")),
			parser: func(line string) (string, error) {
				return line[:1], nil
			},
			result: []string{"A", "B", "C"},
			err:    false,
		},
		{
			name:    "scanner with empty string",
			scanner: bufio.NewScanner(strings.NewReader("")),
			parser: func(line string) (string, error) {
				return line, nil
			},
			result: []string{},
			err:    false,
		}, {
			name:    "scanner with lines but faulty parser",
			scanner: bufio.NewScanner(strings.NewReader("line1\nline2\nline3")),
			parser: func(line string) (string, error) {
				return "", fmt.Errorf("parser err")
			},
			err: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseOutputCommand(tt.scanner, tt.parser)
			if (err != nil) != tt.err {
				t.Errorf("expected to have err: %v, got: %s", tt.err, err.Error())
				return
			}
			if !slices.Equal(result, tt.result) {
				t.Errorf("expected %v, got %v", tt.result, result)
			}
		})
	}

}

func TestParseLineToUpgradablePackage(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected *models.Package
		err      bool
	}{
		{
			name: "valid line",
			line: "pkg-name/focal-updates 2:8.1.2269-1ubuntu5.23 amd64 [upgradable from: 2:8.1.2269-1ubuntu5.22]",
			expected: &models.Package{
				Name:             "pkg-name",
				InstalledVersion: "2:8.1.2269-1ubuntu5.22",
				Version:          "2:8.1.2269-1ubuntu5.23",
				Installed:        true,
			},
			err: false,
		},
		{
			name:     "line with missing fields, fails with error",
			line:     "pkgname description version",
			expected: nil,
			err:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseLineToPackage(tt.line)
			if (err != nil) != tt.err {
				t.Errorf("expected err: %v, got: %v", tt.err, err)
				return
			}
			if !tt.expected.Equals(result) {
				t.Errorf("expected: %v, got: %v", tt.expected, result)
			}
		})

	}
}

func TestParseLineToPackage(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected *models.Package
		err      error
	}{
		{
			name: "upgradable packagemanager",
			line: "libc6/now 2.27-3ubuntu1.2 amd64 [upgradable from: 2.27-3ubuntu1.1]",
			expected: &models.Package{
				Name:             "libc6",
				InstalledVersion: "2.27-3ubuntu1.1",
				Version:          "2.27-3ubuntu1.2",
				Installed:        true,
			},
			err: nil,
		},
		{
			name: "installed packagemanager",
			line: "yudit-common/noble,noble,now 3.1.0-1 all [installed,automatic]",
			expected: &models.Package{
				Name:             "yudit-common",
				InstalledVersion: "3.1.0-1",
				Version:          "3.1.0-1",
				Installed:        true,
			},
			err: nil,
		},
		{
			name: "installed packagemanager with same version",
			line: "yaru-theme-gtk/noble,noble,now 24.04.2-0ubuntu1 all [installed]",
			expected: &models.Package{
				Name:             "yaru-theme-gtk",
				InstalledVersion: "24.04.2-0ubuntu1",
				Version:          "24.04.2-0ubuntu1",
				Installed:        true,
			},
			err: nil,
		},
		{
			name:     "empty line",
			line:     "",
			expected: nil,
			err:      SkippingLineError,
		},
		{
			name:     "warning line",
			line:     "WARNING: some warning message",
			expected: nil,
			err:      SkippingLineError,
		},
		{
			name:     "invalid format",
			line:     "invalid formatLine",
			expected: nil,
			err:      fmt.Errorf("unexpected number of fields in line: invalid formatLine"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseLineToPackage(tt.line)
			if (err != nil && tt.err == nil) || (err == nil && tt.err != nil) || (err != nil && tt.err != nil && err.Error() != tt.err.Error()) {
				t.Errorf("parseLineToPackage() error = %v, want %v", err, tt.err)
				return
			}
			if !result.Equals(tt.expected) {
				t.Errorf("parseLineToPackage() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestShouldSkipLine(t *testing.T) {
	tests := []struct {
		line     string
		expected bool
	}{
		{"", true},
		{"WARNING", true},
		{"Listing...", true},
		{"some packagemanager", false},
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			actual := shouldSkipLine(tt.line)
			if actual != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, actual)
			}
		})
	}
}

func TestExtractPackageName(t *testing.T) {
	tests := []struct {
		line     string
		expected string
		err      bool
	}{
		{"pkgname/description version dist", "pkgname", false},
		{"pkgname-withdash/description version dist", "pkgname-withdash", false},
		{"pkgname/subpkg/description version dist", "pkgname", false},
		{"pkgname description version", "", true},
		{"12", "", true},
		{"", "", true},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("line=%s", tt.line), func(t *testing.T) {
			result, err := extractPackageName(tt.line)
			if (err != nil) != tt.err {
				t.Errorf("expected err: %v, got: %v", tt.err, err)
			}
			if result != tt.expected {
				t.Errorf("expected: %s, got: %s", tt.expected, result)
			}
		})
	}
}

func TestExtractVersion(t *testing.T) {
	tests := []struct {
		name     string
		fields   []string
		expected string
		err      error
	}{
		{
			name:     "valid version",
			fields:   []string{"packagemanager-name", "1.0.0"},
			expected: "1.0.0",
			err:      nil,
		},
		{
			name:     "missing version field",
			fields:   []string{"packagemanager-name"},
			expected: "",
			err:      fmt.Errorf("version field is missing"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := extractVersion(tt.fields)
			if (err != nil && tt.err == nil) || (err == nil && tt.err != nil) || (err != nil && tt.err != nil && err.Error() != tt.err.Error()) {
				t.Errorf("extractVersion() error = %v, want %v", err, tt.err)
				return
			}
			if result != tt.expected {
				t.Errorf("extractVersion() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestExtractCurrentVersion(t *testing.T) {
	tests := []struct {
		name     string
		fields   []string
		expected string
		err      bool
	}{
		{
			name:     "upgradable packagemanager",
			fields:   []string{"packagemanager-name", "1.0.0", "dist", "[upgradable", "from:", "0.9.0]"},
			expected: "0.9.0",
			err:      false,
		},
		{
			name:     "installed packagemanager",
			fields:   []string{"packagemanager-name", "1.0.0", "[installed]"},
			expected: "1.0.0",
			err:      false,
		},
		{
			name:     "missing installed version",
			fields:   []string{"packagemanager-name", "1.0.0"},
			expected: "1.0.0",
			err:      false,
		},
		{
			name:     "invalid upgradable field",
			fields:   []string{"packagemanager-name"},
			expected: "",
			err:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := extractInstalledVersion(tt.fields)
			if (err != nil) != tt.err {
				t.Errorf("extractInstalledVersion() error = %v, want %v", err, tt.err)
				return
			}
			if result != tt.expected {
				t.Errorf("extractInstalledVersion() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsInstalled(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected bool
	}{
		{
			name:     "installed packagemanager",
			line:     "packagemanager-name 1.0.0 [installed]",
			expected: true,
		},
		{
			name:     "not installed packagemanager",
			line:     "packagemanager-name 1.0.0",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isInstalled(tt.line)
			if result != tt.expected {
				t.Errorf("isInstalled() = %v, want %v", result, tt.expected)
			}
		})
	}
}
