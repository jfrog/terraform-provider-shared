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
