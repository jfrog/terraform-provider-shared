package sdk

import (
	"reflect"
	"testing"
)

func TestCastToStringArr(t *testing.T) {
	testSlice := []interface{}{"a", "b", "c"}

	result := CastToStringArr(testSlice)

	expected := []string{"a", "b", "c"}
	if reflect.DeepEqual(expected, result) == false {
		t.Errorf("Incorrect slice. Expected %v: got: %v", expected, result)
	}
}

func TestCastToStringArr_exclude_nils(t *testing.T) {
	testSlice := []interface{}{"a", "b", nil}

	result := CastToStringArr(testSlice)

	if len(result) != 2 {
		t.Errorf("Incorrect length. Expected: 2, got: %d", len(result))
	}
}

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
