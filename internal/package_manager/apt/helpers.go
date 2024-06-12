package apt

import (
	"bufio"
	"errors"
	"fmt"
	"sahand.dev/chisme/internal/persistence/models"
	"strings"
)

var SkippingLineError = errors.New("skipping line")

// parseOutputCommand accepts a scanner and a function to parse each line of the output of a command
func parseOutputCommand[T any](scanner *bufio.Scanner, parseFunc func(string) (T, error)) ([]T, error) {
	var output []T

	for scanner.Scan() {
		line := scanner.Text()

		parsed, err := parseFunc(line)
		if err != nil {
			if errors.Is(err, SkippingLineError) {
				continue
			}
			return nil, fmt.Errorf("failed to parse line: %w", err)
		}

		output = append(output, parsed)
	}

	return output, nil
}

// parseLineToUpgradablePackage parses a line of output from `apt list --upgradable` into a Package struct
func parseLineToPackage(line string) (*models.Package, error) {
	const expectedMinFields = 3

	if shouldSkipLine(line) {
		return nil, SkippingLineError
	}

	// Split the line into fields based on whitespace
	fields := strings.Fields(line)
	if len(fields) < expectedMinFields {
		return nil, fmt.Errorf("unexpected number of fields in line: %s", line)
	}

	packageName, err := extractPackageName(line)
	if err != nil {
		return nil, err
	}

	version, err := extractVersion(fields)
	if err != nil {
		return nil, err
	}

	installedVersion := ""
	installed := isInstalled(line)
	if installed {
		installedVersion, err = extractInstalledVersion(fields)
		if err != nil {
			return nil, err
		}
	}

	return &models.Package{
		Name:             packageName,
		Version:          version,
		InstalledVersion: installedVersion,
		Installed:        installed,
	}, nil
}

// shouldSkipLine checks if the line should be skipped based on it's size and prefix
func shouldSkipLine(line string) bool {
	if len(line) == 0 {
		return true
	}

	for _, prefix := range []string{"WARNING", "Listing..."} {
		if strings.HasPrefix(line, prefix) {
			return true
		}

	}

	return false
}

// extractPackageName extracts the package_manager name from the first field of the line, which is in the format pkgname/description
func extractPackageName(line string) (string, error) {
	parts := strings.Split(line, "/")

	if len(parts) < 2 {
		return "", fmt.Errorf("failed to extract package_manager name from line: %s", line)
	}

	return parts[0], nil
}

// extractVersion extracts the new version from the fields of the line
func extractVersion(fields []string) (string, error) {
	if len(fields) >= 2 {
		return fields[1], nil
	}
	return "", fmt.Errorf("version field is missing")
}

// extractInstalledVersion extracts the installed version from the fields of the line
func extractInstalledVersion(fields []string) (string, error) {
	if len(fields) == 6 {
		parts := strings.Split(fields[5], "]")
		if len(parts) < 1 {
			return "", fmt.Errorf("failed to extract installed version from field: %s", fields[5])
		}
		return parts[0], nil
	}

	// For installed packages, the installed version is the new version
	if len(fields) >= 2 {
		return strings.TrimSpace(strings.TrimLeft(fields[1], " ")), nil
	}

	return "", errors.New("cannot find the version")
}

// isInstalled checks if the package_manager is installed
func isInstalled(line string) bool {
	return strings.Contains(line, "installed") || strings.Contains(line, "upgradable")
}
