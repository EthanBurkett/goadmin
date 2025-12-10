package plugins

import (
	"fmt"
	"strconv"
	"strings"
)

// SemVer represents a semantic version
type SemVer struct {
	Major int
	Minor int
	Patch int
}

// ParseSemVer parses a semantic version string (e.g., "1.2.3")
func ParseSemVer(version string) (*SemVer, error) {
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid semantic version format: %s (expected major.minor.patch)", version)
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid major version: %s", parts[0])
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid minor version: %s", parts[1])
	}

	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, fmt.Errorf("invalid patch version: %s", parts[2])
	}

	return &SemVer{
		Major: major,
		Minor: minor,
		Patch: patch,
	}, nil
}

// String returns the string representation of the version
func (v *SemVer) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// Compare compares two semantic versions
// Returns: -1 if v < other, 0 if v == other, 1 if v > other
func (v *SemVer) Compare(other *SemVer) int {
	if v.Major != other.Major {
		if v.Major < other.Major {
			return -1
		}
		return 1
	}

	if v.Minor != other.Minor {
		if v.Minor < other.Minor {
			return -1
		}
		return 1
	}

	if v.Patch != other.Patch {
		if v.Patch < other.Patch {
			return -1
		}
		return 1
	}

	return 0
}

// GreaterThan checks if v > other
func (v *SemVer) GreaterThan(other *SemVer) bool {
	return v.Compare(other) > 0
}

// LessThan checks if v < other
func (v *SemVer) LessThan(other *SemVer) bool {
	return v.Compare(other) < 0
}

// Equals checks if v == other
func (v *SemVer) Equals(other *SemVer) bool {
	return v.Compare(other) == 0
}

// IsCompatible checks if this version is compatible with a required version range
// minVersion and maxVersion are optional (can be empty string)
func (v *SemVer) IsCompatible(minVersion, maxVersion string) (bool, error) {
	if minVersion != "" {
		min, err := ParseSemVer(minVersion)
		if err != nil {
			return false, fmt.Errorf("invalid minimum version: %w", err)
		}
		if v.LessThan(min) {
			return false, nil
		}
	}

	if maxVersion != "" {
		max, err := ParseSemVer(maxVersion)
		if err != nil {
			return false, fmt.Errorf("invalid maximum version: %w", err)
		}
		if v.GreaterThan(max) {
			return false, nil
		}
	}

	return true, nil
}

// ValidateAPICompatibility validates plugin API compatibility
func ValidateAPICompatibility(apiVersion string, metadata PluginMetadata) error {
	if metadata.MinAPIVersion == "" && metadata.MaxAPIVersion == "" {
		return nil // No version constraints
	}

	currentVersion, err := ParseSemVer(apiVersion)
	if err != nil {
		return fmt.Errorf("invalid API version: %w", err)
	}

	compatible, err := currentVersion.IsCompatible(metadata.MinAPIVersion, metadata.MaxAPIVersion)
	if err != nil {
		return err
	}

	if !compatible {
		return fmt.Errorf("plugin %s requires API version %s-%s, but current version is %s",
			metadata.ID, metadata.MinAPIVersion, metadata.MaxAPIVersion, apiVersion)
	}

	return nil
}
