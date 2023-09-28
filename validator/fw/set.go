package fw

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	utilfw "github.com/jfrog/terraform-provider-shared/util/fw"
	"golang.org/x/exp/slices"
)

// Ensure our implementation satisfies the validator.Set interface.
var _ validator.Set = &stringInSliceValidator{}

type stringInSliceValidator struct {
	validValues []string
}

// Description returns a plaintext string describing the validator.
func (v stringInSliceValidator) Description(_ context.Context) string {
	return fmt.Sprintf("value must contain one of these: %v", v.validValues)
}

// MarkdownDescription returns a Markdown formatted string describing the validator.
func (v stringInSliceValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateBool performs the validation logic for the validator.
func (v stringInSliceValidator) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
	// If the current attribute configuration is null or unknown, there
	// cannot be any value comparisons, so exit early without error.
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	setValues := utilfw.StringSetToStrings(req.ConfigValue)
	for _, setValue := range setValues {
		if !slices.Contains(v.validValues, setValue) {
			resp.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
				req.Path,
				v.Description(ctx),
				setValue,
			))
		}
	}
}

// SetInSlice tests if the provided value matches the value of an element
// in the valid slice
func StringInSlice(validValues []string) validator.Set {
	return &stringInSliceValidator{
		validValues: validValues,
	}
}
