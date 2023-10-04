package string

import (
	"context"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"golang.org/x/exp/slices"
)

// Ensure our implementation satisfies the validator.String interface.
var _ validator.String = &urlHttpOrHttpsValidator{}

type urlHttpOrHttpsValidator struct{}

func (v urlHttpOrHttpsValidator) Description(_ context.Context) string {
	return "value must be a valid URL with host and http or https scheme"
}

func (v urlHttpOrHttpsValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v urlHttpOrHttpsValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()
	u, err := url.Parse(value)
	if err != nil || u.Host == "" || !slices.Contains([]string{"http", "https"}, u.Scheme) {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
			request.Path,
			v.Description(ctx),
			value,
		))
		return
	}
}

func IsURLHttpOrHttps() validator.String {
	return urlHttpOrHttpsValidator{}
}
