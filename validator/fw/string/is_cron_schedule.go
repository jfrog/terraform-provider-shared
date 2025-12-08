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

package string

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	// Embed IANA timezone database for consistent timezone validation
	// across all platforms (especially Windows which lacks native IANA support)
	_ "time/tzdata"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/robfig/cron/v3"
)

// Ensure our implementation satisfies the validator.String interface.
var _ validator.String = &cronScheduleValidator{}

type cronScheduleValidator struct{}

func (v cronScheduleValidator) Description(_ context.Context) string {
	return "value must be ONE of the following formats (cannot mix):\n" +
		"1. Standard cron expression: 'minute hour day-of-month month day-of-week'\n" +
		"   - minute must be 00, 15, 30, or 45\n" +
		"   - hour must be 00-23 (2 digits)\n" +
		"   - Example: '30 09 * * MON'\n" +
		"2. Predefined descriptor (starts with @):\n" +
		"   - @hourly   - Run once an hour at the beginning of the hour\n" +
		"   - @daily    - Run once a day at midnight\n" +
		"   - @midnight - Same as @daily\n" +
		"   - @weekly   - Run once a week at midnight on Sunday\n" +
		"   - @monthly  - Run once a month at midnight of first day\n" +
		"   - @yearly   - Run once a year at midnight of Jan 1\n" +
		"   - @annually - Same as @yearly\n" +
		"   - @every <duration> - Run at fixed intervals, e.g. @every 1h30m (must be positive duration)"
}

func (v cronScheduleValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v cronScheduleValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()

	// Check for mixing formats
	if strings.HasPrefix(value, "@") && strings.Contains(value, " * ") {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path,
			"Invalid cron format",
			"Cannot mix descriptor format (@hourly) with standard cron format (* * * * *). Use only one format.",
		))
		return
	}

	// First try to parse as a descriptor
	if strings.HasPrefix(value, "@") {
		// Special handling for @every to validate duration
		if strings.HasPrefix(value, "@every ") {
			durationStr := strings.TrimPrefix(value, "@every ")
			duration, err := time.ParseDuration(durationStr)
			if err != nil {
				response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
					request.Path,
					"Invalid duration in @every",
					fmt.Sprintf("Duration must be a valid Go duration string (e.g., 1h30m). Got: %s. Error: %s", durationStr, err),
				))
				return
			}
			if duration <= 0 {
				response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
					request.Path,
					"Invalid duration in @every",
					fmt.Sprintf("Duration must be positive. Got: %s", durationStr),
				))
				return
			}
		}

		_, err := cron.ParseStandard(value)
		if err != nil {
			response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
				request.Path,
				"Invalid cron descriptor",
				fmt.Sprintf("Must be one of: @hourly, @daily, @midnight, @weekly, @monthly, @yearly, @annually, or @every <duration>. Got: %s. Error: %s", value, err),
			))
		}
		return
	}

	// If not a descriptor, validate as standard cron expression
	parts := strings.Fields(value)

	if len(parts) != 5 {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path,
			"Invalid cron expression format",
			fmt.Sprintf("%d parts: %s, Standard cron expression must have exactly 5 parts (minute hour day-of-month month day-of-week)", len(parts), value),
		))
		return
	}

	// First validate minute and hour constraints
	validators := []struct {
		name     string
		validate func(string) (bool, string)
		allowed  string
	}{
		{"minute", validateMinute, "00, 15, 30, 45"},
		{"hour", validateHour, "00-23"},
	}

	for i, part := range parts[:2] { // Only validate minute and hour
		v := validators[i]
		if valid, errMsg := v.validate(part); !valid {
			response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
				request.Path,
				fmt.Sprintf("Invalid %s in cron expression", v.name),
				fmt.Sprintf("The %s field must be one of: %s. Got: %s. %s", v.name, v.allowed, part, errMsg),
			))
			return
		}
	}

	// Then validate using robfig/cron's parser
	_, err := cron.ParseStandard(value)
	if err != nil {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path,
			"Invalid cron expression",
			fmt.Sprintf("Failed to parse standard cron expression: %s", err),
		))
	}
}

func validateMinute(value string) (bool, string) {
	allowedMinutes := map[string]bool{
		"00": true,
		"15": true,
		"30": true,
		"45": true,
	}

	if !allowedMinutes[value] {
		return false, "Minutes must be exactly 00, 15, 30, or 45"
	}
	return true, ""
}

func validateHour(value string) (bool, string) {
	// Allow * for "every"
	if value == "*" {
		return true, ""
	}

	// Single value
	num, err := strconv.Atoi(value)
	if err != nil {
		return false, "Hour must be a number"
	}

	if num < 0 || num > 23 {
		return false, "Hour must be between 00-23"
	}

	// Check if the format is correct (2 digits)
	if len(value) == 1 {
		return false, "Use leading zero for single digit hours (e.g., 03 instead of 3)"
	}

	return true, ""
}

// IsCronSchedule returns a validator which ensures that any configured string
// value is ONE of the following (cannot mix formats):
//
// 1. A standard cron expression with:
//   - minutes limited to 00, 15, 30, or 45
//   - hours in 2-digit format (00-23)
//     Example: "30 09 * * MON"
//
// 2. A predefined descriptor:
//   - @hourly   - Run once an hour at the beginning of the hour
//   - @daily    - Run once a day at midnight
//   - @midnight - Same as @daily
//   - @weekly   - Run once a week at midnight on Sunday
//   - @monthly  - Run once a month at midnight of first day
//   - @yearly   - Run once a year at midnight of Jan 1
//   - @annually - Same as @yearly
//   - @every <duration> - Run at fixed intervals (e.g. @every 1h30m)
//     Duration must be positive (> 0)
func IsCronSchedule() validator.String {
	return &cronScheduleValidator{}
}

// Ensure our implementation satisfies the validator.String interface.
var _ validator.String = &cronScheduleTimezoneValidator{}

type cronScheduleTimezoneValidator struct{}

func (v cronScheduleTimezoneValidator) Description(_ context.Context) string {
	return "value must be a valid IANA timezone name (e.g., UTC, America/New_York, Europe/London)"
}

func (v cronScheduleTimezoneValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v cronScheduleTimezoneValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()

	if value == "" {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path,
			"Invalid timezone",
			"empty value. Must be a valid IANA timezone. For valid timezone formats, see: https://timeapi.io/documentation/iana-timezones",
		))
		return
	}

	// Check if the timezone actually exists in the IANA database
	_, err := time.LoadLocation(value)
	if err != nil {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path,
			"Invalid timezone",
			fmt.Sprintf("%s. Must be a valid IANA timezone. For valid timezone formats, see: https://timeapi.io/documentation/iana-timezones", value),
		))
	}
}

// IsCronScheduleTimezone returns a validator which ensures that any configured string
// value is a valid IANA timezone name (e.g., UTC, America/New_York, Europe/London).
func IsCronScheduleTimezone() validator.String {
	return &cronScheduleTimezoneValidator{}
}
