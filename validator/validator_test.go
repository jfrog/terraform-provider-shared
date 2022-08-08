package validator

import (
	"regexp"
	"testing"
)

func TestProjectKey(t *testing.T) {
	validProjectKeys := []string{
		"abc", // min 3
		"abcde12345", // max 10
		"abc-123", // hyphen is supported but not documented
	}

	for _, projectKey := range validProjectKeys {
		t.Run(projectKey, func(t *testing.T) {
			diag := ProjectKey(projectKey, nil)
			if diag.HasError() {
				t.Errorf("ProjectKey validation failed. diag: %v", diag)
			}
		})
	}
}

func TestProjectKey_invalidKeys(t *testing.T) {
	invalidProjectKeys := []string{
		"ab", // 2 characters, too short
		"abcdefghijk", // 11 characters, too long
		"abc,", // invalid characters
		"-abcde1234", // can't start with hyphen
		"ABC123", // lower case only
		"123abc", // must begin with alpha
	}

	for _, projectKey := range invalidProjectKeys {
		t.Run(projectKey, func(t *testing.T) {
			diag := ProjectKey(projectKey, nil)
			if !diag.HasError() {
				t.Errorf("ProjectKey '%s' should fail", projectKey)
			}

			errorRegex := regexp.MustCompile(`.*project_key must be 3 - 10 lowercase alphanumeric and hyphen characters.*`)
			if !errorRegex.MatchString(diag[0].Summary) {
				t.Fail()
			}
		})
	}
}
