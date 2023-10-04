package string

import (
	"context"

	"github.com/gorhill/cronexpr"
	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure our implementation satisfies the validator.String interface.
var _ validator.String = &cronValidator{}

type cronValidator struct{}

func (v cronValidator) Description(_ context.Context) string {
	return "value must be a valid cron expression"
}

func (v cronValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v cronValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()

	_, err := cronexpr.Parse(value)
	if err != nil {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
			request.Path,
			v.Description(ctx),
			value,
		))
	}
}

func IsCron() validator.String {
	return cronValidator{}
}
