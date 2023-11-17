package util

import "testing"

func TestCheckVersion_not_supported(t *testing.T) {
	versionToCheck := "1.0.0"
	supportedVersion := "2.0.0"

	isSupported, _ := CheckVersion(versionToCheck, supportedVersion)

	if isSupported {
		t.Errorf("Incorrect version support. Version to check: %s, supported version: %s", versionToCheck, supportedVersion)
	}
}

func TestCheckVersion_supported(t *testing.T) {
	versionToCheck := "1.1.0"
	supportedVersion := "1.0.0"

	isSupported, _ := CheckVersion(versionToCheck, supportedVersion)

	if !isSupported {
		t.Errorf("Incorrect version support. Version to check: %s, supported version: %s", versionToCheck, supportedVersion)
	}
}
