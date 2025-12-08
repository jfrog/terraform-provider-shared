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
