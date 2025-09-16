package string_test

import (
	"context"
	"testing"

	validatorfw_string "github.com/jfrog/terraform-provider-shared/validator/fw/string"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestIsCronSchedule(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         types.String
		expectError bool
	}
	tests := map[string]testCase{
		// Null/Unknown handling
		"unknown": {
			val: types.StringUnknown(),
		},
		"null": {
			val: types.StringNull(),
		},

		// Standard cron expressions - Valid Minutes
		"valid_minute_00": {
			val: types.StringValue("00 12 * * MON"),
		},
		"valid_minute_15": {
			val: types.StringValue("15 12 * * MON"),
		},
		"valid_minute_30": {
			val: types.StringValue("30 12 * * MON"),
		},
		"valid_minute_45": {
			val: types.StringValue("45 12 * * MON"),
		},

		// Standard cron expressions - Valid Hours
		"valid_hour_00": {
			val: types.StringValue("00 00 * * MON"),
		},
		"valid_hour_09": {
			val: types.StringValue("15 09 * * MON"),
		},
		"valid_hour_15": {
			val: types.StringValue("30 15 * * MON"),
		},
		"valid_hour_23": {
			val: types.StringValue("45 23 * * MON"),
		},

		// Standard cron expressions - Valid Special Characters
		"valid_with_star_hour": {
			val: types.StringValue("15 * * * MON"),
		},
		"valid_with_all_stars": {
			val: types.StringValue("30 * * * *"),
		},
		"valid_with_question_mark_dom": {
			val: types.StringValue("45 12 ? * MON"),
		},

		// Standard cron expressions - Valid Names
		"valid_with_month_name": {
			val: types.StringValue("00 08 * JAN MON"),
		},
		"valid_with_day_name": {
			val: types.StringValue("15 09 * * SUN"),
		},
		"valid_with_both_names": {
			val: types.StringValue("30 10 * DEC SAT"),
		},

		// Descriptors - Basic
		"valid_hourly": {
			val: types.StringValue("@hourly"),
		},
		"valid_daily": {
			val: types.StringValue("@daily"),
		},
		"valid_midnight": {
			val: types.StringValue("@midnight"),
		},
		"valid_weekly": {
			val: types.StringValue("@weekly"),
		},
		"valid_monthly": {
			val: types.StringValue("@monthly"),
		},
		"valid_yearly": {
			val: types.StringValue("@yearly"),
		},
		"valid_annually": {
			val: types.StringValue("@annually"),
		},

		// Descriptors - @every variations
		"valid_every_minute": {
			val: types.StringValue("@every 1m"),
		},
		"valid_every_hour": {
			val: types.StringValue("@every 1h"),
		},
		"valid_every_day": {
			val: types.StringValue("@every 24h"),
		},
		"valid_every_complex": {
			val: types.StringValue("@every 1h30m"),
		},

		// Invalid - Minutes
		"invalid_minute_single_digit": {
			val:         types.StringValue("0 12 * * MON"),
			expectError: true,
		},
		"invalid_minute_not_quarter": {
			val:         types.StringValue("10 12 * * MON"),
			expectError: true,
		},
		"invalid_minute_range": {
			val:         types.StringValue("0-30 12 * * MON"),
			expectError: true,
		},
		"invalid_minute_list": {
			val:         types.StringValue("0,15,30 12 * * MON"),
			expectError: true,
		},
		"invalid_minute_step": {
			val:         types.StringValue("*/15 12 * * MON"),
			expectError: true,
		},

		// Invalid - Hours
		"invalid_hour_single_digit": {
			val:         types.StringValue("00 3 * * MON"),
			expectError: true,
		},
		"invalid_hour_too_large": {
			val:         types.StringValue("00 24 * * MON"),
			expectError: true,
		},
		"invalid_hour_negative": {
			val:         types.StringValue("00 -1 * * MON"),
			expectError: true,
		},
		"invalid_hour_range": {
			val:         types.StringValue("15 9-17 * * MON"),
			expectError: true,
		},
		"invalid_hour_list": {
			val:         types.StringValue("30 8,9,10 * * MON"),
			expectError: true,
		},
		"invalid_hour_step": {
			val:         types.StringValue("45 */2 * * MON"),
			expectError: true,
		},

		// Invalid - Field Count
		"invalid_too_few_parts": {
			val:         types.StringValue("30 12 * *"),
			expectError: true,
		},
		"invalid_too_many_parts": {
			val:         types.StringValue("30 12 * * 5 0"),
			expectError: true,
		},

		// Invalid - Descriptors
		"invalid_descriptor_name": {
			val:         types.StringValue("@invalid"),
			expectError: true,
		},
		"invalid_every_no_duration": {
			val:         types.StringValue("@every"),
			expectError: true,
		},
		"invalid_every_bad_duration": {
			val:         types.StringValue("@every abc"),
			expectError: true,
		},
		"invalid_every_zero": {
			val:         types.StringValue("@every 0h"),
			expectError: true,
		},
		"invalid_every_negative": {
			val:         types.StringValue("@every -1h"),
			expectError: true,
		},

		// Invalid - Mixed Formats
		"invalid_mixed_descriptor_and_stars": {
			val:         types.StringValue("@daily * * * *"),
			expectError: true,
		},
		"invalid_mixed_at_with_cron": {
			val:         types.StringValue("@30 12 * * MON"),
			expectError: true,
		},
		"invalid_mixed_descriptor_with_values": {
			val:         types.StringValue("@hourly 12 * * MON"),
			expectError: true,
		},

		// Invalid - Special Cases
		"invalid_empty_string": {
			val:         types.StringValue(""),
			expectError: true,
		},
		"invalid_only_spaces": {
			val:         types.StringValue("     "),
			expectError: true,
		},
		"invalid_wrong_day_name": {
			val:         types.StringValue("00 12 * * MONDAY"),
			expectError: true,
		},
		"invalid_wrong_month_name": {
			val:         types.StringValue("00 12 * JANUARY MON"),
			expectError: true,
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			request := validator.StringRequest{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    test.val,
			}
			response := validator.StringResponse{}
			validatorfw_string.IsCronSchedule().ValidateString(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}

func TestIsCronScheduleTimezone(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         types.String
		expectError bool
	}
	tests := map[string]testCase{
		// Null/Unknown handling
		"unknown": {
			val: types.StringUnknown(),
		},
		"null": {
			val: types.StringNull(),
		},

		// Valid timezones
		"valid_utc": {
			val: types.StringValue("UTC"),
		},
		"valid_gmt": {
			val: types.StringValue("GMT"),
		},
		"valid_est": {
			val: types.StringValue("EST"),
		},
		"valid_america_new_york": {
			val: types.StringValue("America/New_York"),
		},
		"valid_europe_london": {
			val: types.StringValue("Europe/London"),
		},
		"valid_asia_tokyo": {
			val: types.StringValue("Asia/Tokyo"),
		},
		"valid_australia_sydney": {
			val: types.StringValue("Australia/Sydney"),
		},
		"valid_pacific_auckland": {
			val: types.StringValue("Pacific/Auckland"),
		},

		// Invalid timezones
		"invalid_empty": {
			val:         types.StringValue(""),
			expectError: true,
		},
		"invalid_spaces": {
			val:         types.StringValue("   "),
			expectError: true,
		},
		"invalid_timezone": {
			val:         types.StringValue("Invalid/Timezone"),
			expectError: true,
		},
		"invalid_format": {
			val:         types.StringValue("UTC+01:00"),
			expectError: true,
		},
		"invalid_numeric": {
			val:         types.StringValue("+0100"),
			expectError: true,
		},
		"invalid_case": {
			val:         types.StringValue("utc"),
			expectError: false, // time.LoadLocation accepts case-insensitive timezone names
		},
		"invalid_partial": {
			val:         types.StringValue("America"),
			expectError: true,
		},
		"invalid_special_chars": {
			val:         types.StringValue("UTC!"),
			expectError: true,
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			request := validator.StringRequest{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    test.val,
			}
			response := validator.StringResponse{}
			validatorfw_string.IsCronScheduleTimezone().ValidateString(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}
