package string_test

import (
	"context"
	"testing"

	validatorfw_string "github.com/jfrog/terraform-provider-shared/validator/fw/string"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestIsURLHttpOrHttps(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         types.String
		expectError bool
	}
	tests := map[string]testCase{
		"unknown": {
			val: types.StringUnknown(),
		},
		"null": {
			val: types.StringNull(),
		},
		"valid http": {
			val: types.StringValue("http://tempurl.org"),
		},
		"valid https": {
			val: types.StringValue("https://tempurl.org"),
		},
		"no host": {
			val:         types.StringValue("http://"),
			expectError: true,
		},
		"invalid scheme": {
			val:         types.StringValue("test://tempurl.org"),
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
			validatorfw_string.IsURLHttpOrHttps().ValidateString(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}
