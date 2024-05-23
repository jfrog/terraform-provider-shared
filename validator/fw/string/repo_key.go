package string

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure our implementation satisfies the validator.String interface.
var _ validator.String = &repoKeyValidator{}

type repoKeyValidator struct{}

func (v repoKeyValidator) Description(_ context.Context) string {
	return "value must be a valid email address"
}

func (v repoKeyValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v repoKeyValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()

	if len(value) == 0 || len(value) > 64 {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueLengthDiagnostic(
			request.Path,
			"must be 1 - 64 alphanumeric and hyphen characters",
			value,
		))
	}

	if strings.ContainsAny(value, " !@#$%^&*()+={}[]:;<>,/?~`|\\") {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
			request.Path,
			"cannot contain spaces or special characters: !@#$%^&*()+={}[]:;<>,/?~`|\\",
			value,
		))
	}
}

func RepoKey() validator.String {
	return repoKeyValidator{}
}
