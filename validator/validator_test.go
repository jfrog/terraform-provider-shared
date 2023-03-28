package validator

import (
	"regexp"
	"testing"
)

func TestProjectKey(t *testing.T) {
	validProjectKeys := []string{
		"ab",                               // min 2
		"abcde123456789012345678901234567", // max 32
		"abc-123",                          // hyphen is supported but not documented
		"abc123-",                          // hyphen can be anywhere
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
		"a", // 1 character, too short
		"abcdefghijklmnopqrstu12345678901234567890", // 21 characters, too long
		"abc,",       // invalid characters
		"-abcde1234", // can't start with hyphen
		"ABC123",     // lower case only
		"123abc",     // must begin with alpha
	}

	for _, projectKey := range invalidProjectKeys {
		t.Run(projectKey, func(t *testing.T) {
			diag := ProjectKey(projectKey, nil)
			if !diag.HasError() {
				t.Errorf("ProjectKey '%s' should fail", projectKey)
			}

			errorRegex := regexp.MustCompile(`.*project_key must be 2 - 32 lowercase alphanumeric and hyphen characters.*`)
			if !errorRegex.MatchString(diag[0].Summary) {
				t.Fail()
			}
		})
	}
}
