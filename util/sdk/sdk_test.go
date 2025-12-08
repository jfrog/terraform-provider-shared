// Copyright (c) JFrog Ltd. (2025)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
