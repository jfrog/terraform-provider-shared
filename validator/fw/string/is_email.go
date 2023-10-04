package string

import (
	"context"
	"net/mail"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure our implementation satisfies the validator.String interface.
var _ validator.String = &emailValidator{}

type emailValidator struct{}

func (v emailValidator) Description(_ context.Context) string {
	return "value must be a valid email address"
}

func (v emailValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v emailValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()
	_, err := mail.ParseAddress(value)
	if err != nil {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
			request.Path,
			v.Description(ctx),
			value,
		))
	}
}

func IsEmail() validator.String {
	return emailValidator{}
}
