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

package validator

import (
	"regexp"
	"testing"
)

func TestRepoKey(t *testing.T) {
	t.Parallel()

	validRepoKeys := []string{
		"a", // min 1
		"abcd123456789012345678912345678901234567890123456789012345678901", // max 64
		"abc-123",   // hyphen is supported
		"abc123-",   // hyphen can be anywhere
		"123abc123", // begin with number
	}

	for _, repoKey := range validRepoKeys {
		k := repoKey
		t.Run(k, func(t *testing.T) {
			t.Parallel()

			diag := RepoKey(k, nil)
			if diag.HasError() {
				t.Errorf("RepoKey validation failed. diag: %v", diag)
			}
		})
	}
}

func TestRepoKey_invalidKeys(t *testing.T) {
	t.Parallel()

	invalidRepoKeys := []string{
		"",  // empty
		" ", // has space
		"abcd1234567890123456789123456789012345678901234567890123456789012", // 65 chars, too long
		"abc,", // invalid characters
	}

	for _, repoKey := range invalidRepoKeys {
		k := repoKey
		t.Run(k, func(t *testing.T) {
			t.Parallel()

			diag := RepoKey(k, nil)
			if !diag.HasError() {
				t.Errorf("RepoKey '%s' should fail", k)
			}
		})
	}
}

func TestProjectKey(t *testing.T) {
	t.Parallel()

	validProjectKeys := []string{
		"ab",                               // min 2
		"abcde123456789012345678901234567", // max 32
		"abc-123",                          // hyphen is supported but not documented
		"abc123-",                          // hyphen can be anywhere
	}

	for _, projectKey := range validProjectKeys {
		projectKey := projectKey
		t.Run(projectKey, func(t *testing.T) {
			t.Parallel()

			diag := ProjectKey(projectKey, nil)
			if diag.HasError() {
				t.Errorf("ProjectKey validation failed. diag: %v", diag)
			}
		})
	}
}

func TestProjectKey_invalidKeys(t *testing.T) {
	t.Parallel()

	invalidProjectKeys := []string{
		"a",                                 // 1 character, too short
		"abcdefghijklmnopqrstu123456789012", // 33 characters, too long
		"abc,",                              // invalid characters
		"-abcde1234",                        // can't start with hyphen
		"ABC123",                            // lower case only
		"123abc",                            // must begin with alpha
	}

	for _, projectKey := range invalidProjectKeys {
		projectKey := projectKey
		t.Run(projectKey, func(t *testing.T) {
			t.Parallel()

			diag := ProjectKey(projectKey, nil)
			if !diag.HasError() {
				t.Errorf("ProjectKey '%s' should fail", projectKey)
			}

			errorRegex := regexp.MustCompile(`.*key must be 2 - 32 lowercase alphanumeric and hyphen characters.*`)
			if !errorRegex.MatchString(diag[0].Summary) {
				t.Fail()
			}
		})
	}
}
